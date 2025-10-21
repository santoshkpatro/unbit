package cmd

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/santoshkpatro/unbit/internal/config"
	"github.com/spf13/cobra"
)

var startServerCmd = &cobra.Command{
	Use:   "start_server",
	Short: "Start Server",
	RunE: func(cmd *cobra.Command, args []string) error {
		return startServer()
	},
}

func startServer() error {
	e := echo.New()

	errCh := make(chan error, 1)

	go func() {
		errCh <- e.Start(":" + config.Env.Port)
	}()

	// OS signal handling for graceful shutdown
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	select {
	case sig := <-sigCh:
		// received termination signal — shut down gracefully
		e.Logger.Infof("shutting down on signal: %v", sig)
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := e.Shutdown(ctx); err != nil {
			return err
		}
		return nil
	case err := <-errCh:
		// server returned an error (could be non-nil on failure)
		if err != nil && err != echo.ErrInternalServerError {
			return err
		}
		return nil
	}
}
