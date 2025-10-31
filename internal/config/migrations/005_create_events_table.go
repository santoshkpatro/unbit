package migrations

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
)

func init() {
	RegisterMigration(Migration{
		Version: 5,
		Up: func(ctx context.Context, tx *sqlx.Tx) error {
			_, err := tx.ExecContext(ctx, `
				CREATE TABLE IF NOT EXISTS events (
					id TEXT PRIMARY KEY,
					issue_id TEXT NOT NULL REFERENCES issues(id) ON DELETE CASCADE,
					message TEXT NOT NULL,
					level TEXT NOT NULL,
					timestamp TIMESTAMPTZ NOT NULL,
					stacktrace JSONB NOT NULL,
					created_at TIMESTAMPTZ DEFAULT NOW(),
					updated_at TIMESTAMPTZ DEFAULT NOW()
				);
				CREATE INDEX IF NOT EXISTS idx_events_issue_id ON events(issue_id);
			`)
			if err != nil {
				return fmt.Errorf("failed to apply migration: %w", err)

			}
			return nil
		},
		Down: func(ctx context.Context, tx *sqlx.Tx) error {
			_, err := tx.ExecContext(ctx, `
				DROP TABLE IF EXISTS events;
			`)
			if err != nil {
				return fmt.Errorf("failed to revert migration version: %w", err)
			}
			return nil
		},
	})
}
