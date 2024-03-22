package storage

import (
	"log"
	"log/slog"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func Connect() (*sqlx.DB, error) {
	db, err := sqlx.Connect("postgres", os.Getenv("POSTGRES_CONNSTR"))
	if err != nil {
		log.Fatal(err)
	}

	// Test the connection to the database
	if err := db.Ping(); err != nil {
		return nil, err
	} else {
		slog.Info("Successfully Connected")
	}

	return db, nil
}
