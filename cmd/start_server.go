package cmd

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/santoshkpatro/unbit/internal/config"
	"github.com/santoshkpatro/unbit/internal/worker"
	"github.com/spf13/cobra"
)

type CustomValidator struct {
	validator *validator.Validate
}

var startServerCmd = &cobra.Command{
	Use:   "start_server",
	Short: "Start Server",
	RunE: func(cmd *cobra.Command, args []string) error {
		return startServer()
	},
}

func startServer() error {
	e := echo.New()

	// e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.StaticWithConfig(middleware.StaticConfig{Root: "dist", HTML5: true}))

	// Set up custom validator for Request validation
	e.Validator = &CustomValidator{validator: validator.New()}

	ctx := context.Background()

	db, err := config.NewPostgresConnection(ctx)
	if err != nil {
		log.Fatalf("❌ failed to connect to postgres: %v", err)
	}
	defer db.Close()

	cache, err := config.NewRedisConnection(ctx)
	if err != nil {
		log.Fatalf("❌ failed to connect to redis: %v", err)
	}
	defer cache.Close()

	config.RegisterRoutes(e, db, cache)

	// Start server
	go func() {
		if err := e.Start(":" + config.Env.Port); err != nil {
			log.Printf("server stopped: %v", err)
		}
	}()

	// Start Worker
	go func() {
		queueName := "issues"
		worker.StartWorker(cache, db, queueName)
	}()

	// Wait for Ctrl+C or SIGTERM
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Println("⚙️  Shutting down gracefully...")

	ctxShutdown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := e.Shutdown(ctxShutdown); err != nil {
		log.Fatalf("❌ Server forced to shutdown: %v", err)
	}

	log.Println("✅ Server exited properly")

	return nil
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}
