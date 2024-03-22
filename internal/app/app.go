package app

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/Corray333/stories/internal/config"
	"github.com/Corray333/stories/internal/domains/story"
	"github.com/Corray333/stories/internal/domains/user"
	"github.com/Corray333/stories/internal/storage"
	"github.com/Corray333/stories/pkg/server/logger"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/jmoiron/sqlx"
)

type App struct {
	db     *sqlx.DB
	server *http.Server
}

func NewApp(cfg *config.Config) *App {
	db, err := storage.Connect()
	if err != nil {
		slog.Error("Failed to connect to the database: " + err.Error())
		panic(err)
	}

	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins: cfg.AllowedOrigins,
		AllowedMethods: cfg.AllowedMethods,
		AllowedHeaders: []string{"Authorization"},
		MaxAge:         300,
	}))

	if cfg.Env == config.EnvDev {
		router.Use(middleware.RequestID)
		router.Use(logger.New(slog.Default()))
	}

	if err := user.Init(db, router); err != nil {
		slog.Error(err.Error())
		panic(err)

	}
	if err := story.Init(db, router); err != nil {
		slog.Error(err.Error())
		panic(err)

	}

	return &App{
		db: db,
		server: &http.Server{
			Addr:    os.Getenv("APP_IP") + ":" + os.Getenv("APP_PORT"),
			Handler: router,
		},
	}
}

func (app *App) Run() {
	slog.Info("Server started on port " + os.Getenv("APP_PORT"))
	if err := app.server.ListenAndServe(); err != nil {
		slog.Error("Server failed to start: " + err.Error())
		panic(err)
	}
}
