package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	App      AppConfig
	HTTP     HTTPConfig
	DB       DBConfig
	Redis    RedisConfig
	JWT      JWTConfig
	Stripe   StripeConfig
	Security SecurityConfig
}

type AppConfig struct {
	Env  string
	Name string
}

type HTTPConfig struct {
	Host         string
	Port         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

type DBConfig struct {
	DSN          string
	MaxOpenConns int
	MaxIdleConns int
	MaxIdleTime  time.Duration
}

type RedisConfig struct {
	Addr     string
	Password string
	DB       int
}

type JWTConfig struct {
	AccessSecret  string
	RefreshSecret string
	AccessTTL     time.Duration
	RefreshTTL    time.Duration
}

type StripeConfig struct {
	SecretKey      string
	WebhookSecret  string
	PublishableKey string
}

type SecurityConfig struct {
	BcryptCost     int
	RateLimitRPM   int
	AllowedOrigins []string
}

func Load() (*Config, error) {
	return &Config{
		App: AppConfig{
			Env:  getEnv("APP_ENV", "development"),
			Name: getEnv("APP_NAME", "mfwebapp"),
		},
		HTTP: HTTPConfig{
			Host:         getEnv("HTTP_HOST", "0.0.0.0"),
			Port:         getEnv("HTTP_PORT", "8080"),
			ReadTimeout:  getDuration("HTTP_READ_TIMEOUT", 10*time.Second),
			WriteTimeout: getDuration("HTTP_WRITE_TIMEOUT", 10*time.Second),
		},
		DB: DBConfig{
			DSN:          getEnv("DATABASE_DSN", "postgres://postgres:postgres@localhost:5432/mfwebapp?sslmode=disable"),
			MaxOpenConns: getInt("DB_MAX_OPEN_CONNS", 25),
			MaxIdleConns: getInt("DB_MAX_IDLE_CONNS", 10),
			MaxIdleTime:  getDuration("DB_MAX_IDLE_TIME", 15*time.Minute),
		},
		Redis: RedisConfig{
			Addr:     getEnv("REDIS_ADDR", "localhost:6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getInt("REDIS_DB", 0),
		},
		JWT: JWTConfig{
			AccessSecret:  getEnv("JWT_ACCESS_SECRET", "change-me-access-secret"),
			RefreshSecret: getEnv("JWT_REFRESH_SECRET", "change-me-refresh-secret"),
			AccessTTL:     getDuration("JWT_ACCESS_TTL", 15*time.Minute),
			RefreshTTL:    getDuration("JWT_REFRESH_TTL", 7*24*time.Hour),
		},
		Stripe: StripeConfig{
			SecretKey:      getEnv("STRIPE_SECRET_KEY", ""),
			WebhookSecret:  getEnv("STRIPE_WEBHOOK_SECRET", ""),
			PublishableKey: getEnv("STRIPE_PUBLISHABLE_KEY", ""),
		},
		Security: SecurityConfig{
			BcryptCost:   getInt("BCRYPT_COST", 12),
			RateLimitRPM: getInt("RATE_LIMIT_RPM", 60),
			AllowedOrigins: []string{
				getEnv("CORS_USER_ORIGIN", "http://localhost:3000"),
				getEnv("CORS_ADMIN_ORIGIN", "http://localhost:3001"),
			},
		},
	}, nil
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getInt(key string, fallback int) int {
	if v := os.Getenv(key); v != "" {
		if i, err := strconv.Atoi(v); err == nil {
			return i
		}
	}
	return fallback
}

func getDuration(key string, fallback time.Duration) time.Duration {
	if v := os.Getenv(key); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			return d
		}
	}
	return fallback
}
