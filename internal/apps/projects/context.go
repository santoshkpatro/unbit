package projects

import (
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
)

type ProjectContext struct {
	DB    *sqlx.DB
	Cache *redis.Client
}
