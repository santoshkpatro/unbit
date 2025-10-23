package injest

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (ic *InjestContext) Capture(c echo.Context) error {
	unbitAuth := c.Request().Header.Get("X-Unbit-Auth")
	fmt.Println("Auth", unbitAuth)

	return c.String(http.StatusOK, "Hello, World")
}
