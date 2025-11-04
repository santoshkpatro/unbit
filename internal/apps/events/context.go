package events

import (
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
)

type EventContext struct {
	DB    *sqlx.DB
	Cache *redis.Client
}
