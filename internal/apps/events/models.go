package events

import (
	"encoding/json"
	"time"
)

type Project struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Assignee struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	FirstName string `json:"firstName"`
}

type Group struct {
	ID         string    `json:"id"`
	Status     string    `json:"status"`
	Assignee   *Assignee `json:"assignee"`
	EventCount int       `json:"eventCount"`
}

type Event struct {
	ID        string    `json:"id"`
	Type      string    `json:"type"`
	Level     string    `json:"level"`
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
}

type DayCount struct {
	Day   string `json:"day"`
	Count int    `json:"count"`
}

type Issue struct {
	Event     Event      `json:"event"`
	Project   Project    `json:"project"`
	Group     Group      `json:"group"`
	Last7Days []DayCount `json:"last7Days"`
}

type issueRow struct {
	EventID     string          `db:"event_id"`
	GroupID     string          `db:"group_id"`
	Message     string          `db:"message"`
	Type        string          `db:"type"`
	Level       string          `db:"level"`
	Status      string          `db:"status"`
	AssigneeID  *string         `db:"assignee_id"`
	FirstName   *string         `db:"first_name"`
	Email       *string         `db:"email"`
	UserID      *string         `db:"user_id"`
	ProjectID   string          `db:"project_id"`
	ProjectName string          `db:"project_name"`
	Last7Days   json.RawMessage `db:"last_7_days"`
	EventCount  int             `db:"event_count"`
}

func (ir *issueRow) ToIssue() Issue {
	var arr []DayCount
	if ir.Last7Days != nil {
		_ = json.Unmarshal(ir.Last7Days, &arr)
	}

	// ← default nil
	var assignee *Assignee = nil

	if ir.AssigneeID != nil {
		assignee = &Assignee{}
		assignee.ID = *ir.AssigneeID

		if ir.FirstName != nil {
			assignee.FirstName = *ir.FirstName
		}
		if ir.Email != nil {
			assignee.Email = *ir.Email
		}
	}

	return Issue{
		Event: Event{
			ID:      ir.EventID,
			Type:    ir.Type,
			Level:   ir.Level,
			Message: ir.Message,
		},
		Project: Project{
			ID:   ir.ProjectID,
			Name: ir.ProjectName,
		},
		Group: Group{
			ID:         ir.GroupID,
			Status:     ir.Status,
			Assignee:   assignee,
			EventCount: ir.EventCount,
		},
		Last7Days: arr,
	}
}
