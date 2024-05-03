package transport

import (
	"encoding/json"
	"log/slog"
	"mime/multipart"
	"net/http"

	"github.com/Corray333/univer_cs/internal/domains/project/types"
	"github.com/Corray333/univer_cs/pkg/server/auth"
	"github.com/go-chi/chi/v5"
)

const MaxFileSize = 5 << 20

type Storage interface {
	InsertProject(uid int, cover multipart.File, project types.Project) error
	GetProjects(project_id string, uid int) ([]types.Project, error)
}

func NewProject(store Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseMultipartForm(MaxFileSize); err != nil {
			http.Error(w, "Failed to read file", http.StatusBadRequest)
			slog.Error("Failed to read file: " + err.Error())
			return
		}
		cover, _, err := r.FormFile("cover")
		if err != nil {
			http.Error(w, "Failed to read cover", http.StatusBadRequest)
			slog.Error("Failed to read cover: " + err.Error())
			return
		}
		data := r.FormValue("data")

		var project types.Project

		if err := json.Unmarshal([]byte(data), &project); err != nil {
			http.Error(w, "Failed to read project data", http.StatusBadRequest)
			slog.Error("Failed to read project data: " + err.Error())
			return
		}

		creds, err := auth.ExtractCredentials(r.Header.Get("Authorization"))

		if err != nil {
			http.Error(w, "Failed to extract credentials", http.StatusUnauthorized)
			slog.Error("Failed to extract credentials: " + err.Error())
			return
		}

		if err := store.InsertProject(creds.ID, cover, project); err != nil {
			http.Error(w, "Failed to insert project", http.StatusInternalServerError)
			slog.Error("Failed to insert project: " + err.Error())
			return
		}
		w.WriteHeader(http.StatusCreated)
	}
}

func GetProjects(store Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		project_id := chi.URLParam(r, "project_id")
		creds, err := auth.ExtractCredentials(r.Header.Get("Authorization"))

		if err != nil {
			http.Error(w, "Failed to extract credentials", http.StatusUnauthorized)
			slog.Error("Failed to extract credentials: " + err.Error())
			return
		}

		projects, err := store.GetProjects(project_id, creds.ID)

		if err != nil {
			http.Error(w, "Failed to get projects", http.StatusInternalServerError)
			slog.Error("Failed to get projects: " + err.Error())
			return
		}

		if err := json.NewEncoder(w).Encode(struct {
			Projects []types.Project `json:"projects"`
		}{
			Projects: projects,
		}); err != nil {
			http.Error(w, "Failed to encode projects", http.StatusInternalServerError)
			slog.Error("Failed to encode projects: " + err.Error())
			return
		}
	}
}
