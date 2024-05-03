package transport

import (
	"encoding/json"
	"io"
	"log/slog"
	"mime/multipart"
	"net/http"
	"strconv"

	"github.com/Corray333/univer_cs/internal/domains/story/types"
	"github.com/Corray333/univer_cs/pkg/server/auth"
	"github.com/go-chi/chi/v5"
)

const MaxFileSize = 5 << 20

type Storage interface {
	SelectStories(project_id, story_id, banner_id, creator, offset, lang string) ([]types.Story, error)
	InsertBanner(project_id string, story_id string, uid int, banner types.Banner, file multipart.File, fileName string) error
	InsertView(user_id int, banner_id string) error
	UpdateBanner(banner types.Banner, expires_at string, file multipart.File, fileName string) error
}

func GetStories(store Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		project_id := chi.URLParam(r, "project_id")
		story_id := r.URL.Query().Get("story_id")
		banner_id := r.URL.Query().Get("banner_id")
		creator := r.URL.Query().Get("creator")
		offset := r.URL.Query().Get("offset")
		lang := r.URL.Query().Get("lang")

		stories, err := store.SelectStories(project_id, story_id, banner_id, creator, offset, lang)
		if err != nil {
			http.Error(w, "Failed to get stories", http.StatusInternalServerError)
			slog.Error("Failed to get stories: " + err.Error())
			return
		}
		if err := json.NewEncoder(w).Encode(struct {
			Stories []types.Story `json:"stories"`
		}{Stories: stories}); err != nil {
			http.Error(w, "Failed to  response", http.StatusInternalServerError)
			slog.Error("Failed to  response: " + err.Error())
			return
		}
	}
}

// NewBanner creates a new banner in the database and saves the image
func NewBanner(store Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		story_id := r.URL.Query().Get("story_id")
		project_id := chi.URLParam(r, "project_id")
		// Limit max input length
		if err := r.ParseMultipartForm(MaxFileSize); err != nil {
			slog.Error(err.Error())
			http.Error(w, "Failed to read file", http.StatusBadRequest)
		}
		file, fileHeader, err := r.FormFile("file")
		if err != nil {
			slog.Error(err.Error())
			http.Error(w, "Failed to read file", http.StatusBadRequest)
		}
		defer file.Close()

		var langs []types.BannerLang

		// Unmarshal the banner
		langsRaw := r.FormValue("langs")

		if err := json.Unmarshal([]byte(langsRaw), &langs); err != nil {
			http.Error(w, "Failed to unmarshal banner", http.StatusInternalServerError)
			slog.Error("Failed to unmarshal banner: " + err.Error())
			return
		}

		banner := types.Banner{Langs: langs}

		creds, err := auth.ExtractCredentials(r.Header.Get("Authorization"))
		if err != nil {
			http.Error(w, "Failed to get user id", http.StatusInternalServerError)
			slog.Error("Failed to get user id: " + err.Error())
			return
		}

		if err := store.InsertBanner(project_id, story_id, creds.ID, banner, file, fileHeader.Filename); err != nil {
			http.Error(w, "Failed to insert banner", http.StatusInternalServerError)
			slog.Error("Failed to insert banner: " + err.Error())
			return
		}

		w.WriteHeader(http.StatusCreated)
	}
}

func NewView(store Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		creds, err := auth.ExtractCredentials(token)
		if err != nil {
			http.Error(w, "Failed to get user id", http.StatusInternalServerError)
			slog.Error("Failed to get user id: " + err.Error())
			return
		}
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Failed reading body", http.StatusInternalServerError)
			slog.Error("Failed reading body: " + err.Error())
			return
		}
		bodyUnmarshalled := struct {
			BannerId int64 `json:"banner_id"`
		}{}
		if err := json.Unmarshal(body, &bodyUnmarshalled); err != nil {
			http.Error(w, "Failed to unmarshal body", http.StatusInternalServerError)
			slog.Error("Failed to unmarshal body: " + err.Error())
			return
		}

		if err = store.InsertView(creds.ID, strconv.Itoa(int(bodyUnmarshalled.BannerId))); err != nil {
			http.Error(w, "Failed to insert view", http.StatusInternalServerError)
			slog.Error("Failed to insert view: " + err.Error())
			return
		}

	}
}

func UpdateBanner(store Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body := r.FormValue("banner")
		var banner types.Banner
		if err := json.Unmarshal([]byte(body), &banner); err != nil {
			http.Error(w, "Failed unmarshalling body", http.StatusInternalServerError)
			slog.Error("Failed unmarshalling body: " + err.Error())
			return
		}

		expires_at := r.FormValue("expires_at")

		file, fileHeader, err := r.FormFile("file")
		fileName := ""
		if err != nil {
			if err.Error() != "http: no such file" {
				slog.Error(err.Error())
				http.Error(w, "Failed to read file", http.StatusBadRequest)
				return
			}
		} else {
			fileName = fileHeader.Filename
			defer file.Close()
		}

		if err := store.UpdateBanner(banner, expires_at, file, fileName); err != nil {
			http.Error(w, "Failed to update banner", http.StatusInternalServerError)
			slog.Error("Failed to update banner: " + err.Error())
			return
		}
	}
}
