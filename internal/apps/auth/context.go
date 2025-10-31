package auth

import (
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
)

type AuthContext struct {
	DB    *sqlx.DB
	Cache *redis.Client
}
