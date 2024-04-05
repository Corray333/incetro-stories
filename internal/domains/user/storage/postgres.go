package storage

import (
	"fmt"

	"github.com/Corray333/univer_cs/internal/domains/user/types"
	"github.com/Corray333/univer_cs/pkg/server/auth"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Storage struct {
	db *sqlx.DB
}

// New creates a new storage and tables
func NewStorage(db *sqlx.DB) (*Storage, error) {

	_, err := db.Query(`
		CREATE TABLE IF NOT EXISTS users
		(
			user_id bigint NOT NULL GENERATED ALWAYS AS IDENTITY ( INCREMENT 1 START 1 MINVALUE 1 MAXVALUE 9223372036854775807 CACHE 1 ),
			name text COLLATE pg_catalog."default" NOT NULL,
			email text COLLATE pg_catalog."default" NOT NULL,
			password character varying(60) COLLATE pg_catalog."default" NOT NULL,
			avatar text COLLATE pg_catalog."default",
			CONSTRAINT users_pkey PRIMARY KEY (user_id),
			CONSTRAINT users_email_key UNIQUE (email)
		);
		CREATE TABLE IF NOT EXISTS user_token
		(
			user_id bigint NOT NULL,
			token text NOT NULL,
			PRIMARY KEY (user_id),
			FOREIGN KEY (user_id)
				REFERENCES public.users (user_id) MATCH SIMPLE
				ON UPDATE NO ACTION
				ON DELETE NO ACTION
				NOT VALID
		);
	`)
	return &Storage{db: db}, err
}

// InsertUser inserts a new user into the database and returns the id
func (s *Storage) InsertUser(user types.User) (int64, string, error) {
	passHash, err := auth.Hash(user.Password)
	if err != nil {
		return -1, "", err
	}
	user.Password = passHash

	rows := s.db.QueryRow(`
		INSERT INTO users (username, email, password, avatar) VALUES ($1, $2, $3, $4) RETURNING user_id;
	`, user.Username, user.Email, user.Password, "http://localhost:3001/images/avatars/default_avatar.png")

	if err := rows.Scan(&user.ID); err != nil {
		return -1, "", err
	}

	refresh, err := auth.CreateToken(user.ID, auth.RefreshTokenLifeTime)
	if err != nil {
		return -1, "", err
	}

	_, err = s.db.Queryx(`
		INSERT INTO user_token (user_id, token) VALUES ($1, $2);
	`, user.ID, refresh)
	if err != nil {
		return -1, "", err
	}

	return user.ID, refresh, nil
}

// LoginUser checks if the user exists and the password is correct
func (s *Storage) LoginUser(user types.User) (int64, string, error) {
	password := user.Password

	rows := s.db.QueryRow(`
		SELECT user_id, password FROM users WHERE email = $1;
	`, user.Email)

	if err := rows.Scan(&user.ID, &user.Password); err != nil {
		return -1, "", err
	}
	if !auth.Verify(user.Password, password) {
		return -1, "", fmt.Errorf("invalid password")
	}

	// Auto update refresh token
	refresh, err := auth.CreateToken(user.ID, auth.RefreshTokenLifeTime)
	if err != nil {
		return -1, "", err
	}

	_, err = s.db.Queryx(`
		INSERT INTO user_token (user_id, token) VALUES ($1, $2) ON CONFLICT (user_id) DO UPDATE SET token = $3;
	`, user.ID, refresh, refresh)
	if err != nil {
		return -1, "", err
	}

	return user.ID, refresh, nil
}

// CheckAndUpdateRefresh checks if the refresh token is valid and updates it
func (s *Storage) CheckAndUpdateRefresh(id int64, refresh string) (string, error) {
	rows, err := s.db.Queryx(`
		SELECT token FROM user_token WHERE user_id = $1 AND token = $2;
	`, id, refresh)
	if err != nil {
		return "", err
	}
	if !rows.Next() {
		return "", fmt.Errorf("invalid refresh token")
	}
	newRefresh, err := auth.CreateToken(id, auth.RefreshTokenLifeTime)
	if err != nil {
		return "", err
	}
	_, err = s.db.Queryx(`
		UPDATE user_token SET token = $1 WHERE user_id = $2;
	`, newRefresh, id)
	if err != nil {
		return "", err
	}
	return newRefresh, nil
}

func (s *Storage) SelectUser(id string) (types.User, error) {
	var user types.User
	rows, err := s.db.Queryx(`
		SELECT * FROM users WHERE user_id = $1;
	`, id)
	if err != nil {
		return user, err
	}
	if !rows.Next() {
		return user, fmt.Errorf("user not found")
	}
	if err := rows.StructScan(&user); err != nil {
		return user, err
	}
	user.Password = ""
	return user, nil
}

func (s *Storage) UpdateUser(user types.User) error {
	fmt.Println()
	fmt.Println(user)
	fmt.Println()
	_, err := s.db.Queryx(`
		UPDATE users SET username = $1, email = $2, avatar = $3 WHERE user_id = $4;
	`, user.Username, user.Email, user.Avatar, user.ID)
	return err
}
