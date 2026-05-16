package main

import (
	"context"
	"errors"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/hajimohammadinet/dabir/internal/bootstrap"
)

func main() {
	ctx := context.Background()

	app, err := bootstrap.New(ctx)
	if err != nil {
		log.Fatalf("failed to bootstrap application: %v", err)
	}

	serverErrors := make(chan error, 1)

	go func() {
		serverErrors <- app.Start()
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-serverErrors:
		if err != nil {
			log.Fatalf("server error: %v", err)
		}

	case sig := <-shutdown:
		app.Logger.Info("shutdown signal received", "signal", sig.String())

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := app.Shutdown(shutdownCtx); err != nil {
			if !errors.Is(err, context.DeadlineExceeded) {
				app.Logger.Error("graceful shutdown failed", "error", err)
			}
		}

		app.Logger.Info("server stopped")
	}
}
