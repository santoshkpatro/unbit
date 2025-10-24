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
		// Organization settings
		{"org.site_name", `"Unbit"`},
		{"org.root_url", `"https://unbit.app"`},
		{"org.support_email", `"no-reply@unbit.app"`},
		{"org.logo_url", `""`},
		{"org.description", `"Next-generation collaboration platform"`},
		{"org.timezone", `"UTC"`},

		// Authentication settings
		{"auth.allow_registration", `true`},
		{"auth.require_email_verification", `true`},
		{"auth.password_min_length", `8`},
		{"auth.max_login_attempts", `5`},

		// UI settings
		{"ui.theme", `"light"`},
		{"ui.language", `"en-US"`},

		// Security settings
		{"security.enable_2fa", `false`},
		{"security.session_timeout_minutes", `30`},
		{"security.allowed_ip_ranges", `[]`},

		// Maintenance settings
		{"system.maintenance_mode", `false`},
		{"system.maintenance_message", `"System under maintenance. Please check back later."`},

		// Miscellaneous
		{"feature.invite_users", `true`},
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
