// Package migrator applies and undoes database migrations.
// It is a thin wrapper around the golang-migrate library, so the rest of our code
// never has to know how that library works.
package migrator

import (
	"errors"
	"fmt"
	"io/fs"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	// The blank import registers the "pgx5" database driver with golang-migrate.
	_ "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

// Migrator runs migrations against one database.
type Migrator struct {
	m *migrate.Migrate
}

// New connects to the database and reads migration files from files.
// The files must sit at the root of files (the way embed.FS in package migrations does).
func New(databaseURL string, files fs.FS) (*Migrator, error) {
	src, err := iofs.New(files, ".")
	if err != nil {
		return nil, fmt.Errorf("read migration files: %w", err)
	}
	m, err := migrate.NewWithSourceInstance("iofs", src, driverURL(databaseURL))
	if err != nil {
		return nil, fmt.Errorf("connect to the database: %w", err)
	}
	return &Migrator{m: m}, nil
}

// Up applies every migration that has not been applied yet.
// Having nothing to apply is fine, so running Up twice is harmless.
func (m *Migrator) Up() error {
	err := m.m.Up()
	if errors.Is(err, migrate.ErrNoChange) {
		return nil
	}
	return err
}

// Down undoes the last n migrations, newest first.
//
// It asks for a number on purpose. There is no "undo everything" shortcut,
// so a typo cannot wipe the schema.
func (m *Migrator) Down(n int) error {
	if n < 1 {
		return fmt.Errorf("down needs a number of 1 or more, got %d", n)
	}
	version, dirty, err := m.Version()
	if err != nil {
		return err
	}
	if dirty {
		return fmt.Errorf("database is dirty at version %d: a migration failed half way and must be repaired by hand first", version)
	}
	// Versions are consecutive (migrations_test.go checks it), so the current
	// version number is also the number of migrations applied so far.
	if n > int(version) {
		return fmt.Errorf("cannot undo %d migrations: only %d applied", n, version)
	}
	return m.m.Steps(-n)
}

// Version returns the number of the newest applied migration.
// An empty database returns 0. dirty is true when the last migration failed
// half way; nothing more can be applied until someone repairs the database.
func (m *Migrator) Version() (version uint, dirty bool, err error) {
	version, dirty, err = m.m.Version()
	if errors.Is(err, migrate.ErrNilVersion) {
		return 0, false, nil
	}
	return version, dirty, err
}

// Close releases the database connection.
func (m *Migrator) Close() error {
	srcErr, dbErr := m.m.Close()
	return errors.Join(srcErr, dbErr)
}

// driverURL changes postgres:// into pgx5://, the name golang-migrate gives its
// pgx driver. Our config uses the normal postgres:// form everywhere else.
func driverURL(databaseURL string) string {
	for _, prefix := range []string{"postgres://", "postgresql://"} {
		if rest, ok := strings.CutPrefix(databaseURL, prefix); ok {
			return "pgx5://" + rest
		}
	}
	return databaseURL
}
