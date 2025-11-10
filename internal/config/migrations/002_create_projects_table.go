package migrations

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
)

func init() {
	RegisterMigration(Migration{
		Version: 2,
		Up: func(ctx context.Context, tx *sqlx.Tx) error {
			_, err := tx.ExecContext(ctx, `
				CREATE TABLE IF NOT EXISTS projects (
					id TEXT PRIMARY KEY,
					name TEXT NOT NULL UNIQUE,
					description TEXT,
					dsn_token TEXT,
					total_events BIGINT DEFAULT 0,
					created_at TIMESTAMPTZ DEFAULT NOW(),
					updated_at TIMESTAMPTZ DEFAULT NOW()
				);

				CREATE TABLE IF NOT EXISTS project_members (
					id TEXT PRIMARY KEY,
					project_id TEXT NOT NULL,
					user_id TEXT NOT NULL,
					role TEXT NOT NULL DEFAULT 'member',
					joined_at TIMESTAMPTZ,
					created_at TIMESTAMPTZ DEFAULT NOW(),
					updated_at TIMESTAMPTZ DEFAULT NOW(),
					FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE,
					FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
				);
				CREATE UNIQUE INDEX IF NOT EXISTS idx_unique_project_member ON project_members(project_id, user_id);
			`)
			if err != nil {
				return fmt.Errorf("failed to apply migration: %w", err)

			}
			return nil
		},
		Down: func(ctx context.Context, tx *sqlx.Tx) error {
			_, err := tx.ExecContext(ctx, `
				DROP TABLE IF EXISTS project_members;
				DROP TABLE IF EXISTS projects;
			`)
			if err != nil {
				return fmt.Errorf("failed to revert migration version: %w", err)
			}
			return nil
		},
	})
}
