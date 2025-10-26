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

func (v *ViewContext) RespondOK(c echo.Context, data interface{}, msg *string) error {
	// if data is nil, make it an empty object instead of null
	if data == nil {
		data = map[string]interface{}{}
	}

	var message string
	if msg != nil {
		message = *msg
	}
	response := map[string]interface{}{
		"success": true,
		"message": message,
		"data":    data,
	}
	return c.JSON(http.StatusOK, response)
}

func (v *ViewContext) RespondFail(c echo.Context, code *int, msg *string, err interface{}) error {
	status := http.StatusBadRequest
	if code != nil {
		status = *code
	}

	message := ""
	if msg != nil {
		message = *msg
	}

	var errVal interface{}
	if err != nil {
		errVal = err
	} else {
		errVal = map[string]interface{}{}
	}

	response := map[string]interface{}{
		"success": false,
		"message": message,
		"error":   errVal,
	}

	return c.JSON(status, response)
}
