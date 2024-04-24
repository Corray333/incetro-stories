package main

import (
	"os"

	"github.com/Corray333/univer_cs/internal/app"
	"github.com/Corray333/univer_cs/internal/config"
)

func main() {
	// Load the configuration
	cfg := config.Configure(os.Args[1])

	app := app.NewApp(cfg)
	app.Run()

}
