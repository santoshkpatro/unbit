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
					fingerprint TEXT NOT NULL,
					assignee_id TEXT REFERENCES users(id),
					status TEXT NOT NULL DEFAULT 'unresolved',
					event_count BIGINT DEFAULT 0,
					created_at TIMESTAMPTZ DEFAULT NOW(),
					updated_at TIMESTAMPTZ DEFAULT NOW()
				);
				CREATE INDEX IF NOT EXISTS idx_issues_project_id ON issues(project_id);
				CREATE INDEX IF NOT EXISTS idx_issues_assignee_id ON issues(assignee_id);
				CREATE INDEX IF NOT EXISTS idx_issues_status ON issues(status);
				CREATE UNIQUE INDEX IF NOT EXISTS idx_issues_project_fingerprint ON issues(project_id, fingerprint);

				CREATE TABLE IF NOT EXISTS events (
					id TEXT PRIMARY KEY,
					project_id TEXT NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
					issue_id TEXT NOT NULL REFERENCES issues(id) ON DELETE CASCADE,
					timestamp TIMESTAMPTZ NOT NULL,
					properties JSONB,
					event_type TEXT,
					created_at TIMESTAMPTZ DEFAULT NOW(),
					updated_at TIMESTAMPTZ DEFAULT NOW()
				);
				CREATE INDEX IF NOT EXISTS idx_events_issue_id ON events(issue_id);
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
				DROP TABLE IF EXISTS issues;
			`)
			if err != nil {
				return fmt.Errorf("failed to revert migration version: %w", err)
			}
			return nil
		},
	})
}
