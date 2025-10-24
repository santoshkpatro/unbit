package migrations

import (
	"context"
	"fmt"
	"sort"

	"github.com/jmoiron/sqlx"
)

type Migration struct {
	Version int
	Up      func(ctx context.Context, tx *sqlx.Tx) error
	Down    func(ctx context.Context, tx *sqlx.Tx) error
}

var migrationRegistry = make(map[int]Migration)

func RegisterMigration(m Migration) {
	if _, exists := migrationRegistry[m.Version]; exists {
		panicf("migration version %d already registered", m.Version)
	}

	migrationRegistry[m.Version] = m
}

func MigrationList() []Migration {
	vers := make([]int, 0, len(migrationRegistry))
	for v := range migrationRegistry {
		vers = append(vers, v)
	}
	sort.Ints(vers)

	out := make([]Migration, 0, len(vers))
	for _, v := range vers {
		out = append(out, migrationRegistry[v])
	}
	return out
}

func panicf(format string, args ...interface{}) {
	// We don't import fmt to keep migrations files small; this still prints a useful error
	panic(fmt.Sprintf(format, args...))
}
