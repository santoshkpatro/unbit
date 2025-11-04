package events

import (
	"database/sql"
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
	Message   string    `json:"message"`
	Level     string    `json:"level"`
	Timestamp time.Time `json:"timestamp"`
	CreatedAt time.Time `json:"createdAt"`
	Project   Project   `json:"project"`
	Group     Group     `json:"group"`
}

type eventRow struct {
	EventID   string    `db:"event_id"`
	Type      string    `db:"event_type"`
	Message   string    `db:"event_message"`
	Level     string    `db:"event_level"`
	Timestamp time.Time `db:"event_timestamp"`
	CreatedAt time.Time `db:"created_at"`

	ProjectID   string `db:"project_id"`
	ProjectName string `db:"project_name"`

	GroupID           string         `db:"group_id"`
	GroupStatus       string         `db:"group_status"`
	GroupEventCount   int            `db:"group_event_count"`
	AssigneeID        sql.NullString `db:"assignee_id"`
	AssigneeEmail     sql.NullString `db:"assignee_email"`
	AssigneeFirstName sql.NullString `db:"assignee_first_name"`
}

func (r eventRow) toEvent() Event {
	var a *Assignee
	if r.AssigneeID.Valid {
		a = &Assignee{
			ID:        r.AssigneeID.String,
			Email:     r.AssigneeEmail.String,
			FirstName: r.AssigneeFirstName.String,
		}
	}

	return Event{
		ID:        r.EventID,
		Type:      r.Type,
		Message:   r.Message,
		Level:     r.Level,
		Timestamp: r.Timestamp,
		CreatedAt: r.CreatedAt,
		Project: Project{
			ID:   r.ProjectID,
			Name: r.ProjectName,
		},
		Group: Group{
			ID:         r.GroupID,
			Status:     r.GroupStatus,
			Assignee:   a,
			EventCount: r.GroupEventCount,
		},
	}
}
