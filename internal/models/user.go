package models

import "time"

type User struct {
	BaseModel

	Email        string    `json:"email" db:"email"`
	FullName     string    `json:"fullName" db:"full_name"`
	PasswordHash string    `json:"-" db:"password_hash"`
	Salt         string    `json:"-" db:"salt"`
	IsActive     bool      `json:"isActive" db:"is_active"`
	CreatedAt    time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt    time.Time `json:"updatedAt" db:"updated_at"`
}
