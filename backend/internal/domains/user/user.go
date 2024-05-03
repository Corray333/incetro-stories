package user

import (
	"github.com/Corray333/univer_cs/internal/domains/user/storage"
	"github.com/Corray333/univer_cs/internal/domains/user/transport"
	"github.com/Corray333/univer_cs/pkg/server/auth"
	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
)

func Init(db *sqlx.DB, router *chi.Mux) error {
	store := storage.NewStorage(db)
	router.Post("/api/users/login", transport.LogIn(store))
	router.Post("/api/users/signup", transport.SignUp(store))
	router.Get("/api/users/refresh", transport.RefreshAccessToken(store))
	router.With(auth.NewMiddleware()).Put("/api/users/{id}", transport.UpdateUser(store))
	router.With(auth.NewMiddleware()).Get("/api/users/{id}", transport.GetUser(store))
	return nil
}
