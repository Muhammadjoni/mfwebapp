package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"regexp"
	"strings"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
	"github.com/muhammadjoni/mfwebapp/internal/config"
	v1 "github.com/muhammadjoni/mfwebapp/internal/handler/http/v1"
	"github.com/muhammadjoni/mfwebapp/internal/handler/middleware"
	"github.com/muhammadjoni/mfwebapp/internal/infrastructure/database"
	pgRepo "github.com/muhammadjoni/mfwebapp/internal/repository/postgres"
	"github.com/muhammadjoni/mfwebapp/internal/service"
	"github.com/muhammadjoni/mfwebapp/migrations"
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

	// ── Auto-migrate ────────────────────────────────────────────────────────
	if err := runMigrations(ctx, cfg.DB.DSN, log); err != nil {
		log.Fatal("migration failed", zap.Error(err))
	}
	log.Info("migrations applied")

	// ── Database pool ────────────────────────────────────────────────────────
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

	// ── Infrastructure ───────────────────────────────────────────────────────
	jwtMgr := jwt.NewManager(
		cfg.JWT.AccessSecret,
		cfg.JWT.RefreshSecret,
		cfg.JWT.AccessTTL,
		cfg.JWT.RefreshTTL,
	)
	hasher := hash.NewHasher(cfg.Security.BcryptCost)

	// ── Repositories ─────────────────────────────────────────────────────────
	userRepo := pgRepo.NewUserRepository(pool)
	productRepo := pgRepo.NewProductRepository(pool)
	categoryRepo := pgRepo.NewCategoryRepository(pool)
	orderRepo := pgRepo.NewOrderRepository(pool)
	cartRepo := pgRepo.NewCartRepository(pool)
	sellerRepo := pgRepo.NewSellerRepository(pool)

	// ── Services ─────────────────────────────────────────────────────────────
	authSvc := service.NewAuthService(userRepo, jwtMgr, hasher, &cfg.JWT)
	productSvc := service.NewProductService(productRepo, categoryRepo)
	orderSvc := service.NewOrderService(orderRepo, productRepo, cartRepo)
	cartSvc := service.NewCartService(cartRepo)
	userSvc := service.NewUserService(userRepo)
	sellerSvc := service.NewSellerService(sellerRepo)

	// ── Handlers ─────────────────────────────────────────────────────────────
	authHandler := v1.NewAuthHandler(authSvc)
	productHandler := v1.NewProductHandler(productSvc)
	orderHandler := v1.NewOrderHandler(orderSvc)
	cartHandler := v1.NewCartHandler(cartSvc)
	adminHandler := v1.NewAdminHandler(userSvc, sellerSvc, orderSvc, productSvc)

	// ── Middleware & Router ───────────────────────────────────────────────────
	authMW := middleware.NewAuthMiddleware(jwtMgr)
	router := v1.NewRouter(
		authHandler,
		productHandler,
		orderHandler,
		cartHandler,
		adminHandler,
		authMW,
		cfg.Security.AllowedOrigins,
		cfg.Security.RateLimitRPM,
	)

	// ── HTTP server ───────────────────────────────────────────────────────────
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

// runMigrations checks if the schema exists and applies it if not.
// It uses a single pgx connection (not the pool) so it can run before the pool.
func runMigrations(ctx context.Context, dsn string, log *zap.Logger) error {
	conn, err := pgx.Connect(ctx, dsn)
	if err != nil {
		return fmt.Errorf("connect: %w", err)
	}
	defer conn.Close(ctx)

	var exists bool
	if err := conn.QueryRow(ctx, `
		SELECT EXISTS (
			SELECT FROM information_schema.tables
			WHERE table_schema='public' AND table_name='users'
		)`).Scan(&exists); err != nil {
		return fmt.Errorf("check schema: %w", err)
	}
	if exists {
		log.Info("schema already exists, skipping migration")
		return nil
	}

	log.Info("applying initial migration")
	sqlBytes, err := migrations.FS.ReadFile("001_init_schema.sql")
	if err != nil {
		return fmt.Errorf("read migration: %w", err)
	}

	sql := extractGooseUp(string(sqlBytes))
	if _, err := conn.Exec(ctx, sql); err != nil {
		return fmt.Errorf("execute migration: %w", err)
	}
	return nil
}

var gooseDirective = regexp.MustCompile(`-- \+goose\s+\S+\s*\n?`)

// extractGooseUp strips goose directives and returns only the Up section.
func extractGooseUp(content string) string {
	downIdx := strings.Index(content, "-- +goose Down")
	if downIdx > 0 {
		content = content[:downIdx]
	}
	upIdx := strings.Index(content, "-- +goose Up")
	if upIdx >= 0 {
		content = content[upIdx:]
	}
	return strings.TrimSpace(gooseDirective.ReplaceAllString(content, ""))
}
