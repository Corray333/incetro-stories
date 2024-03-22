package transport

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"os"
	"strconv"

	"github.com/Corray333/stories/internal/domains/story/types"
	"github.com/Corray333/stories/pkg/server/auth"
)

type Storage interface {
	SelectStories(filter string) ([]types.Story, error)
	InsertStory(story types.Story) (int64, error)
	InsertBanner(storyId string, banner types.Banner) (int64, error)
	InsertView(userId int64, bannerId string) error
}

func GetStories(store Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		story_id := r.URL.Query().Get("stories_id")
		banner_id := r.URL.Query().Get("banner_id")
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
		if filter == "WHERE " {
			filter = ""
		}
		filter += " GROUP BY stories.stories_id, banners.banner_id "
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

// NewStories creates a new story in the database and sends back the id
func NewStories(store Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var story types.Story
		// credentials, err := auth.ExtractCredentials(r.Header.Get("Authorization"))
		// if err != nil {
		// 	http.Error(w, "Failed to extract credentials", http.StatusBadRequest)
		// 	slog.Error("Failed to extract credentials: " + err.Error())
		// 	return
		// }
		// story.UserID = credentials.ID
		id, err := store.InsertStory(story)
		if err != nil {
			http.Error(w, "Failed to insert story", http.StatusInternalServerError)
			slog.Error("Failed to insert story: " + err.Error())
			return
		}
		if _, err := w.Write([]byte(`{"id":` + strconv.Itoa(int(id)) + `}`)); err != nil {
			http.Error(w, "Failed to response", http.StatusInternalServerError)
			slog.Error("Failed to response: " + err.Error())
			return
		}
	}
}

// NewBanner creates a new banner in the database and saves the image
func NewBanner(store Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Limit max input length
		if err := r.ParseMultipartForm(32 << 20); err != nil {
			slog.Error(err.Error())
			http.Error(w, "Failed to read file", http.StatusBadRequest)
		}
		file, _, err := r.FormFile("img")

		if err != nil {
			slog.Error(err.Error())
			http.Error(w, "Failed to read file", http.StatusBadRequest)
		}
		defer file.Close()

		var banner types.Banner

		banner.Name = r.FormValue("name")
		banner.Description = r.FormValue("description")

		id, err := store.InsertBanner(r.URL.Query().Get("stories_id"), banner)
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
