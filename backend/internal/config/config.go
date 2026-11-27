package config

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"strconv"
	"time"
)

const (
	EnvDevelopment = "development"
	EnvTest        = "test"
	EnvProduction  = "production"
)

type Config struct {
	Env string // EnvDevelopment, EnvTest, or EnvProduction

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

// Load reads every setting, then checks them together.
// It returns all problems at once, so the operator can fix them in one pass.
func Load() (Config, error) {
	l := &loader{}
	var cfg Config

	cfg.Env = l.str("APP_ENV", EnvDevelopment)

	cfg.Server.Port = strconv.Itoa(l.intRange("PORT", 8080, 1, 65535))
	cfg.Server.ReadTimeout = l.duration("SERVER_READ_TIMEOUT", 10*time.Second)
	cfg.Server.WriteTimeout = l.duration("SERVER_WRITE_TIMEOUT", 15*time.Second)
	cfg.Server.IdleTimeout = l.duration("SERVER_IDLE_TIMEOUT", 60*time.Second)
	cfg.Server.ShutdownTimeout = l.duration("SERVER_SHUTDOWN_TIMEOUT", 10*time.Second)

	cfg.DB.URL = l.required("DATABASE_URL")
	cfg.DB.MaxConns = int32(l.intRange("DB_MAX_CONNS", 25, 1, 1000))
	cfg.DB.MinConns = int32(l.intRange("DB_MIN_CONNS", 5, 0, 1000))
	cfg.DB.MaxConnLifetime = l.duration("DB_MAX_CONN_LIFETIME", 30*time.Minute)

	cfg.Auth.JWTSecret = l.required("JWT_SECRET")
	cfg.Auth.AccessTokenTTL = l.duration("JWT_ACCESS_TTL", 15*time.Minute)
	cfg.Auth.RefreshTokenTTL = l.duration("JWT_REFRESH_TTL", 7*24*time.Hour)

	l.errs = append(l.errs, cfg.validate()...)
	if len(l.errs) > 0 {
		return Config{}, errors.Join(l.errs...)
	}
	return cfg, nil
}

// validate checks rules that involve more than one setting.
// It skips empty values, because the loader has already reported those.
func (c Config) validate() []error {
	var errs []error

	switch c.Env {
	case EnvDevelopment, EnvTest, EnvProduction:
	default:
		errs = append(errs, fmt.Errorf("APP_ENV must be development, test, or production, got %q", c.Env))
	}

	if c.Auth.JWTSecret != "" && len(c.Auth.JWTSecret) < 32 {
		errs = append(errs, errors.New("JWT_SECRET must be at least 32 characters"))
	}
	if c.DB.MinConns > c.DB.MaxConns {
		errs = append(errs, fmt.Errorf("DB_MIN_CONNS (%d) must not exceed DB_MAX_CONNS (%d)", c.DB.MinConns, c.DB.MaxConns))
	}
	if c.Auth.RefreshTokenTTL <= c.Auth.AccessTokenTTL {
		errs = append(errs, errors.New("JWT_REFRESH_TTL must be longer than JWT_ACCESS_TTL"))
	}
	return errs
}

// String implements fmt.Stringer, so fmt.Println and %v never print secrets.
// Note: %#v and json.Marshal bypass this method, so don't log the struct that way.
func (c Config) String() string {
	return fmt.Sprintf(
		"Config{Env:%s Port:%s DBURL:%s DBMaxConns:%d DBMinConns:%d JWTSecret:[REDACTED] AccessTokenTTL:%s RefreshTokenTTL:%s}",
		c.Env, c.Server.Port, redactURL(c.DB.URL), c.DB.MaxConns, c.DB.MinConns,
		c.Auth.AccessTokenTTL, c.Auth.RefreshTokenTTL,
	)
}

// redactURL hides the password in a URL like postgres://user:pass@host/db.
func redactURL(raw string) string {
	u, err := url.Parse(raw)
	if err != nil {
		return "[REDACTED]"
	}
	if u.User == nil {
		return u.String()
	}
	if _, hasPassword := u.User.Password(); hasPassword {
		u.User = url.UserPassword(u.User.Username(), "REDACTED")
	}
	return u.String()
}

// loader collects errors while reading environment variables.
type loader struct {
	errs []error
}

func (l *loader) fail(format string, args ...any) {
	l.errs = append(l.errs, fmt.Errorf(format, args...))
}

// str returns the variable or the fallback when it's unset or empty.
func (l *loader) str(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

// required returns the variable and records an error if it's unset or empty.
func (l *loader) required(key string) string {
	v := os.Getenv(key)
	if v == "" {
		l.fail("%s is required", key)
	}
	return v
}

// intRange parses a whole number. A value that is set but invalid is an error,
// not a silent fallback.
func (l *loader) intRange(key string, fallback, min, max int) int {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	n, err := strconv.Atoi(v)
	if err != nil || n < min || n > max {
		l.fail("%s must be a whole number from %d to %d, got %q", key, min, max, v)
		return fallback
	}
	return n
}

// duration parses a positive duration such as "15m" or "168h".
func (l *loader) duration(key string, fallback time.Duration) time.Duration {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	d, err := time.ParseDuration(v)
	if err != nil || d <= 0 {
		l.fail("%s must be a positive duration such as 15m or 168h, got %q", key, v)
		return fallback
	}
	return d
}