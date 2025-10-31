package utils

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// Response is the standard shape for API responses.
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   interface{} `json:"error,omitempty"`
}

// OK sends a 200-success JSON response.
// If data is nil, an empty object is used.
func RespondOK(c echo.Context, data interface{}, message string) error {
	if data == nil {
		data = map[string]interface{}{}
	}

	res := Response{
		Success: true,
		Message: message,
		Data:    data,
	}

	return c.JSON(http.StatusOK, res)
}

// Fail sends an error JSON response.
// If status is 0, it defaults to 400.
// If err is nil, an empty object is returned.
func RespondFail(c echo.Context, status int, message string, err interface{}) error {
	if status == 0 {
		status = http.StatusBadRequest
	}
	if err == nil {
		err = map[string]interface{}{}
	}

	res := Response{
		Success: false,
		Message: message,
		Error:   err,
	}

	return c.JSON(status, res)
}
