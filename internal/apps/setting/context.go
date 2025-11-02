package setting

import (
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
)

type SettingContext struct {
	DB    *sqlx.DB
	Cache *redis.Client
}
