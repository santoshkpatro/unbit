package views

import (
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
)

type ViewContext struct {
	DB    *sqlx.DB
	Cache *redis.Client
}

// RespondOK sends a 200 JSON response. Pass nil for data to get an empty object.
// Pass "" for message if you don't want a message.
func (v *ViewContext) RespondOK(c echo.Context, data interface{}, message string) error {
	if data == nil {
		data = map[string]interface{}{}
	}

	response := map[string]interface{}{
		"success": true,
		"message": message,
		"data":    data,
	}
	return c.JSON(http.StatusOK, response)
}

// RespondFail sends an error JSON response. Use status=0 to default to 400.
// Pass nil for err to get an empty object. Pass "" for message if not needed.
func (v *ViewContext) RespondFail(c echo.Context, status int, message string, err interface{}) error {
	if status == 0 {
		status = http.StatusBadRequest
	}
	if err == nil {
		err = map[string]interface{}{}
	}

	response := map[string]interface{}{
		"success": false,
		"message": message,
		"error":   err,
	}
	return c.JSON(status, response)
}
