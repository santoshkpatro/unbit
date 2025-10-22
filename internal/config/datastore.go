package config

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
)

func NewPostgresConnection(ctx context.Context) (*sqlx.DB, error) {
	db, err := sqlx.Connect("postgres", Env.PgUrl)
	if err != nil {
		log.Fatalln(err)
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping postgres: %w", err)
	}

	return db, nil
}

func NewRedisConnection(ctx context.Context) (*redis.Client, error) {
	opt, err := redis.ParseURL(Env.RedisUrl)
	if err != nil {
		panic(err)
	}

	rdb := redis.NewClient(opt)
	if err := rdb.Ping(ctx).Err(); err != nil {
		rdb.Close()
		return nil, fmt.Errorf("failed to connect to redis: %w", err)
	}

	return rdb, nil
}
