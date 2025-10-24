package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (hc *HandlerContext) SettingMetaHandler(c echo.Context) error {
	var settingMeta struct {
		OrgName         string `db:"org_name" json:"orgName"`
		OrgURL          string `db:"org_url" json:"orgUrl"`
		SupportEmail    string `db:"support_email" json:"supportEmail"`
		IsMaintenanceOn bool   `db:"is_maintenance_on" json:"isMaintenanceOn"`
	}
	err := hc.DB.Get(&settingMeta, "SELECT org_name, org_url, support_email, is_maintenance_on FROM setting LIMIT 1")
	if err != nil {
		return hc.ErrorResponse(c, http.StatusInternalServerError, "failed to fetch settings", err, nil)
	}

	return hc.SuccessResponse(c, http.StatusOK, "fetched settings successfully", settingMeta)
}
