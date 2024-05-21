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
	SelectStories(project_id, story_id, banner_id, creator, offset, lang, all string) ([]types.Story, error)
	InsertBanner(project_id string, story_id string, uid int, banner types.Banner, file multipart.File, fileName string) error
	InsertView(user_id int, banner_id string) error
	UpdateBanner(banner types.Banner, expires_at string, file multipart.File, fileName string) error
}

type GetStoriesResponse struct {
	Stories []types.Story `json:"stories"`
}

// GetStories godoc
// @Summary Get stories
// @Description Get stories by project ID
// @Tags stories
// @Produce json
// @Param project_id path string true "Project ID"
// @Param story_id query string false "Story ID"
// @Param banner_id query string false "Banner ID"
// @Param creator query string false "Creator"
// @Param offset query string false "Offset"
// @Param lang query string false "Language"
// @Param all query string false "All"
// @Success 200 {object} GetStoriesResponse "Stories"
// @Failure 500 {string} string "Failed to get stories"
// @Router /api/projects/{project_id}/stories [get]
func GetStories(store Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		project_id := chi.URLParam(r, "project_id")
		story_id := r.URL.Query().Get("story_id")
		banner_id := r.URL.Query().Get("banner_id")
		creator := r.URL.Query().Get("creator")
		offset := r.URL.Query().Get("offset")
		lang := r.URL.Query().Get("lang")
		all := r.URL.Query().Get("all")

		stories, err := store.SelectStories(project_id, story_id, banner_id, creator, offset, lang, all)
		if err != nil {
			http.Error(w, "Failed to get stories", http.StatusInternalServerError)
			slog.Error("Failed to get stories: " + err.Error())
			return
		}
		if err := json.NewEncoder(w).Encode(GetStoriesResponse{Stories: stories}); err != nil {
			http.Error(w, "Failed to response", http.StatusInternalServerError)
			slog.Error("Failed to response: " + err.Error())
			return
		}
	}
}

// NewBanner godoc
// @Summary Create a new banner
// @Description Create a new banner and save the image
// @Tags banners
// @Accept multipart/form-data
// @Param project_id path string true "Project ID"
// @Param story_id query string true "Story ID"
// @Param file formData file true "File"
// @Param langs formData string true "Langs"
// @Success 201 {string} string "Created"
// @Failure 400 {string} string "Failed to read file"
// @Failure 500 {string} string "Failed to unmarshal banner"
// @Router /api/projects/{project_id}/banners [post]
func NewBanner(store Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		story_id := r.URL.Query().Get("story_id")
		project_id := chi.URLParam(r, "project_id")
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

// NewView godoc
// @Summary Create a new view
// @Description Create a new view for a banner
// @Tags views
// @Accept json
// @Param Authorization header string true "Bearer token"
// @Param banner_id body int true "Banner ID"
// @Success 200 {string} string "Created"
// @Failure 500 {string} string "Failed to insert view"
// @Router /api/views [post]
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

// UpdateBanner godoc
// @Summary Update a banner
// @Description Update a banner with new data and optionally a new file
// @Tags banners
// @Accept multipart/form-data
// @Param banner formData string true "Banner data"
// @Param expires_at formData string false "Expiration date"
// @Param file formData file false "File"
// @Success 200 {string} string "Updated"
// @Failure 400 {string} string "Failed to read file"
// @Failure 500 {string} string "Failed to update banner"
// @Router /api/banners [put]
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
