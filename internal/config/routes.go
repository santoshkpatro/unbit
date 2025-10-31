package config

import (
	"github.com/gorilla/sessions"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
	"github.com/santoshkpatro/unbit/internal/apps/auth"
	"github.com/santoshkpatro/unbit/internal/apps/ingest"
	"github.com/santoshkpatro/unbit/internal/apps/projects"
)

func RegisterRoutes(e *echo.Echo, db *sqlx.DB, cache *redis.Client) {
	api := e.Group("/api")
	api.Use(session.Middleware(sessions.NewCookieStore([]byte(Env.SecretKey))))

	// Ingest routes
	ingestContext := &ingest.IngestContext{
		DB:    db,
		Cache: cache,
	}
	api.POST("/ingest/event", ingestContext.NewEvent)

	// Auth routes
	authContext := &auth.AuthContext{
		DB:    db,
		Cache: cache,
	}
	api.POST("/auth/login", authContext.LoginUser)
	api.GET("/auth/profile", authContext.Profile)

	// Project routes
	projectContext := &projects.ProjectContext{
		DB:    db,
		Cache: cache,
	}
	api.GET("/projects", projectContext.ProjectListView)
}
