package models

import "time"

type Project struct {
	ID          string    `db:"id" json:"id"`
	Name        string    `db:"name" json:"name"`
	Description string    `db:"description" json:"description"`
	DsnToken    *string   `db:"dsn_token" json:"dsnToken"`
	CreatedAt   time.Time `db:"created_at" json:"createdAt"`
}

type ProjectMemmber struct {
	ID        string    `db:"id" json:"id"`
	ProjectID string    `db:"project_id" json:"projectId"`
	UserID    string    `db:"user_id" json:"userId"`
	Role      string    `db:"role" json:"role"`
	JoinedAt  time.Time `db:"joined_at" json:"joinedAt"`
	CreatedAt time.Time `db:"created_at" json:"createdAt"`
}
