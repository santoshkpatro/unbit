package models

type User struct {
	ID           string `db:"id"`
	Email        string `db:"email"`
	FirstName    string `db:"first_name"`
	LastName     string `db:"last_name"`
	IsAdmin      bool   `db:"is_admin"`
	PasswordHash string `db:"password_hash"`
	Salt         string `db:"salt"`
	CreatedAt    int64  `db:"created_at"`
	UpdatedAt    int64  `db:"updated_at"`
}
