package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (hc *HandlerContext) LoginHandler(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World")
}

func (hc *HandlerContext) ProfileHandler(c echo.Context) error {
	return c.String(http.StatusOK, "Profile")
}
