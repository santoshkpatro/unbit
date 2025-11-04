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

func (v *ProjectContext) ProjectCreateView(c echo.Context) error {
	userID, _ := utils.CheckAuthentication(c)

	var newProject ProjectNew
	if err := c.Bind(&newProject); err != nil {
		return utils.RespondFail(c, http.StatusBadRequest, "Invalid request payload", err)
	}

	var newProjectId = utils.GenerateID("prj")
	_, err := v.DB.Exec(`
		INSERT INTO projects (id, name, description, dsn_token)
		VALUES ($1, $2, $3, $4)
	`, newProjectId, newProject.Name, newProject.Description, utils.GenerateID("dsn"))
	if err != nil {
		return utils.RespondFail(c, http.StatusInternalServerError, "Failed to create project", err)
	}

	_, err = v.DB.Exec(`
		INSERT INTO project_members (id, project_id, user_id, role)
		VALUES ($1, $2, $3, 'owner')
	`, utils.GenerateID("prm_"), newProjectId, userID)
	if err != nil {
		return utils.RespondFail(c, http.StatusInternalServerError, "Failed to add project member", err)
	}

	var createdProject Project
	err = v.DB.Get(&createdProject, `
		SELECT *
		FROM projects
		WHERE id = $1
	`, newProjectId)
	if err != nil {
		return utils.RespondFail(c, http.StatusInternalServerError, "Failed to fetch created project", err)
	}

	return utils.RespondOK(c, createdProject, "Project created successfully")
}
