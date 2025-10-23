package injest

import (
	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
)

type InjestContext struct {
	RDB *redis.Client
}

func (ic *InjestContext) RegisterInjestRoutes(g *echo.Group) {
	g.POST("/capture", ic.Capture)
}
