package events

import (
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/santoshkpatro/unbit/internal/utils"
)

func (v *EventContext) EventList(c echo.Context) error {
	userID, _ := utils.CheckAuthentication(c)
	query := `
		WITH latest_per_group AS (
			SELECT DISTINCT ON (e.group_id)
				g.id          AS group_id,
				g.status      AS group_status,
				g.event_count AS group_event_count,
				p.id          AS project_id,
				p.name        AS project_name,
				e.id          AS event_id,
				e.type        AS event_type,
				e.message     AS event_message,
				e.timestamp   AS event_timestamp,
				e.created_at  AS created_at,
				u.id          AS assignee_id,
				u.email       AS assignee_email,
				u.first_name  AS assignee_first_name
			FROM events e
			JOIN groups   g ON g.id = e.group_id
			JOIN projects p ON p.id = e.project_id
			LEFT JOIN users u ON u.id = g.assignee_id
			WHERE p.id IN (
				SELECT project_id FROM project_members WHERE user_id = $1
			)
			ORDER BY e.group_id, e.timestamp DESC, e.created_at DESC, e.id DESC
			)
			SELECT *
			FROM latest_per_group
			ORDER BY event_timestamp DESC, created_at DESC
			LIMIT 100;
	`

	var rows []eventRow
	err := v.DB.Select(&rows, query, userID)
	if err != nil {
		fmt.Println("Err", err)
		return utils.RespondFail(c, 500, "Failed to fetch results", nil)
	}
	events := make([]Event, 0, len(rows))
	for _, row := range rows {
		events = append(events, row.toEvent())
	}

	return utils.RespondOK(c, events, "")
}
