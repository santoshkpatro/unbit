package views

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (v *ViewContext) LoginUser(c echo.Context) error {

	return nil
}

func (v *ViewContext) AuthenticationStatus(c echo.Context) error {
	return c.String(http.StatusOK, "Working Good!")
}
