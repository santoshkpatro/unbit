package auth

import (
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/santoshkpatro/unbit/internal/utils"
)

func (v *AuthContext) LoginUser(c echo.Context) error {
	var data loginData
	if err := c.Bind(&data); err != nil {
		return utils.RespondFail(c, http.StatusBadRequest, "Invalid request data", err.Error())
	}

	if err := c.Validate(&data); err != nil {
		return utils.RespondFail(c, http.StatusBadRequest, "Validation failed", err.Error())
	}

	var userExists bool
	query := "SELECT EXISTS(SELECT 1 FROM users WHERE email=$1)"
	err := v.DB.Get(&userExists, query, data.Email)
	if err != nil {
		return utils.RespondFail(c, http.StatusInternalServerError, "Database error", err.Error())
	}
	if !userExists {
		return utils.RespondFail(c, http.StatusUnauthorized, "Invalid email or password", nil)
	}

	var user User
	err = v.DB.Get(&user, "SELECT * FROM users WHERE email = $1", data.Email)
	if err != nil {
		return utils.RespondFail(c, http.StatusInternalServerError, "Database error", err.Error())
	}

	sess, err := session.Get("session", c)
	if err != nil {
		return utils.RespondFail(c, http.StatusInternalServerError, "Session error", err.Error())
	}
	sess.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7,
		HttpOnly: true,
	}
	sess.Values["loggedInUser"] = user.ID
	if err := sess.Save(c.Request(), c.Response()); err != nil {
		return utils.RespondFail(c, http.StatusInternalServerError, "Failed to save session", err.Error())
	}

	return utils.RespondOK(c, user, "Login successful")
}

func (v *AuthContext) Profile(c echo.Context) error {
	userID, _ := utils.CheckAuthentication(c)
	var user User
	err := v.DB.Get(&user, "SELECT * FROM users WHERE id = $1", userID)
	if err != nil {
		return utils.RespondFail(c, http.StatusInternalServerError, "Database error", err.Error())
	}

	return utils.RespondOK(c, user, "User profile fetched successfully")
}
