package events

import (
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/santoshkpatro/unbit/internal/utils"
)

func (v *EventContext) EventIssues(c echo.Context) error {
	userID, _ := utils.CheckAuthentication(c)
	query := `
		WITH
			issues AS (
				SELECT DISTINCT
					ON (e.group_id) e.id AS event_id,
					e.group_id,
					e.message,
					e.type,
					e.level,
					g.status,
					g.assignee_id,
					e.project_id,
					g.event_count,
					p.name AS project_name,
					e.timestamp AS event_at
				FROM
					events e
					JOIN groups g ON g.id = e.group_id
					JOIN projects p ON p.id = e.project_id
				ORDER BY
					e.group_id,
					e.timestamp DESC,
					e.id DESC
			),
			date_spine_7d AS (
				SELECT
					generate_series(CURRENT_DATE - interval '6 day', CURRENT_DATE, interval '1 day')::date AS DAY
			),
			daily_event_counts AS (
				SELECT
					e.group_id,
					e.timestamp::date AS DAY,
					COUNT(*) AS event_count
				FROM
					events e
				WHERE
					e.timestamp::date BETWEEN CURRENT_DATE - interval '6 day' AND CURRENT_DATE
				GROUP BY
					e.group_id,
					e.timestamp::date
			)
		SELECT
			i.event_id,
			i.group_id,
			i.message,
			i.type,
			i.level,
			i.status,
			i.assignee_id,
			i.event_count,
			u.first_name,
			u.email,
			u.id AS user_id,
			i.project_id,
			i.project_name,
			(
				SELECT
					jsonb_agg(
						jsonb_build_object('day', d.day, 'count', COALESCE(dec.event_count, 0))
						ORDER BY
							d.day DESC
					)
				FROM
					date_spine_7d d
					LEFT JOIN daily_event_counts dec ON dec.group_id = i.group_id
					AND dec.day = d.day
			) AS last_7_days
		FROM
			issues i
			LEFT JOIN users u ON u.id = i.assignee_id
		WHERE
			i.project_id IN (
				SELECT
					project_id
				FROM
					project_members
				WHERE
					user_id = $1
			)
		ORDER BY
			i.event_at DESC;
	`
	var rows []issueRow
	err := v.DB.Select(&rows, query, userID)
	if err != nil {
		fmt.Println("Err", err)
		return utils.RespondFail(c, 500, "Failed to fetch results", nil)
	}
	issues := make([]Issue, 0, len(rows))
	for _, row := range rows {
		issues = append(issues, row.ToIssue())
	}

	return utils.RespondOK(c, issues, "")
}
