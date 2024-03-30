package transport

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"os"
	"strconv"

	"github.com/Corray333/univer_cs/internal/domains/user/storage"
	"github.com/Corray333/univer_cs/internal/domains/user/types"
	"github.com/Corray333/univer_cs/pkg/server/auth"
	"github.com/go-chi/chi/v5"
)

type Storage interface {
	InsertUser(user types.User) (int64, string, error)
	LoginUser(user types.User) (int64, string, error)
	CheckAndUpdateRefresh(id int64, refresh string) (string, error)
	SelectUser(id string) (types.User, error)
	UpdateUser(user types.User) error
}

// SignUp creates a new user it it doesn't exits and sends back refresh token, access token and user info
func SignUp(store *storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := types.User{}
		user.Avatar = "http://localhost:3001/images/avatars/default_avatar.png"
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
			http.Error(w, "Wrong password or email", http.StatusForbidden)
			slog.Error("Failed to login user: " + err.Error())
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
		userId := chi.URLParam(r, "id")
		user, err := store.SelectUser(userId)
		if err != nil {
			http.Error(w, "Failed to get user", http.StatusInternalServerError)
			slog.Error("Failed to get user: " + err.Error())
			return
		}
		if err := json.NewEncoder(w).Encode(struct {
			User types.User `json:"user"`
		}{
			User: user,
		}); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			slog.Error("Failed to encode response: " + err.Error())
			return
		}
	}
}

func UpdateUser(store Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		creds, err := auth.ExtractCredentials(r.Header.Get("Authorization"))
		if err != nil {
			http.Error(w, "Failed to extract credentials", http.StatusBadRequest)
			slog.Error("Failed to extract credentials: " + err.Error())
			return
		}
		user, err := store.SelectUser(strconv.Itoa(int(creds.ID)))
		if err != nil {
			http.Error(w, "Failed to get user", http.StatusInternalServerError)
			slog.Error("Failed to get user: " + err.Error())
			return
		}
		file, _, err := r.FormFile("avatar")
		if err != nil && err.Error() != "http: no such file" {
			http.Error(w, "Failed to read file", http.StatusBadRequest)
			slog.Error("Failed to read file: " + err.Error())
			return
		}
		if file != nil {
			newFile, err := os.Create("../files/images/avatars/avatar" + strconv.Itoa(int(user.ID)) + ".png")
			if err != nil {
				http.Error(w, "Failed to create file", http.StatusInternalServerError)
				slog.Error("Failed to create file: " + err.Error())
				return
			}
			data, err := io.ReadAll(file)
			if err != nil {
				http.Error(w, "Failed to read file", http.StatusInternalServerError)
				slog.Error("Failed to read file: " + err.Error())
				return
			}
			if _, err := newFile.Write(data); err != nil {
				http.Error(w, "Failed to write file", http.StatusInternalServerError)
				slog.Error("Failed to write file: " + err.Error())
				return
			}
			user.Avatar = "http://localhost:3001/images/avatars/avatar" + strconv.Itoa(int(user.ID)) + ".png"
		}
		user.Username = r.FormValue("username")
		if err := store.UpdateUser(user); err != nil {
			http.Error(w, "Failed to update user", http.StatusInternalServerError)
			slog.Error("Failed to update user: " + err.Error())
			return
		}
	}
}
