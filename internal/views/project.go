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
