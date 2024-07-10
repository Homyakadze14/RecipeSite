package main

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/Homyakadze14/RecipeSite/RecipeSite/internal/common/middlewares"
	"github.com/Homyakadze14/RecipeSite/RecipeSite/internal/config"
	"github.com/Homyakadze14/RecipeSite/RecipeSite/internal/database"
	"github.com/gorilla/mux"
)

func main() {
	// Logger
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	// Config
	cfg, err := config.Parse()
	if err != nil {
		log.Fatal(err)
	}

	// Connect to database
	_, err = database.New(cfg)
	log.Printf("Connect to database on %s", cfg.DB_Host)
	if err != nil {
		log.Fatalf("Databse failed: %s", err)
	}

	// Handler
	handler := mux.NewRouter()
	handler.Use(middlewares.Logging)
	handler.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("Hello, world!")) })

	// Run server
	addr := fmt.Sprintf("%s:%v", cfg.Address, cfg.Port)
	server := http.Server{
		Addr:         addr,
		Handler:      handler,
		ReadTimeout:  20 * time.Second,
		WriteTimeout: 20 * time.Second,
	}
	log.Printf("Server start working on %s", addr)

	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Server failed: %s", err)
	}
}
