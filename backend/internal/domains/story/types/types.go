package types

type Project struct {
	ID          int    `json:"id" db:"project_id"`
	Name        string `json:"name" db:"name"`
	Description string `json:"description" db:"description"`
	CreatedAt   string `json:"created_at" db:"created_at"`
	Creator     int    `json:"creator" db:"creator"`
	Cover       string `json:"cover" db:"cover"`
	ProjectID   int    `json:"project_id" db:"project_id"`
}

type Banner struct {
	ID        int          `json:"id" db:"banner_id"`
	Langs     []BannerLang `json:"langs" db:"langs"`
	CreatedAt int          `json:"created_at" db:"created_at"`
	MediaURL  string       `json:"media_url" db:"media_url"`
	Views     int          `json:"views"`
}

type BannerLang struct {
	Lang        string `json:"lang" db:"lang"`
	Title       string `json:"title" db:"title"`
	Description string `json:"description" db:"description"`
}

type Story struct {
	ID        int      `json:"id" db:"story_id"`
	CreatedAt int      `json:"created_at" db:"created_at"`
	ExpiresAt int      `json:"expires_at" db:"expires_at"`
	Banners   []Banner `json:"banners" db:"banners"`
	Creator   int      `json:"creator" db:"creator"`
}
