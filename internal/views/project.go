package views

import (
	"github.com/labstack/echo/v4"
	"github.com/santoshkpatro/unbit/internal/models"
	"github.com/santoshkpatro/unbit/internal/utils"
)

type ProjectCreateData struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
}

func (v *ViewContext) ProjectCreateView(c echo.Context) error {
	userId, _ := v.CheckAuthentication(c)

	var data ProjectCreateData
	if err := c.Bind(&data); err != nil {
		return v.RespondFail(c, 400, "Invalid request data", err.Error())
	}

	if err := c.Validate(&data); err != nil {
		return v.RespondFail(c, 400, "Validation failed", err.Error())
	}

	tx := v.DB.MustBegin()
	newProjectID := utils.GenerateID("proj")
	tx.MustExec("INSERT INTO projects (id, name, description, dsn_token) VALUES ($1, $2, $3, $4)", newProjectID, data.Name, data.Description, utils.GenerateID("dsn"))
	tx.MustExec("INSERT INTO project_members (id, project_id, user_id, role) VALUES ($1, $2, $3, $4)", utils.GenerateID("prme"), newProjectID, userId, "owner")
	if err := tx.Commit(); err != nil {
		return v.RespondFail(c, 500, "Database error", err.Error())
	}

	var newProject models.Project
	err := v.DB.Get(&newProject, "SELECT id, name, description, dsn_token, created_at FROM projects WHERE id=$1", newProjectID)
	if err != nil {
		return v.RespondFail(c, 500, "Database error", err.Error())
	}

	return v.RespondOK(c, newProject, "New project created")
}

func (v *ViewContext) ProjectListView(c echo.Context) error {
	userId, _ := v.CheckAuthentication(c)

	var projects []models.Project
	err := v.DB.Select(&projects, `
		SELECT p.id, p.name, p.description, p.dsn_token, p.created_at
		FROM projects p
		JOIN project_members pm ON p.id = pm.project_id
		WHERE pm.user_id = $1
	`, userId)
	if err != nil {
		return v.RespondFail(c, 500, "Database error", err.Error())
	}

	return v.RespondOK(c, projects, "User projects fetched")
}

func (v *ViewContext) ProjectDetailView(c echo.Context) error {
	userId, _ := v.CheckAuthentication(c)
	projectId := c.Param("projectId")

	var project models.Project
	err := v.DB.Get(&project, `
		SELECT p.id, p.name, p.description, p.dsn_token, p.created_at
		FROM projects p
		JOIN project_members pm ON p.id = pm.project_id
		WHERE pm.user_id = $1 AND p.id = $2
	`, userId, projectId)
	if err != nil {
		return v.RespondFail(c, 404, "Project not found or access denied", err.Error())
	}

	var members []models.ProjectMember
	err = v.DB.Select(&members, `
		SELECT 
			pm.id, 
			pm.project_id, 
			pm.user_id, 
			pm.role, 
			pm.joined_at, 
			pm.created_at,
			json_build_object(
				'id', u.id,
				'email', u.email,
				'first_name', u.first_name,
				'last_name', u.last_name,
				'is_admin', u.is_admin
			) AS user
		FROM project_members pm
		JOIN users u ON u.id = pm.user_id
		WHERE pm.project_id = $1
		ORDER BY pm.joined_at ASC
	`, projectId)
	if err != nil {
		return v.RespondFail(c, 500, "Database error", err.Error())
	}
	project.Members = members

	return v.RespondOK(c, project, "Project details fetched")
}
