package projects

type Project struct {
	ID          string `db:"id" json:"id"`
	Name        string `db:"name" json:"name"`
	Description string `db:"description" json:"description"`
	DsnToken    string `db:"dsn_token" json:"dsnToken"`
	CreatedAt   string `db:"created_at" json:"createdAt"`
	UpdatedAt   string `db:"updated_at" json:"-"`
}

type ProjectNew struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
}
