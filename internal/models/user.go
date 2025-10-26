package models

type User struct {
	ID        string  `db:"id" json:"id"`
	Email     string  `db:"email" json:"email"`
	FirstName string  `db:"first_name" json:"firstName"`
	LastName  *string `db:"last_name" json:"lastName"`
	IsAdmin   bool    `db:"is_admin" json:"isAdmin"`
}
