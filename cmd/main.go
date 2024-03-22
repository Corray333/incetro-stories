package main

import (
	"os"

	"github.com/Corray333/stories/internal/app"
	"github.com/Corray333/stories/internal/config"
)

func main() {
	// Load the configuration
	cfg := config.Configure(os.Args[1])

	app := app.NewApp(cfg)
	app.Run()

}
