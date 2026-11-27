package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type Config struct {
	Env string // "development" or "production"

	Server struct {
		Port            string
		ReadTimeout     time.Duration
		WriteTimeout    time.Duration
		IdleTimeout     time.Duration
		ShutdownTimeout time.Duration
	}

	DB struct {
		URL             string
		MaxConns        int32
		MinConns        int32
		MaxConnLifetime time.Duration
	}

	Auth struct {
		JWTSecret       string
		AccessTokenTTL  time.Duration
		RefreshTokenTTL time.Duration
	}
}

func Load() (Config, error) {
	var cfg Config

	cfg.Env = getEnv("APP_ENV", "development")

	// Server
	cfg.Server.Port = getEnv("PORT", "8080")
	cfg.Server.ReadTimeout = getDuration("SERVER_READ_TIMEOUT", 10*time.Second)
	cfg.Server.WriteTimeout = getDuration("SERVER_WRITE_TIMEOUT", 15*time.Second)
	cfg.Server.IdleTimeout = getDuration("SERVER_IDLE_TIMEOUT", 60*time.Second)
	cfg.Server.ShutdownTimeout = getDuration("SERVER_SHUTDOWN_TIMEOUT", 10*time.Second)

	// Database
	cfg.DB.URL = os.Getenv("DATABASE_URL")
	if cfg.DB.URL == "" {
		return Config{}, fmt.Errorf("config: DATABASE_URL is required")
	}
	cfg.DB.MaxConns = int32(getInt("DB_MAX_CONNS", 25))
	cfg.DB.MinConns = int32(getInt("DB_MIN_CONNS", 5))
	cfg.DB.MaxConnLifetime = getDuration("DB_MAX_CONN_LIFETIME", 30*time.Minute)

	// Auth
	cfg.Auth.JWTSecret = os.Getenv("JWT_SECRET")
	if cfg.Auth.JWTSecret == "" {
		return Config{}, fmt.Errorf("config: JWT_SECRET is required")
	}
	if len(cfg.Auth.JWTSecret) < 32 {
		return Config{}, fmt.Errorf("config: JWT_SECRET must be at least 32 characters")
	}
	cfg.Auth.AccessTokenTTL = getDuration("JWT_ACCESS_TTL", 15*time.Minute)
	cfg.Auth.RefreshTokenTTL = getDuration("JWT_REFRESH_TTL", 7*24*time.Hour)

	return cfg, nil
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getInt(key string, fallback int) int {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	parsed, err := strconv.Atoi(v)
	if err != nil {
		return fallback
	}
	return parsed
}

func getDuration(key string, fallback time.Duration) time.Duration {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	parsed, err := time.ParseDuration(v)
	if err != nil {
		return fallback
	}
	return parsed
}