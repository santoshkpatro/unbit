package auth

import "time"

type loginData struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type User struct {
	ID           string    `db:"id"`
	Email        string    `db:"email" json:"email"`
	FirstName    string    `db:"first_name" json:"firstName"`
	LastName     *string   `db:"last_name" json:"lastName"`
	IsAdmin      bool      `db:"is_admin" json:"isAdmin"`
	IsActive     bool      `db:"is_active" json:"-"`
	PasswordHash string    `db:"password_hash" json:"-"`
	Salt         string    `db:"salt" json:"-"`
	CreatedAt    time.Time `db:"created_at" json:"-"`
	UpdatedAt    time.Time `db:"updated_at" json:"-"`
}
