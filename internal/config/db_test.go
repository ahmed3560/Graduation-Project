package config

import (
	"strings"
	"testing"
)

// setDBOnly clears every variable, then sets only DATABASE_URL.
// allKeys and testDBURL come from config_test.go in this package.
func setDBOnly(t *testing.T) {
	t.Helper()
	for _, k := range allKeys {
		t.Setenv(k, "")
	}
	t.Setenv("DATABASE_URL", testDBURL)
}

func TestLoadDB_WorksWithoutJWTSecret(t *testing.T) {
	setDBOnly(t)

	cfg, err := LoadDB()
	if err != nil {
		t.Fatalf("LoadDB() error = %v", err)
	}
	if cfg.DB.URL != testDBURL {
		t.Errorf("DB.URL = %q, want %q", cfg.DB.URL, testDBURL)
	}
	if cfg.Env != EnvDevelopment {
		t.Errorf("Env = %q, want %q", cfg.Env, EnvDevelopment)
	}
}

func TestLoadDB_RequiresDatabaseURL(t *testing.T) {
	setDBOnly(t)
	t.Setenv("DATABASE_URL", "")

	_, err := LoadDB()
	if err == nil || !strings.Contains(err.Error(), "DATABASE_URL is required") {
		t.Fatalf("err = %v, want DATABASE_URL is required", err)
	}
}

func TestLoadDB_RejectsUnknownEnv(t *testing.T) {
	setDBOnly(t)
	t.Setenv("APP_ENV", "staging")

	_, err := LoadDB()
	if err == nil || !strings.Contains(err.Error(), "APP_ENV must be") {
		t.Fatalf("err = %v, want APP_ENV error", err)
	}
}

func TestLoadDB_ReportsAllProblemsAtOnce(t *testing.T) {
	setDBOnly(t)
	t.Setenv("DATABASE_URL", "")
	t.Setenv("APP_ENV", "staging")

	_, err := LoadDB()
	if err == nil {
		t.Fatal("want an error")
	}
	for _, want := range []string{"DATABASE_URL is required", "APP_ENV must be"} {
		if !strings.Contains(err.Error(), want) {
			t.Errorf("error %q is missing %q", err, want)
		}
	}
}
