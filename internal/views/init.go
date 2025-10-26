package views

import (
	"errors"
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo-contrib/session"
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

// CheckAuthentication ensures the user is logged in.
// If the session is valid, returns (userID, nil).
// Otherwise, sends a 401 JSON error and returns ("", error).
func (v *ViewContext) CheckAuthentication(c echo.Context) (string, error) {
	sess, err := session.Get("session", c)
	if err != nil {
		// session middleware misconfigured or cookie invalid
		_ = v.RespondFail(c, http.StatusUnauthorized, "Unauthorized", "invalid session")
		return "", err
	}

	userID, ok := sess.Values["loggedInUser"].(string)
	if !ok || userID == "" {
		_ = v.RespondFail(c, http.StatusUnauthorized, "Unauthorized", "user not logged in")
		return "", errors.New("user not logged in")
	}

	return userID, nil
}
