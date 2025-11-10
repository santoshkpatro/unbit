package models

import "time"

type Frame struct {
	Function string `json:"function"`
	File     string `json:"file"`
	Line     int    `json:"line"`
	Code     string `json:"code"`
}

type Issue struct {
	ID          string    `db:"id"`
	ProjectID   string    `db:"project_id"`
	Fingerprint string    `db:"fingerprint"`
	AssigneeID  *string   `db:"assignee_id"`
	Status      string    `db:"status"`
	EventCount  int64     `db:"event_count"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

type Properties struct {
	Type       string  `json:"type"`
	Message    string  `json:"message"`
	Level      string  `json:"level"`
	Stacktrace []Frame `json:"stacktrace"`
}

type Event struct {
	Timestamp  time.Time  `json:"timestamp"`
	Properties Properties `json:"properties"`
}

type Payload struct {
	DSNToken string `json:"dsnToken"`
	Event    Event  `json:"event"`
}
