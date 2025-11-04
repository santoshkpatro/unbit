package migrations

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
)

func init() {
	RegisterMigration(Migration{
		Version: 4,
		Up: func(ctx context.Context, tx *sqlx.Tx) error {
			_, err := tx.ExecContext(ctx, `
				CREATE TABLE IF NOT EXISTS groups (
					id TEXT PRIMARY KEY,
					project_id TEXT NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
					fingerprint TEXT NOT NULL,
					assignee_id TEXT REFERENCES users(id),
					status TEXT NOT NULL DEFAULT 'unresolved',
					event_count BIGINT DEFAULT 0,
					created_at TIMESTAMPTZ DEFAULT NOW(),
					updated_at TIMESTAMPTZ DEFAULT NOW()
				);
				CREATE INDEX IF NOT EXISTS idx_groups_project_id ON groups(project_id);
				CREATE UNIQUE INDEX IF NOT EXISTS idx_groups_project_fingerprint ON groups(project_id, fingerprint);

				CREATE TABLE IF NOT EXISTS events (
					id TEXT PRIMARY KEY,
					project_id TEXT NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
					group_id TEXT NOT NULL REFERENCES groups(id) ON DELETE CASCADE,
					type TEXT NOT NULL,
					message TEXT NOT NULL,
					level TEXT NOT NULL,
					timestamp TIMESTAMPTZ NOT NULL,
					stacktrace JSONB,
					properties JSONB,
					created_at TIMESTAMPTZ DEFAULT NOW(),
					updated_at TIMESTAMPTZ DEFAULT NOW()
				);
				CREATE INDEX IF NOT EXISTS idx_events_group_id ON events(group_id);
				CREATE INDEX IF NOT EXISTS idx_events_project_id ON events(project_id);
			`)
			if err != nil {
				return fmt.Errorf("failed to apply migration: %w", err)

			}
			return nil
		},
		Down: func(ctx context.Context, tx *sqlx.Tx) error {
			_, err := tx.ExecContext(ctx, `
				DROP TABLE IF EXISTS events;
				DROP TABLE IF EXISTS groups;
			`)
			if err != nil {
				return fmt.Errorf("failed to revert migration version: %w", err)
			}
			return nil
		},
	})
}
