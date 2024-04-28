package types

type Project struct {
	ID          int    `json:"id,omitempty" db:"project_id"`
	Name        string `json:"name" db:"name"`
	Description string `json:"description" db:"description"`
	CreatedAt   string `json:"created_at,omitempty" db:"created_at"`
	Creator     int    `json:"creator,omitempty" db:"creator"`
	Cover       string `json:"cover" db:"cover"`
}
