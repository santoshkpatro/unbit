package issues

import (
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
)

type IssueContext struct {
	DB    *sqlx.DB
	Cache *redis.Client
}
