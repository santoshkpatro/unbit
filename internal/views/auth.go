package views

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type loginData struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

func (v *ViewContext) LoginUser(c echo.Context) error {
	var data loginData
	if err := c.Bind(&data); err != nil {
		return v.RespondFail(c, http.StatusBadRequest, "Invalid request data", err.Error())
	}

	if err := c.Validate(&data); err != nil {
		return v.RespondFail(c, http.StatusBadRequest, "Validation failed", err.Error())
	}

	var userExists bool
	query := "SELECT EXISTS(SELECT 1 FROM users WHERE email=$1)"
	err := v.DB.Get(&userExists, query, data.Email)
	if err != nil {
		return v.RespondFail(c, http.StatusInternalServerError, "Database error", err.Error())
	}
	if !userExists {
		return v.RespondFail(c, http.StatusUnauthorized, "Invalid email or password", nil)
	}

	return c.String(200, "Ok")
}
