package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/muhammadjoni/mfwebapp/internal/config"
	v1 "github.com/muhammadjoni/mfwebapp/internal/handler/http/v1"
	"github.com/muhammadjoni/mfwebapp/internal/handler/middleware"
	"github.com/muhammadjoni/mfwebapp/internal/infrastructure/database"
	"github.com/muhammadjoni/mfwebapp/pkg/hash"
	"github.com/muhammadjoni/mfwebapp/pkg/jwt"
	"github.com/muhammadjoni/mfwebapp/pkg/logger"
	"go.uber.org/zap"
)

func main() {
	_ = godotenv.Load()

	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "config error: %v\n", err)
		os.Exit(1)
	}

	log, err := logger.New(cfg.App.Env)
	if err != nil {
		fmt.Fprintf(os.Stderr, "logger error: %v\n", err)
		os.Exit(1)
	}
	defer log.Sync() //nolint:errcheck

	ctx := context.Background()

	pool, err := database.NewPool(ctx, database.Config{
		DSN:          cfg.DB.DSN,
		MaxOpenConns: cfg.DB.MaxOpenConns,
		MaxIdleConns: cfg.DB.MaxIdleConns,
		MaxIdleTime:  cfg.DB.MaxIdleTime,
	})
	if err != nil {
		log.Fatal("database connection failed", zap.Error(err))
	}
	defer pool.Close()
	log.Info("database connected")

	jwtMgr := jwt.NewManager(
		cfg.JWT.AccessSecret,
		cfg.JWT.RefreshSecret,
		cfg.JWT.AccessTTL,
		cfg.JWT.RefreshTTL,
	)
	hasher := hash.NewHasher(cfg.Security.BcryptCost)

	// Suppress unused warnings until real repos are wired
	_ = hasher
	_ = pool

	authMW := middleware.NewAuthMiddleware(jwtMgr)

	// TODO: wire real postgres repository implementations then inject into services
	// e.g. userRepo := postgres.NewUserRepository(pool)
	//      authSvc  := service.NewAuthService(userRepo, jwtMgr, hasher, &cfg.JWT)
	var authHandler *v1.AuthHandler
	var productHandler *v1.ProductHandler
	var orderHandler *v1.OrderHandler

	router := v1.NewRouter(
		authHandler,
		productHandler,
		orderHandler,
		authMW,
		cfg.Security.AllowedOrigins,
		cfg.Security.RateLimitRPM,
	)

	srv := &http.Server{
		Addr:         cfg.HTTP.Host + ":" + cfg.HTTP.Port,
		Handler:      router.Build(),
		ReadTimeout:  cfg.HTTP.ReadTimeout,
		WriteTimeout: cfg.HTTP.WriteTimeout,
		IdleTimeout:  60 * time.Second,
	}

	done := make(chan struct{})
	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit
		log.Info("graceful shutdown initiated")
		shutCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		if err := srv.Shutdown(shutCtx); err != nil {
			log.Error("shutdown error", zap.Error(err))
		}
		close(done)
	}()

	log.Info("server listening", zap.String("addr", srv.Addr))
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal("server error", zap.Error(err))
	}
	<-done
	log.Info("server stopped")
}
