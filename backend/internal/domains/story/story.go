package story

import (
	"net/http"

	"github.com/Corray333/univer_cs/internal/domains/story/storage"
	"github.com/Corray333/univer_cs/internal/domains/story/transport"
	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
)

func Init(db *sqlx.DB, router *chi.Mux) error {
	store := storage.NewStorage(db)

	router.Group(func(subRouter chi.Router) {
		// subRouter.Use(auth.NewMiddleware())

		subRouter.Get("/api/projects/{project_id}/stories", transport.GetStories(store))
		subRouter.Post("/api/projects/{project_id}/banners", transport.NewBanner(store))
		subRouter.Post("/api/stories/views", transport.NewView(store))

		// subRouter.Post("/api/banners/{id}/media", transport.UpdateBannerMedia(store))
		subRouter.Put("/api/banners", transport.UpdateBanner(store))
	})
	fs := http.FileServer(http.Dir("../files/images"))
	router.Handle("/api/images/*", http.StripPrefix("/api/images", fs))

	return nil
}
