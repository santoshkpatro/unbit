package projects

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/santoshkpatro/unbit/internal/utils"
)

func (v *ProjectContext) ProjectListView(c echo.Context) error {
	userID, _ := utils.CheckAuthentication(c)
	var projects []Project
	err := v.DB.Select(&projects, `
		SELECT p.*
		FROM projects p
		JOIN project_members pm ON p.id = pm.project_id
		WHERE pm.user_id = $1
	`, userID)
	if err != nil {
		return utils.RespondFail(c, http.StatusInternalServerError, "Failed to fetch projects", err)
	}

	return utils.RespondOK(c, projects, "")
}
