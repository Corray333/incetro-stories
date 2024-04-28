package project

import (
	"github.com/Corray333/univer_cs/internal/domains/project/storage"
	"github.com/Corray333/univer_cs/internal/domains/project/transport"
	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
)

func Init(db *sqlx.DB, router *chi.Mux) error {
	store := storage.NewStorage(db)

	router.Group(func(subRouter chi.Router) {
		// subRouter.Use(auth.NewMiddleware())
		subRouter.Get("/api/projects", transport.GetProjects(store))
		subRouter.Get("/api/projects/{project_id}", transport.GetProjects(store))
		subRouter.Post("/api/projects", transport.NewProject(store))
	})

	return nil
}
