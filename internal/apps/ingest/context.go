package ingest

import (
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
)

type IngestContext struct {
	DB    *sqlx.DB
	Cache *redis.Client
}
