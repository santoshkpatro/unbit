package handlers

import (
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
)

type HandlerContext struct {
	DB  *sqlx.DB
	RDB *redis.Client
}

type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   interface{} `json:"error,omitempty"`
}

// SuccessResponse sends a consistent success JSON response
func (hc *HandlerContext) SuccessResponse(c echo.Context, status int, message string, data interface{}) error {
	if message == "" {
		message = "success"
	}
	resp := APIResponse{
		Success: true,
		Message: message,
		Data:    data,
	}
	return c.JSON(status, resp)
}

// ErrorResponse sends a consistent error JSON response
func (hc *HandlerContext) ErrorResponse(c echo.Context, status int, message string, err error, details interface{}) error {
	if message == "" {
		message = http.StatusText(status)
	}

	// Log internal error safely
	if err != nil {
		c.Logger().Error(err)
	}

	resp := APIResponse{
		Success: false,
		Message: message,
		Error:   details,
	}

	return c.JSON(status, resp)
}

func (hc *HandlerContext) RegisterAPIRoutes(g *echo.Group) {
	// Auth Routes
	g.POST("/auth/login", hc.LoginHandler)
	g.GET("/auth/profile", hc.ProfileHandler)

	// Setting Routes
	g.GET("/setting/meta", hc.SettingMetaHandler)

	// Injest Routes
	g.POST("/injest/capture", hc.InjestCaptureHandler)
}
