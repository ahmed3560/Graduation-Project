package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/ahmed3560/Graduation-Project/backend/internal/config"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	if err := run(logger); err != nil {
		logger.Error("api exited with error", "err", err)
		os.Exit(1)
	}
}

func run(logger *slog.Logger) error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	r := chi.NewRouter()
	r.Use(middleware.Recoverer) // turns a panic in a handler into a 500 instead of crashing the process
	r.Get("/healthz", healthz)

	srv := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      r,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	// ctx is cancelled on Ctrl+C (SIGINT) or on SIGTERM, which Docker sends on stop.
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	serverErr := make(chan error, 1)
	go func() {
		logger.Info("api starting", "port", cfg.Server.Port, "env", cfg.Env)
		serverErr <- srv.ListenAndServe()
	}()

	select {
	case err := <-serverErr:
		// The server stopped on its own, which means it failed to start or crashed.
		return err
	case <-ctx.Done():
	}

	logger.Info("shutting down", "timeout", cfg.Server.ShutdownTimeout.String())

	// Stop accepting new connections and wait for in-flight requests to finish.
	shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.Server.ShutdownTimeout)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		return err
	}

	// After Shutdown, ListenAndServe returns http.ErrServerClosed. That's expected, not a failure.
	if err := <-serverErr; err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	logger.Info("api stopped")
	return nil
}

func healthz(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{"status":"ok"}`))
}