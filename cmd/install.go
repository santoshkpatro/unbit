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

	// create settings table
	if err := createSettingsTable(ctx, db); err != nil {
		return fmt.Errorf("failed to create settings table: %w", err)
	}

	// seed a few settings
	if err := seedSettings(ctx, db); err != nil {
		return fmt.Errorf("failed to seed settings: %w", err)
	}

	// create schema_migrations table
	if err := createSchemaMigrationsTable(ctx, db); err != nil {
		return fmt.Errorf("failed to create schema_migrations table: %w", err)
	}

	fmt.Println("✅ Installation complete!")
	return nil
}

// createSchemaMigrationsTable creates schema_migrations with an integer 'timestamp' and applied_at timestamptz.
// Note: column name 'timestamp' is used as requested (type BIGINT).
func createSchemaMigrationsTable(ctx context.Context, db *sqlx.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS schema_migrations (
		timestamp BIGINT PRIMARY KEY,
		applied_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
	);
	`
	_, err := db.ExecContext(ctx, query)
	return err
}

// createSettingsTable creates settings table with key (unique), value (jsonb) and updated_at timestamptz.
func createSettingsTable(ctx context.Context, db *sqlx.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS settings (
		key TEXT PRIMARY KEY,
		value JSONB NOT NULL,
		updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
	);
	`
	_, err := db.ExecContext(ctx, query)
	return err
}

// seedSettings inserts a few default key -> JSON values (no-op if keys already exist).
func seedSettings(ctx context.Context, db *sqlx.DB) error {
	queries := []struct {
		Key   string
		Value string // JSON literal
	}{

		// Current Migration
		{"migrationVersion", `0`},

		// Organization settings
		{"org.siteName", `"Unbit"`},
		{"org.rootUrl", `"https://unbit.app"`},
		{"org.supportEmail", `"no-reply@unbit.app"`},
		{"org.logoUrl", `""`},
		{"org.description", `"Next-generation collaboration platform"`},
		{"org.timezone", `"UTC"`},

		// Authentication settings
		{"auth.allowRegistration", `true`},
		{"auth.requireEmailVerification", `true`},
		{"auth.passwordMinLength", `8`},
		{"auth.maxLoginAttempts", `5`},

		// UI settings
		{"ui.theme", `"light"`},
		{"ui.language", `"en-US"`},

		// Security settings
		{"security.enable2fa", `false`},
		{"security.sessionTimeoutMinutes", `30`},
		{"security.allowedIPRanges", `[]`},

		// Maintenance settings
		{"system.maintenanceMode", `false`},
		{"system.maintenanceMessage", `"System under maintenance. Please check back later."`},

		// Miscellaneous
		{"feature.inviteUsers", `true`},
	}

	tx, err := db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	stmt := `
	INSERT INTO settings (key, value)
	VALUES ($1, $2::jsonb)
	ON CONFLICT (key) DO NOTHING;
	`
	for _, q := range queries {
		if _, err := tx.ExecContext(ctx, stmt, q.Key, q.Value); err != nil {
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}
