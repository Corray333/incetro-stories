package transport

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"

	"github.com/Corray333/stories/internal/domains/user/storage"
	"github.com/Corray333/stories/internal/domains/user/types"
	"github.com/Corray333/stories/pkg/server/auth"
)

type Storage interface {
	InsertUser(user types.User) (int64, string, error)
	LoginUser(user types.User) (int64, string, error)
	CheckAndUpdateRefresh(id int64, refresh string) (string, error)
	SelectUser(id int64) (types.User, error)
}

// SignUp creates a new user it it doesn't exits and sends back refresh token, access token and user info
func SignUp(store *storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := types.User{}
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Failed to read request body", http.StatusBadRequest)
			slog.Error("Failed to read request body: " + err.Error())
			return
		}
		if err := json.Unmarshal(body, &user); err != nil {
			http.Error(w, "Failed to unmarshal request body", http.StatusBadRequest)
			slog.Error("Failed to unmarshal request body: " + err.Error())
			return
		}
		id, refresh, err := store.InsertUser(user)
		if err != nil {
			http.Error(w, "Failed to insert user", http.StatusInternalServerError)
			slog.Error("Failed to insert user: " + err.Error())
			return
		}
		user.ID = id

		token, err := auth.CreateToken(user.ID, auth.AccessTokenLifeTime)
		if err != nil {
			http.Error(w, "Failed to create token", http.StatusInternalServerError)
			slog.Error("Failed to create token: " + err.Error())
			return
		}
		user.Password = ""
		// fmt.Printf("User: %+v\n", user)
		if err := json.NewEncoder(w).Encode(struct {
			Authorization string     `json:"authorization"`
			Refresh       string     `json:"refresh"`
			User          types.User `json:"user"`
		}{Authorization: token,
			Refresh: refresh,
			User:    user,
		}); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			slog.Error("Failed to send response: " + err.Error())
			return
		}
	}
}

// LogIn logs in the user and sends back refresh token, access token and user info
func LogIn(store *storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := types.User{}
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Failed to read request body", http.StatusBadRequest)
			slog.Error("Failed to read request body: " + err.Error())
			return
		}
		if err := json.Unmarshal(body, &user); err != nil {
			http.Error(w, "Failed to unmarshal request body", http.StatusBadRequest)
			slog.Error("Failed to unmarshal request body: " + err.Error())
			return
		}
		id, refresh, err := store.LoginUser(user)
		if err != nil {
			http.Error(w, "Failed to insert user", http.StatusInternalServerError)
			slog.Error("Failed to insert user: " + err.Error())
			return
		}
		user.ID = id

		token, err := auth.CreateToken(user.ID, auth.AccessTokenLifeTime)
		if err != nil {
			http.Error(w, "Failed to create token", http.StatusInternalServerError)
			slog.Error("Failed to create token: " + err.Error())
			return
		}
		user.Password = ""
		// fmt.Printf("User: %+v\n", user)
		if err := json.NewEncoder(w).Encode(struct {
			Authorization string     `json:"authorization"`
			Refresh       string     `json:"refresh"`
			User          types.User `json:"user"`
		}{Authorization: token,
			Refresh: refresh,
			User:    user,
		}); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			slog.Error("Failed to send response: " + err.Error())
			return
		}
	}
}

// RefreshAccessToken creates a new access token
func RefreshAccessToken(store Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		refresh := r.Header.Get("Refresh")
		access, refresh, err := auth.RefreshAccessToken(store, refresh)
		if err != nil {
			http.Error(w, "Failed to refresh token", http.StatusInternalServerError)
			slog.Error("Failed to refresh token: " + err.Error())
			return
		}
		if err := json.NewEncoder(w).Encode(struct {
			Authorization string `json:"authorization"`
			Refresh       string `json:"refresh"`
		}{
			Authorization: access,
			Refresh:       refresh,
		}); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			slog.Error("Failed to encode response: " + err.Error())
			return
		}
	}
}

func GetUser(store Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		creds, err := auth.ExtractCredentials(token)
		if err != nil {
			http.Error(w, "Failed to get user id", http.StatusInternalServerError)
			slog.Error("Failed to get user id: " + err.Error())
			return
		}
		user, err := store.SelectUser(creds.ID)
		if err != nil {
			http.Error(w, "Failed to get user", http.StatusInternalServerError)
			slog.Error("Failed to get user: " + err.Error())
			return
		}
		if err := json.NewEncoder(w).Encode(user); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			slog.Error("Failed to encode response: " + err.Error())
			return
		}
	}
}
