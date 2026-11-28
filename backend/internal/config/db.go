package config

import (
	"errors"
	"fmt"
)

// LoadDB reads only the settings a database-only program needs, such as cmd/migrate.
// Load asks for JWT_SECRET because the API needs it; migrations do not, so a person
// running migrations should not have to set it.
//
// Only Env and DB.URL are filled in. Every other field stays at its zero value.
func LoadDB() (Config, error) {
	l := &loader{}
	var cfg Config

	cfg.Env = l.str("APP_ENV", EnvDevelopment)
	cfg.DB.URL = l.required("DATABASE_URL")

	// Same APP_ENV rule as Config.validate. We cannot call validate itself here,
	// because it also checks the JWT and pool settings that this loader never reads.
	switch cfg.Env {
	case EnvDevelopment, EnvTest, EnvProduction:
	default:
		l.errs = append(l.errs, fmt.Errorf("APP_ENV must be development, test, or production, got %q", cfg.Env))
	}

	if len(l.errs) > 0 {
		return Config{}, errors.Join(l.errs...)
	}
	return cfg, nil
}
