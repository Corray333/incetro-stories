package types

type Banner struct {
	ID          int    `json:"id,omitempty" db:"banner_id"`
	Title       string `json:"title" db:"title"`
	Description string `json:"description" db:"description"`
	CreatedAt   string `json:"created_at,omitempty" db:"created_at"`
	Views       int    `json:"views"`
	Lang        string `json:"lang" db:"-"`
}

type Story struct {
	ID        int      `json:"id,omitempty" db:"story_id"`
	CreatedAt string   `json:"created_at,omitempty" db:"created_at"`
	Banners   []Banner `json:"banners,omitempty" db:"banners"`
	Creator   int      `json:"creator,omitempty" db:"creator"`
}
