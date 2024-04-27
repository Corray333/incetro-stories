package transport

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"os"
	"strconv"

	"github.com/Corray333/univer_cs/internal/domains/story/types"
	"github.com/Corray333/univer_cs/pkg/server/auth"
)

const MaxFileSize = 64 << 20

type Storage interface {
	SelectStories(story_id, banner_id, creator, offset, lang string) ([]types.Story, error)
	InsertBanner(story_id string, uid int, banners []types.Banner) (int, error)
	InsertView(user_id int, banner_id string) error
	// UpdateBanner(banner types.Banner) error
}

func GetStories(store Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		story_id := r.URL.Query().Get("story_id")
		banner_id := r.URL.Query().Get("banner_id")
		creator := r.URL.Query().Get("creator")
		offset := r.URL.Query().Get("offset")
		lang := r.URL.Query().Get("lang")

		stories, err := store.SelectStories(story_id, banner_id, creator, offset, lang)
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

		var banners []types.Banner

		// Unmarshal the banner
		data := r.FormValue("data")
		if err := json.Unmarshal([]byte(data), &banners); err != nil {
			http.Error(w, "Failed to unmarshal banner", http.StatusInternalServerError)
			slog.Error("Failed to unmarshal banner: " + err.Error())
			return
		}

		creds, err := auth.ExtractCredentials(r.Header.Get("Authorization"))
		if err != nil {
			http.Error(w, "Failed to get user id", http.StatusInternalServerError)
			slog.Error("Failed to get user id: " + err.Error())
			return
		}

		id, err := store.InsertBanner(story_id, creds.ID, banners)
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

// func UpdateBanner(store Storage) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		body, err := io.ReadAll(r.Body)
// 		if err != nil {
// 			http.Error(w, "Failed reading body", http.StatusInternalServerError)
// 			slog.Error("Failed reading body: " + err.Error())
// 			return
// 		}
// 		var banner types.Banner
// 		if err := json.Unmarshal(body, &banner); err != nil {
// 			http.Error(w, "Failed unmarshalling body", http.StatusInternalServerError)
// 			slog.Error("Failed unmarshalling body: " + err.Error())
// 			return
// 		}
// 		id, err := strconv.Atoi(chi.URLParam(r, "id"))
// 		if err != nil {
// 			http.Error(w, "Failed to get banner id", http.StatusInternalServerError)
// 			slog.Error("Failed to get banner id: " + err.Error())
// 			return
// 		}
// 		banner.ID = id
// 		if err = store.UpdateBanner(banner); err != nil {
// 			http.Error(w, "Failed to update banner", http.StatusInternalServerError)
// 			slog.Error("Failed to update banner: " + err.Error())
// 			return
// 		}
// 	}
// }
