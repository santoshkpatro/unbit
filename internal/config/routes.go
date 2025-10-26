package config

import (
	"github.com/gorilla/sessions"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
	"github.com/santoshkpatro/unbit/internal/views"
)

func RegisterRoutes(e *echo.Echo, db *sqlx.DB, cache *redis.Client) {
	view := &views.ViewContext{
		DB:    db,
		Cache: cache,
	}

	api := e.Group("/api")
	api.Use(session.Middleware(sessions.NewCookieStore([]byte(Env.SecretKey))))

	api.POST("/auth/login", view.LoginUser)
	api.GET("/auth/profile", view.Profile)
}
