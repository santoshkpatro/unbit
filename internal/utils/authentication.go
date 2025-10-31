package utils

import (
	"net/http"

	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

func CheckAuthentication(c echo.Context) (string, error) {
	sess, err := session.Get("session", c)
	if err != nil {
		RespondFail(c, http.StatusInternalServerError, "Session error", err.Error())
		return "", err
	}

	userID, ok := sess.Values["loggedInUser"]
	if !ok || userID == nil {
		RespondFail(c, http.StatusUnauthorized, "User not authenticated", nil)
		return "", echo.NewHTTPError(http.StatusUnauthorized, "User not authenticated")
	}

	id, ok := userID.(string)
	if !ok {
		RespondFail(c, http.StatusUnauthorized, "Invalid session data", nil)
		return "", echo.NewHTTPError(http.StatusUnauthorized, "Invalid session data")
	}

	return id, nil
}
