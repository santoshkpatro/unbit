package issues

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/santoshkpatro/unbit/internal/utils"
)

func (v *IssueContext) IssueList(c echo.Context) error {
	userID, _ := utils.CheckAuthentication(c)

	query := `
		SELECT 
			i.id, i.project_id, i.summary, i.assignee_id, i.status, 
			i.last_seen_at, i.event_count, i.created_at,
			p.name AS project_name,
			u.id AS assignee_id, u.email AS assignee_email, u.first_name AS assignee_first_name
		FROM issues i
		INNER JOIN projects p ON i.project_id = p.id
		LEFT JOIN users u ON i.assignee_id = u.id
		WHERE i.project_id IN (
			SELECT project_id FROM project_members WHERE user_id = $1
		)
		ORDER BY i.last_seen_at DESC
		LIMIT 100
	`
	var issueRows []IssueRow
	err := v.DB.Select(&issueRows, query, userID)
	if err != nil {
		fmt.Println("Err", err)
		return utils.RespondFail(c, http.StatusInternalServerError, "Failed to fetch results", nil)
	}

	issues := make([]Issue, 0, len(issueRows))
	for _, row := range issueRows {
		issue := Issue{
			ID: row.ID,
			Project: &Project{
				ID:   row.ProjectId,
				Name: row.ProjectName,
			},
			Summary:    row.Summary,
			AssigneeID: row.AssigneeID,
			Status:     row.Status,
			LastSeenAt: row.LastSeenAt,
			EventCount: row.EventCount,
			CreatedAt:  row.CreatedAt,
		}

		// Only populate assignee if assignee_id is not null
		if row.AssigneeID.Valid {
			issue.Assignee = &Assignee{
				ID:        row.AssigneeID.String,
				Email:     row.AssigneeEmail.String,
				FirstName: row.AssigneeFirstName.String,
			}
		}

		issues = append(issues, issue)
	}

	return utils.RespondOK(c, issues, "")
}
