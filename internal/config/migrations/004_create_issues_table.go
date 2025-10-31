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
				CREATE TABLE IF NOT EXISTS issues (
					id TEXT PRIMARY KEY,
					project_id TEXT NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
					summary TEXT NOT NULL,
					fingerprint TEXT NOT NULL,
					assignee_id TEXT REFERENCES users(id) ON DELETE SET NULL,
					status TEXT NOT NULL DEFAULT 'unresolved',
					last_seen_at TIMESTAMPTZ NOT NULL,
					event_count BIGINT DEFAULT 1,
					created_at TIMESTAMPTZ DEFAULT NOW(),
					updated_at TIMESTAMPTZ DEFAULT NOW()
				);

				CREATE INDEX IF NOT EXISTS idx_issues_project_id ON issues(project_id);
				CREATE INDEX IF NOT EXISTS idx_issues_assignee_id ON issues(assignee_id);
				CREATE UNIQUE INDEX IF NOT EXISTS idx_unique_issue_fingerprint_per_project ON issues(project_id, fingerprint);
				CREATE INDEX IF NOT EXISTS idx_issues_status ON issues(status);
			`)
			if err != nil {
				return fmt.Errorf("failed to apply migration: %w", err)

			}
			return nil
		},
		Down: func(ctx context.Context, tx *sqlx.Tx) error {
			_, err := tx.ExecContext(ctx, `
				DROP TABLE IF EXISTS issues;
			`)
			if err != nil {
				return fmt.Errorf("failed to revert migration version: %w", err)
			}
			return nil
		},
	})
}
