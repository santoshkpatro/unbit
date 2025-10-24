package cmd

import (
	"context"
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
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
	fmt.Println("🚀 Running installation setup...")

	if err := createSettingsTable(ctx, db); err != nil {
		return fmt.Errorf("failed to create settings table: %w", err)
	}

	if err := insertDefaultSettings(ctx, db); err != nil {
		return fmt.Errorf("failed to insert default settings: %w", err)
	}

	fmt.Println("✅ Installation complete!")
	return nil
}

func createSettingsTable(ctx context.Context, db *sqlx.DB) error {
	query := `
		CREATE TABLE IF NOT EXISTS setting (
			org_name TEXT NOT NULL DEFAULT '',
			org_url TEXT NOT NULL DEFAULT '',
			support_email TEXT NOT NULL DEFAULT '',
			schema_migrations TEXT[] NOT NULL DEFAULT ARRAY[]::TEXT[],
			allow_invite BOOLEAN NOT NULL DEFAULT true,
			allow_login BOOLEAN NOT NULL DEFAULT true,
			is_maintenance_on BOOLEAN NOT NULL DEFAULT false,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
		);
	`
	_, err := db.ExecContext(ctx, query)

	return err
}

func insertDefaultSettings(ctx context.Context, db *sqlx.DB) error {
	query := `
	INSERT INTO setting (org_name, org_url, support_email)
	SELECT 'My Org', 'https://unbit.app', 'support@unbit.com'
	WHERE NOT EXISTS (SELECT 1 FROM setting);
	`
	_, err := db.ExecContext(ctx, query)
	return err
}
