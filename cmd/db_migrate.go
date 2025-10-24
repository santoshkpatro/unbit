package cmd

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/jmoiron/sqlx"
	"github.com/santoshkpatro/unbit/internal/config"
	"github.com/santoshkpatro/unbit/internal/config/migrations"
	"github.com/spf13/cobra"
)

var dbMigrateCmd = &cobra.Command{
	Use:   "db_migrate",
	Short: "DB Migrate",
	RunE: func(cmd *cobra.Command, args []string) error {
		return dbUpgarde()
	},
}

func dbUpgarde() error {
	ctx := context.Background()
	db, err := config.NewPostgresConnection(ctx)
	if err != nil {
		log.Fatalf("❌ failed to connect to postgres: %v", err)
	}
	defer db.Close()

	currentVersion, err := readCurrentMigrationVersion(ctx, db)
	if err != nil {
		return fmt.Errorf("reading current migration version: %w", err)
	}

	all := migrations.MigrationList()

	for _, m := range all {
		if m.Version <= currentVersion {
			continue
		}

		log.Printf("⬆️  applying migration version %d", m.Version)
		tx, err := db.BeginTxx(ctx, nil)
		if err != nil {
			return fmt.Errorf("starting transaction for migration %d: %w", m.Version, err)
		}

		if m.Up != nil {
			if err := m.Up(ctx, tx); err != nil {
				tx.Rollback()
				return fmt.Errorf("migration %d up failed: %w", m.Version, err)
			}
		}

		_, err = tx.ExecContext(ctx, `
			INSERT INTO settings(key, value)
			VALUES ('migrationVersion', $1)
			ON CONFLICT (key) DO UPDATE SET value = EXCLUDED.value;
		`, strconv.Itoa(m.Version))
		if err != nil {
			_ = tx.Rollback()
			return fmt.Errorf("updating migrationVersion for %d: %w", m.Version, err)
		}

		if err := tx.Commit(); err != nil {
			return fmt.Errorf("commit migration %d: %w", m.Version, err)
		}

		log.Printf("✅ migration %04d applied", m.Version)
		currentVersion = m.Version
	}
	log.Printf("All migrations up to %d applied.", currentVersion)
	return nil
}

func readCurrentMigrationVersion(ctx context.Context, db *sqlx.DB) (int, error) {
	var vstr string
	err := db.QueryRowContext(ctx, `SELECT value FROM settings WHERE key = 'migrationVersion'`).Scan(&vstr)
	if err != nil {
		return 0, err
	}
	if vstr == "" {
		return 0, nil
	}
	v, err := strconv.Atoi(vstr)
	if err != nil {
		return 0, fmt.Errorf("bad migrationVersion value %q: %w", vstr, err)
	}
	return v, nil
}
