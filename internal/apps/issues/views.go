package issues

import (
	"fmt"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/santoshkpatro/unbit/internal/utils"
)

func (v *IssueContext) RecentIssueListView(c echo.Context) error {
	userID, _ := utils.CheckAuthentication(c)
	var params []interface{}
	params = append(params, userID) // $1

	var extraWhere []string

	// optional param
	projectID := c.QueryParam("project_id")
	if projectID != "" {
		extraWhere = append(extraWhere, fmt.Sprintf("p.id = $%d", len(params)+1))
		params = append(params, projectID)
	}

	where := ""
	if len(extraWhere) > 0 {
		where = " AND (" + strings.Join(extraWhere, " AND ") + ")"
	}

	query := fmt.Sprintf(`
		WITH
			days AS (
				SELECT
					generate_series(CURRENT_DATE - interval '13 days', CURRENT_DATE, interval '1 day')::date AS DAY
			),
			recent_issues AS (
				SELECT DISTINCT
					ON (e.issue_id) i.id AS issue_id,
					e.id AS event_id,
					e.timestamp,
					i.event_count,
					i.assignee_id,
					u.email AS assignee_email,
					concat_ws(' ', u.first_name, u.last_name) AS assignee_name,
					p.id AS project_id,
					p.name AS project_name,
					e.properties ->> 'message' AS message,
					e.properties ->> 'type' AS type,
					e.properties ->> 'level' AS level,
					e.properties -> 'stacktrace' -> 0 AS first_stack_trace,
					floor(extract(epoch FROM (e.timestamp - i.created_at)))::int AS age
				FROM
					issues i
					JOIN events e ON i.id = e.issue_id
					JOIN projects p ON p.id = e.project_id
					LEFT JOIN users u ON i.assignee_id = u.id
				WHERE
					e.event_type = 'issues'
					AND p.id IN (
						SELECT
							project_id
						FROM
							project_members
						WHERE
							user_id = $1
					)
					%s
				ORDER BY
					e.issue_id,
					e.timestamp DESC
			),
			issue_events AS (
				SELECT
					e.id,
					e.issue_id,
					e.timestamp::date AS date
				FROM
					events e
					JOIN recent_issues ri ON ri.issue_id = e.issue_id
			),
			daily_issue_counts AS (
				SELECT
					d.day,
					ri.issue_id,
					count(ie.id) AS event_count
				FROM
					days d
					CROSS JOIN recent_issues ri
					LEFT JOIN issue_events ie ON ie.date = d.day
					AND ie.issue_id = ri.issue_id
				GROUP BY
					d.day,
					ri.issue_id
			)
		SELECT
			ri.issue_id AS id,
			ri.event_id,
			ri.event_count,
			ri.timestamp,
			ri.message,
			ri.level,
			ri.type,
			ri.first_stack_trace,
			ri.assignee_id,
			ri.assignee_name,
			ri.assignee_email,
			ri.project_id,
			ri.project_name,
			ri.age,
			json_agg(json_build_object('date', dic.day, 'eventCount', dic.event_count) ORDER BY dic.day DESC) AS issue_count_report
		FROM
			recent_issues ri
			JOIN daily_issue_counts dic ON ri.issue_id = dic.issue_id
		GROUP BY
			ri.issue_id,
			ri.event_id,
			ri.event_count,
			ri.timestamp,
			ri.message,
			ri.level,
			ri.type,
			ri.first_stack_trace,
			ri.assignee_id,
			ri.assignee_name,
			ri.assignee_email,
			ri.project_id,
			ri.project_name,
			ri.age
		ORDER BY
			ri.timestamp DESC;
	`, where)
	var rows []issueRow
	err := v.DB.Select(&rows, query, params...)
	if err != nil {
		return utils.RespondFail(c, 500, "Failed to fetch issues", err)
	}
	issues := make([]Issue, len(rows))
	for i, row := range rows {
		issue, _err := row.ToIssue()
		if _err != nil {
			fmt.Println("Error parsing issue:", _err)
			return utils.RespondFail(c, 500, "Failed to parse issue data", _err)
		}
		issues[i] = issue
	}

	return utils.RespondOK(c, issues, "")
}
