package cmd

import (
	"context"
	"log"

	"github.com/santoshkpatro/unbit/internal/config"
	"github.com/spf13/cobra"
)

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install",
	RunE: func(cmd *cobra.Command, args []string) error {
		return install()
	},
}

func install() error {
	ctx := context.Background()
	db, err := config.NewPostgresConnection(ctx)
	if err != nil {
		log.Fatalf("❌ failed to connect to postgres: %v", err)
	}
	defer db.Close()

	// 1. Setup initial schema table in the schema_migrations
	// 2. Setup setting table with some seed data

	return nil
}
