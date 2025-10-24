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
				ALTER TABLE users
				ADD COLUMN IF NOT EXISTS salt TEXT;
			`)
			if err != nil {
				return fmt.Errorf("creating users table: %w", err)
			}
			return nil
		},
		Down: func(ctx context.Context, tx *sqlx.Tx) error {
			_, err := tx.ExecContext(ctx, `
				ALTER TABLE users
				DROP COLUMN IF EXISTS salt;
			`)
			if err != nil {
				return fmt.Errorf("dropping users table: %w", err)
			}
			return nil
		},
	})
}
