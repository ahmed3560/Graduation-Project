// Command migrate creates and changes the database tables.
//
//	go run ./cmd/migrate up
//	go run ./cmd/migrate down 1
//	go run ./cmd/migrate version
package main

import (
	"errors"
	"fmt"
	"log/slog"
	"os"

	"github.com/ahmed3560/Graduation-Project/backend/internal/config"
	"github.com/ahmed3560/Graduation-Project/backend/internal/migrator"
	"github.com/ahmed3560/Graduation-Project/backend/migrations"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	// Check the arguments first, so a typo prints help instead of touching the database.
	cmd, err := parseCommand(os.Args[1:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n\n%s\n", err, usage)
		os.Exit(2)
	}

	if err := run(logger, cmd); err != nil {
		logger.Error("migrate failed", "err", err)
		os.Exit(1)
	}
}

func run(logger *slog.Logger, cmd command) (err error) {
	cfg, err := config.LoadDB()
	if err != nil {
		return err
	}
	if err := cmd.checkAllowed(cfg.Env); err != nil {
		return err
	}

	m, err := migrator.New(cfg.DB.URL, migrations.FS)
	if err != nil {
		return err
	}
	// Close runs last. If closing fails and nothing else failed, report that.
	defer func() { err = errors.Join(err, m.Close()) }()

	switch cmd.name {
	case "up":
		err = m.Up()
	case "down":
		err = m.Down(cmd.steps)
	case "version":
		// Nothing to change; the version is logged below.
	}
	if err != nil {
		return err
	}

	version, dirty, err := m.Version()
	if err != nil {
		return err
	}
	// Never log cfg.DB.URL: it contains the database password.
	logger.Info("migrate finished", "command", cmd.name, "version", version, "dirty", dirty)
	return nil
}
