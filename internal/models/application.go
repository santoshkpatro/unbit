package models

type Application struct {
	BaseModel

	Name        string `json:"name" db:"name"`
	CreatedBy   User
	Description string `json:"description" db:"description"`
}

type ApplicationPermission struct {
	BaseModel

	Application Application
	User        User
	Role        string `json:"role" db:"role"`
}
