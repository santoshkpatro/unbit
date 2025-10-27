package views

import (
	"database/sql"
	"encoding/json"

	"github.com/labstack/echo/v4"
	"github.com/santoshkpatro/unbit/internal/utils"
)

type ProjectCreateData struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
}

type ProjectResponse struct {
	ID          string `json:"id" db:"id"`
	Name        string `json:"name" db:"name"`
	Description string `json:"description" db:"description"`
	DSNToken    string `json:"dsnToken" db:"dsn_token"`
	CreatedAt   string `json:"createdAt" db:"created_at"`
}

type User struct {
	ID        string `json:"id" db:"id"`
	Email     string `json:"email" db:"email"`
	FirstName string `json:"firstName" db:"first_name"`
	LastName  string `json:"lastName" db:"last_name"`
}

type Member struct {
	UserID string `json:"userId" db:"user_id"`
	Role   string `json:"role" db:"role"`
	User   User   `json:"user" db:"-"`
}

type ProjectDetailResponse struct {
	ID          string   `json:"id" db:"id"`
	Name        string   `json:"name" db:"name"`
	Description string   `json:"description" db:"description"`
	DSNToken    string   `json:"dsnToken" db:"dsn_token"`
	CreatedAt   string   `json:"createdAt" db:"created_at"`
	Members     []Member `json:"members" db:"-"`
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

	var newProject ProjectResponse
	err := v.DB.Get(&newProject, "SELECT id, name, description, dsn_token, created_at FROM projects WHERE id = $1", newProjectID)
	if err != nil {
		return v.RespondFail(c, 500, "Database error", err.Error())
	}

	return v.RespondOK(c, newProject, "New project created")
}

func (v *ViewContext) ProjectListView(c echo.Context) error {
	userId, _ := v.CheckAuthentication(c)

	var projects []ProjectResponse
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

	var userRole string
	err := v.DB.Get(&userRole, "SELECT role FROM project_members WHERE project_id = $1 AND user_id = $2", projectId, userId)
	if err != nil {
		if err == sql.ErrNoRows {
			return v.RespondFail(c, 403, "Forbidden", "user is not a member of the project")
		}
		return v.RespondFail(c, 500, "Database error", err.Error())
	}

	// members: array of { userId, role, user: { id, email, firstName, lastName } }
	const q = `
	SELECT
		p.id,
		p.name,
		p.description,
		p.dsn_token,
		p.created_at,
		COALESCE(
		  json_agg(
			json_build_object(
			  'userId', pm.user_id,
			  'role', pm.role,
			  'user', json_build_object(
				'id', u.id,
				'email', u.email,
				'firstName', u.first_name,
				'lastName', u.last_name
			  )
			)
		  ) FILTER (WHERE pm.user_id IS NOT NULL),
		  '[]'
		) AS members
	FROM projects p
	JOIN project_members pm ON p.id = pm.project_id
	LEFT JOIN users u ON u.id = pm.user_id
	WHERE p.id = $1
	GROUP BY p.id
	`

	// intermediate row to receive the DB result
	type projectRow struct {
		ID          string          `db:"id"`
		Name        string          `db:"name"`
		Description string          `db:"description"`
		DSNToken    string          `db:"dsn_token"`
		CreatedAt   string          `db:"created_at"`
		MembersJSON json.RawMessage `db:"members"`
	}

	var row projectRow
	err = v.DB.Get(&row, q, projectId)
	if err != nil {
		if err == sql.ErrNoRows {
			return v.RespondFail(c, 404, "Not found", "project not found")
		}
		return v.RespondFail(c, 500, "Database error", err.Error())
	}

	var members []Member
	if err := json.Unmarshal(row.MembersJSON, &members); err != nil {
		return v.RespondFail(c, 500, "JSON unmarshal error", err.Error())
	}

	resp := ProjectDetailResponse{
		ID:          row.ID,
		Name:        row.Name,
		Description: row.Description,
		DSNToken:    row.DSNToken,
		CreatedAt:   row.CreatedAt,
		Members:     members,
	}

	return v.RespondOK(c, resp, "Project details fetched")
}
