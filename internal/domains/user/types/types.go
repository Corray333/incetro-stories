package types

type User struct {
	ID       int64  `json:"id,omitempty" db:"user_id"`
	Name     string `json:"name" db:"name"`
	Email    string `json:"email" db:"email"`
	Password string `json:"password" db:"password"`
}
