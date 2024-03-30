package story

import (
	"net/http"

	"github.com/Corray333/univer_cs/internal/domains/story/storage"
	"github.com/Corray333/univer_cs/internal/domains/story/transport"
	"github.com/Corray333/univer_cs/pkg/server/auth"
	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
)

func Init(db *sqlx.DB, router *chi.Mux) error {
	store, err := storage.NewStorage(db)
	if err != nil {
		return err
	}

	router.Group(func(subRouter chi.Router) {
		subRouter.Use(auth.NewMiddleware())

		subRouter.Get("/stories", transport.GetStories(store))
		subRouter.Post("/banners", transport.NewBanner(store))
		subRouter.Post("/stories/views", transport.NewView(store))

		subRouter.Post("/banners/{id}/media", transport.UpdateBannerMedia(store))
		subRouter.Post("/story/{id}/timestamp", transport.UpdateStoryTimestamp(store))
		subRouter.Post("/banners/{id}/name", transport.UpdateBannerName(store))
		subRouter.Post("/banners/{id}/description", transport.UpdateBannerDescription(store))
	})
	fs := http.FileServer(http.Dir("../files/images"))
	router.Handle("/images/*", http.StripPrefix("/images", fs))

	return nil
}
