package models

import "time"

type StackFrame struct {
	Function string `json:"function"`
	File     string `json:"file"`
	Line     int    `json:"line"`
}

type Event struct {
	Message    string       `json:"message"`
	Level      string       `json:"level"`
	Timestamp  time.Time    `json:"timestamp"`
	StackTrace []StackFrame `json:"stacktrace"`
}

type Payload struct {
	DSNToken string `json:"dsnToken"`
	Event    Event  `json:"event"`
}

type Issue struct {
	ID          string     `db:"id"`
	ProjectID   string     `db:"project_id"`
	Summary     string     `db:"summary"`
	Fingerprint string     `db:"fingerprint"`
	LastSeenAt  *time.Time `db:"last_seen_at"`
	AssigneeID  *string    `db:"assignee_id"`
	Status      string     `db:"status"`
	EventCount  int64      `db:"event_count"`
	CreatedAt   time.Time  `db:"created_at"`
	UpdatedAt   time.Time  `db:"updated_at"`
}
