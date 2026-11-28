package main

import (
	"errors"
	"fmt"
	"strconv"
)

const usage = `usage:
  migrate up          apply every migration that has not been applied yet
  migrate down N      undo the last N migrations (N is required, on purpose)
  migrate version     show which migration the database is at

The database comes from the DATABASE_URL environment variable.`

// command is what the person typed, checked and ready to run.
type command struct {
	name  string // "up", "down" or "version"
	steps int    // only used by "down": how many migrations to undo
}

// parseCommand turns the program's arguments into a command, or says what is wrong.
// It does not touch the database, so it is easy to test.
func parseCommand(args []string) (command, error) {
	if len(args) == 0 {
		return command{}, errors.New("missing command")
	}

	switch args[0] {
	case "up", "version":
		if len(args) != 1 {
			return command{}, fmt.Errorf("%s takes no arguments", args[0])
		}
		return command{name: args[0]}, nil

	case "down":
		if len(args) != 2 {
			return command{}, errors.New("down needs exactly one number, for example: migrate down 1")
		}
		n, err := strconv.Atoi(args[1])
		if err != nil || n < 1 {
			return command{}, fmt.Errorf("down needs a whole number of 1 or more, got %q", args[1])
		}
		return command{name: "down", steps: n}, nil

	default:
		return command{}, fmt.Errorf("unknown command %q", args[0])
	}
}

// checkAllowed stops commands that must never run in some environments.
// Undoing a migration can destroy data, so it is blocked in production:
// there we fix problems by adding a new migration instead.
func (c command) checkAllowed(env string) error {
	if c.name == "down" && env == "production" {
		return errors.New("down is disabled when APP_ENV=production; fix forward with a new migration instead")
	}
	return nil
}
