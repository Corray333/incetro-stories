package story

import (
	"net/http"

	"github.com/Corray333/stories/internal/domains/story/storage"
	"github.com/Corray333/stories/internal/domains/story/transport"
	"github.com/Corray333/stories/internal/domains/story/types"
	"github.com/Corray333/stories/pkg/server/auth"
	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
)

type Storage interface {
	SelectStories(filter string) ([]types.Story, error)
	InsertStory(story types.Story) (int64, error)
	InsertBanner(storyId string, banner types.Banner) (int64, error)
}

func Init(db *sqlx.DB, router *chi.Mux) error {
	store, err := storage.NewStorage(db)
	if err != nil {
		return err
	}

	router.Group(func(subRouter chi.Router) {
		subRouter.Use(auth.NewMiddleware())

		subRouter.Get("/stories", transport.GetStories(store))
		subRouter.Post("/stories", transport.NewStories(store))
		// Replace with /stories/{id}
		subRouter.Post("/stories/banners", transport.NewBanner(store))
		subRouter.Post("/stories/views", transport.NewView(store))

		// TODO: replace with a proper file server
		fs := http.FileServer(http.Dir("../files/images"))
		subRouter.Handle("/images/*", http.StripPrefix("/images", fs))

	})

	return nil
}
