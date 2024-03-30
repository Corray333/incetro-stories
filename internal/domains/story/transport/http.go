package transport

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/Corray333/univer_cs/internal/domains/story/storage"
	"github.com/Corray333/univer_cs/internal/domains/story/types"
	"github.com/Corray333/univer_cs/pkg/server/auth"
	"github.com/go-chi/chi/v5"
)

const MaxFileSize = 64 << 20

func GetStories(store storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		story_id := r.URL.Query().Get("stories_id")
		banner_id := r.URL.Query().Get("banner_id")
		creator := r.URL.Query().Get("creator")
		filter := "WHERE "
		if story_id != "" {
			filter = filter + "stories.stories_id = " + story_id
		}
		if banner_id != "" {
			if story_id != "" {
				filter = filter + " AND "
			}
			filter = filter + "banners.banner_id = " + banner_id
		}
		if creator != "" {
			if story_id != "" || banner_id != "" {
				filter = filter + " AND "
			}
			filter = filter + "stories.creator = " + creator
		}
		if filter == "WHERE " {
			filter = ""
		}
		filter += " ORDER BY stories.stories_id, banners.created_at"
		offset := r.URL.Query().Get("offset")
		if offset != "" {
			filter = filter + " OFFSET " + offset
		}

		stories, err := store.SelectStories(filter)
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
func NewBanner(store storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		story_id := r.URL.Query().Get("story_id")
		if story_id == "" {
			var story types.Story
			creds, err := auth.ExtractCredentials(r.Header.Get("Authorization"))
			if err != nil {
				http.Error(w, "Failed to get user id", http.StatusInternalServerError)
				slog.Error("Failed to get user id: " + err.Error())
				return
			}
			story.Creator = creds.ID
			id, err := store.InsertStory(story)
			if err != nil {
				http.Error(w, "Failed to insert story", http.StatusInternalServerError)
				slog.Error("Failed to insert story: " + err.Error())
				return
			}
			story_id = strconv.Itoa(int(id))
		}
		// Limit max input length
		if err := r.ParseMultipartForm(64 << 20); err != nil {
			slog.Error(err.Error())
			http.Error(w, "Failed to read file", http.StatusBadRequest)
		}
		file, _, err := r.FormFile("file")
		if err != nil {
			slog.Error(err.Error())
			http.Error(w, "Failed to read file", http.StatusBadRequest)
		}
		defer file.Close()

		var banner types.Banner

		banner.Name = r.FormValue("name")
		banner.Description = r.FormValue("description")

		id, err := store.InsertBanner(story_id, banner)
		if err != nil {
			http.Error(w, "Failed to insert banner", http.StatusInternalServerError)
			slog.Error("Failed to insert banner: " + err.Error())
			return
		}

		// TODO: replace with remote storage
		newFile, _ := os.Create("../files/images/banners/banner" + strconv.Itoa(int(id)) + ".png")
		defer newFile.Close()
		if _, err := io.Copy(newFile, file); err != nil {
			http.Error(w, "Failed to save file", http.StatusInternalServerError)
			slog.Error("Failed to save file: " + err.Error())
			return
		}
		w.WriteHeader(http.StatusCreated)
	}
}

func NewView(store storage.Storage) http.HandlerFunc {
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

// UpdateBannerMedia updates the media attribute of the banner
func UpdateBannerMedia(store storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseMultipartForm(MaxFileSize); err != nil {
			slog.Error(err.Error())
			http.Error(w, "Failed to read file", http.StatusBadRequest)
			return
		}
		file, _, err := r.FormFile("file")
		if err != nil {
			slog.Error(err.Error())
			http.Error(w, "Failed to read file", http.StatusBadRequest)
			return
		}
		defer file.Close()

		// Save banner
		newFile, _ := os.Create("../files/images/banners/banner" + chi.URLParam(r, "id") + ".png")
		defer newFile.Close()
		if _, err := io.Copy(newFile, file); err != nil {
			http.Error(w, "Failed to save file", http.StatusInternalServerError)
			slog.Error("Failed to save file: " + err.Error())
			return
		}
		w.WriteHeader(http.StatusCreated)

		w.WriteHeader(http.StatusOK)

	}
}

func UpdateStoryTimestamp(store storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Failed reading body", http.StatusInternalServerError)
			slog.Error("Failed reading body: " + err.Error())
			return
		}
		bodyUnmarshalled := struct {
			Timestamp time.Time `json:"timestamp"`
		}{}
		if err := json.Unmarshal(body, &bodyUnmarshalled); err != nil {
			http.Error(w, "Failed to unmarshal body", http.StatusInternalServerError)
			slog.Error("Failed to unmarshal body: " + err.Error())
			return
		}

		if err := store.UpdateBannerTimestamp(chi.URLParam(r, "id"), bodyUnmarshalled.Timestamp); err != nil {
			http.Error(w, "Failed to change time of banner expire", http.StatusInternalServerError)
			slog.Error("Failed to change time of banner expire: " + err.Error())
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

// UpdateBannerName updates the name of the banner
func UpdateBannerName(store storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Failed reading body", http.StatusInternalServerError)
			slog.Error("Failed reading body: " + err.Error())
			return
		}
		bodyUnmarshalled := struct {
			Name string `json:"name"`
		}{}
		if err := json.Unmarshal(body, &bodyUnmarshalled); err != nil {
			http.Error(w, "Failed to unmarshal body", http.StatusInternalServerError)
			slog.Error("Failed to unmarshal body: " + err.Error())
			return
		}

		if err := store.UpdateBannerName(chi.URLParam(r, "id"), bodyUnmarshalled.Name); err != nil {
			http.Error(w, "Failed to change name of banner", http.StatusInternalServerError)
			slog.Error("Failed to change name of banner: " + err.Error())
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

// UpdateBannerDescription updates the description of the banner
func UpdateBannerDescription(store storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Failed reading body", http.StatusInternalServerError)
			slog.Error("Failed reading body: " + err.Error())
			return
		}
		bodyUnmarshalled := struct {
			Description string `json:"description"`
		}{}
		if err := json.Unmarshal(body, &bodyUnmarshalled); err != nil {
			http.Error(w, "Failed to unmarshal body", http.StatusInternalServerError)
			slog.Error("Failed to unmarshal body: " + err.Error())
			return
		}

		if err := store.UpdateBannerDescription(chi.URLParam(r, "id"), bodyUnmarshalled.Description); err != nil {
			http.Error(w, "Failed to change description of banner", http.StatusInternalServerError)
			slog.Error("Failed to change description of banner: " + err.Error())
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}
