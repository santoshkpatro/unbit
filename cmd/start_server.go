package cmd

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/santoshkpatro/unbit/internal/config"
	"github.com/santoshkpatro/unbit/internal/handlers"
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

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	ctx := context.Background()

	db, err := config.NewPostgresConnection(ctx)
	if err != nil {
		log.Fatalf("❌ failed to connect to postgres: %v", err)
	}
	defer db.Close()

	rdb, err := config.NewRedisConnection(ctx)
	if err != nil {
		log.Fatalf("❌ failed to connect to redis: %v", err)
	}
	defer rdb.Close()

	hc := &handlers.HandlerContext{
		DB:  db,
		RDB: rdb,
	}

	apiRoutes := e.Group("/api")
	hc.RegisterAPIRoutes(apiRoutes)

	go func() {
		if err := e.Start(":" + config.Env.Port); err != nil {
			log.Printf("server stopped: %v", err)
		}
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
