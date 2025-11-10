package issues

import (
	"encoding/json"
	"time"
)

type Assignee struct {
	ID    string  `json:"id"`
	Name  *string `json:"name"`
	Email *string `json:"email"`
}

type Project struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type IssueCountReport struct {
	Day        string `json:"date"`
	EventCount int    `json:"eventCount"`
}

type FirstStackTrace struct {
	Function string `json:"function"`
	File     string `json:"file"`
	Line     int    `json:"line"`
	Code     string `json:"code"`
}

type Issue struct {
	ID              string             `json:"id"`
	EventID         string             `json:"eventId"`
	EventCount      int                `json:"eventCount"`
	Timestamp       time.Time          `json:"timestamp"`
	Message         string             `json:"message"`
	Level           string             `json:"level"`
	Type            string             `json:"type"`
	Assignee        *Assignee          `json:"assignee"`
	Project         Project            `json:"project"`
	IssueCount      []IssueCountReport `json:"issueCountReport"`
	FirstStackTrace FirstStackTrace    `json:"firstStackTrace"`
	Age             int                `json:"age"`
}

type issueRow struct {
	ID               string          `db:"id"`
	EventID          string          `db:"event_id"`
	EventCount       int             `db:"event_count"`
	Timestamp        time.Time       `db:"timestamp"`
	Message          string          `db:"message"`
	Level            string          `db:"level"`
	Type             string          `db:"type"`
	FirstStackTrace  json.RawMessage `db:"first_stack_trace"`
	AssigneeID       *string         `db:"assignee_id"`
	AssigneeName     *string         `db:"assignee_name"`
	AssigneeEmail    *string         `db:"assignee_email"`
	ProjectID        string          `db:"project_id"`
	ProjectName      string          `db:"project_name"`
	IssueCountReport json.RawMessage `db:"issue_count_report"`
	Age              int             `db:"age"`
}

func (ir *issueRow) ToIssue() (Issue, error) {
	var issueCount []IssueCountReport
	if err := json.Unmarshal(ir.IssueCountReport, &issueCount); err != nil {
		return Issue{}, err
	}

	var firstStackTrace FirstStackTrace
	if err := json.Unmarshal(ir.FirstStackTrace, &firstStackTrace); err != nil {
		return Issue{}, err
	}

	// ← default nil
	var assignee *Assignee = nil

	if ir.AssigneeID != nil {
		assignee = &Assignee{
			ID:    *ir.AssigneeID,
			Name:  ir.AssigneeName,
			Email: ir.AssigneeEmail,
		}
	}

	return Issue{
		ID:         ir.ID,
		EventID:    ir.EventID,
		EventCount: ir.EventCount,
		Timestamp:  ir.Timestamp,
		Message:    ir.Message,
		Level:      ir.Level,
		Type:       ir.Type,
		Assignee:   assignee,
		Project: Project{
			ID:   ir.ProjectID,
			Name: ir.ProjectName,
		},
		IssueCount:      issueCount,
		FirstStackTrace: firstStackTrace,
		Age:             ir.Age,
	}, nil
}
