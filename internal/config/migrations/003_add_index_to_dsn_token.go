package migrations

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
)

func init() {
	RegisterMigration(Migration{
		Version: 3,
		Up: func(ctx context.Context, tx *sqlx.Tx) error {
			_, err := tx.ExecContext(ctx, `
				CREATE UNIQUE INDEX IF NOT EXISTS idx_dsn_token_unique ON projects(dsn_token);
			`)
			if err != nil {
				return fmt.Errorf("failed to apply migration: %w", err)

			}
			return nil
		},
		Down: func(ctx context.Context, tx *sqlx.Tx) error {
			_, err := tx.ExecContext(ctx, `
				DROP INDEX IF EXISTS idx_dsn_token_unique;
			`)
			if err != nil {
				return fmt.Errorf("failed to revert migration version: %w", err)
			}
			return nil
		},
	})
}
