package main

import (
	"log"

	"github.com/Homyakadze14/RecipeSite/config"
	"github.com/Homyakadze14/RecipeSite/internal/app"
)

func main() {
	// Configuration
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	// Run
	app.Run(cfg)
}
