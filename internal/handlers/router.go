package handlers

import (
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
)

type HandlerContext struct {
	DB  *sqlx.DB
	RDB *redis.Client
}

func (hc *HandlerContext) RegisterAPIRoutes(g *echo.Group) {
	g.POST("/auth/login", hc.LoginHandler)
	g.GET("/auth/profile", hc.ProfileHandler)
}
