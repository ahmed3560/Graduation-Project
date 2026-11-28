package config

import (
	"strings"
	"testing"
)

const testDBURL = "postgres://wasla:secret-pass@localhost:5432/wasla?sslmode=disable"

// allKeys lists every variable Load reads, so tests start from a clean slate.
var allKeys = []string{
	"APP_ENV", "PORT",
	"SERVER_READ_TIMEOUT", "SERVER_WRITE_TIMEOUT", "SERVER_IDLE_TIMEOUT", "SERVER_SHUTDOWN_TIMEOUT",
	"DATABASE_URL", "DB_MAX_CONNS", "DB_MIN_CONNS", "DB_MAX_CONN_LIFETIME",
	"JWT_SECRET", "JWT_ACCESS_TTL", "JWT_REFRESH_TTL",
}

// setValid clears every variable, then sets the minimum valid configuration.
// An empty value counts as unset in Load, so t.Setenv(key, "") clears a key.
func setValid(t *testing.T) {
	t.Helper()
	for _, k := range allKeys {
		t.Setenv(k, "")
	}
	t.Setenv("DATABASE_URL", testDBURL)
	t.Setenv("JWT_SECRET", strings.Repeat("a", 32))
}

func TestLoad_DefaultsApplyWhenUnset(t *testing.T) {
	setValid(t)

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if cfg.Env != EnvDevelopment {
		t.Errorf("Env = %q, want %q", cfg.Env, EnvDevelopment)
	}
	if cfg.Server.Port != "8080" {
		t.Errorf("Port = %q, want 8080", cfg.Server.Port)
	}
	if cfg.DB.MaxConns != 25 || cfg.DB.MinConns != 5 {
		t.Errorf("pool sizes = %d/%d, want 25/5", cfg.DB.MaxConns, cfg.DB.MinConns)
	}
	if cfg.Auth.AccessTokenTTL.String() != "15m0s" {
		t.Errorf("AccessTokenTTL = %s, want 15m0s", cfg.Auth.AccessTokenTTL)
	}
}

func TestLoad_SecretLengthBoundary(t *testing.T) {
	setValid(t)

	t.Setenv("JWT_SECRET", strings.Repeat("a", 31))
	if _, err := Load(); err == nil || !strings.Contains(err.Error(), "at least 32") {
		t.Fatalf("31-char secret: err = %v, want length error", err)
	}

	t.Setenv("JWT_SECRET", strings.Repeat("a", 32))
	if _, err := Load(); err != nil {
		t.Fatalf("32-char secret: unexpected error %v", err)
	}
}

func TestLoad_InvalidValuesReturnErrors(t *testing.T) {
	cases := []struct {
		name      string
		key       string
		value     string
		wantInMsg string
	}{
		{"missing database url", "DATABASE_URL", "", "DATABASE_URL is required"},
		{"non-numeric pool size", "DB_MAX_CONNS", "ten", "DB_MAX_CONNS"},
		{"port out of range", "PORT", "70000", "PORT"},
		{"malformed duration", "JWT_ACCESS_TTL", "15 minutes", "JWT_ACCESS_TTL"},
		{"zero duration", "SERVER_READ_TIMEOUT", "0s", "SERVER_READ_TIMEOUT"},
		{"min conns above max", "DB_MIN_CONNS", "30", "DB_MIN_CONNS"},
		{"refresh not longer than access", "JWT_REFRESH_TTL", "10m", "JWT_REFRESH_TTL"},
		{"unknown environment", "APP_ENV", "prod", "APP_ENV"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			setValid(t)
			t.Setenv(tc.key, tc.value)

			_, err := Load()
			if err == nil {
				t.Fatalf("Load() returned no error for %s=%q", tc.key, tc.value)
			}
			if !strings.Contains(err.Error(), tc.wantInMsg) {
				t.Errorf("error %q does not mention %q", err, tc.wantInMsg)
			}
		})
	}
}

func TestConfigString_RedactsSecrets(t *testing.T) {
	setValid(t)

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	s := cfg.String()

	if strings.Contains(s, "secret-pass") {
		t.Errorf("String() leaks database password: %s", s)
	}
	if strings.Contains(s, strings.Repeat("a", 32)) {
		t.Errorf("String() leaks JWT secret: %s", s)
	}
	if !strings.Contains(s, "REDACTED") {
		t.Errorf("String() shows no redaction marker: %s", s)
	}
}