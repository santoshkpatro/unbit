package views

import (
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
)

type ViewContext struct {
	DB    *sqlx.DB
	Cache *redis.Client
}
