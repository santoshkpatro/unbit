package views

import (
	"github.com/labstack/echo/v4"
)

func (v *ViewContext) LoginUser(c echo.Context) error {

	return nil
}

func (v *ViewContext) AuthenticationStatus(c echo.Context) error {
	return v.RespondOK(c, nil, nil)
}

func (v *ViewContext) LogoutUser(c echo.Context) error {
	return v.RespondFail(c, nil, nil, nil)
}
