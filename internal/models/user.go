package models

type User struct {
	BaseModel

	Email        string `json:"email" db:"email"`
	FullName     string `json:"fullName" db:"full_name"`
	PasswordHash string `json:"-" db:"password_hash"`
	Salt         string `json:"-" db:"salt"`
	IsAdmin      bool   `json:"isAdmin" db:"is_admin"`
}
