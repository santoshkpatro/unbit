package issues

import (
	"database/sql"
	"time"

	_ "github.com/lib/pq"
)

type Assignee struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	FirstName string `json:"firstName"`
}

type Project struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Issue struct {
	ID         string         `json:"id"`
	ProjectId  string         `json:"-"`
	Project    *Project       `json:"project"`
	Summary    string         `json:"summary"`
	AssigneeID sql.NullString `json:"-"`
	Assignee   *Assignee      `json:"assignee"`
	Status     string         `json:"status"`
	LastSeenAt time.Time      `json:"lastSeenAt"`
	EventCount int            `json:"eventCount"`
	CreatedAt  time.Time      `json:"createdAt"`
}

type IssueRow struct {
	ID                string         `db:"id"`
	ProjectId         string         `db:"project_id"`
	ProjectName       string         `db:"project_name"`
	Summary           string         `db:"summary"`
	AssigneeID        sql.NullString `db:"assignee_id"`
	Status            string         `db:"status"`
	LastSeenAt        time.Time      `db:"last_seen_at"`
	EventCount        int            `db:"event_count"`
	CreatedAt         time.Time      `db:"created_at"`
	AssigneeEmail     sql.NullString `db:"assignee_email"`
	AssigneeFirstName sql.NullString `db:"assignee_first_name"`
}
