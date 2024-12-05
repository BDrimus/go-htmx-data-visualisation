package main

import (
	"log"
	"net/http"

	"github.com/BDrimus/go-htmx-data-visualisation/internal/config"
	"github.com/BDrimus/go-htmx-data-visualisation/internal/handlers"
)

func main() {
	mux := http.NewServeMux()

	// Routes
	mux.HandleFunc("GET /", handlers.TimeSeriesHandler)
	mux.HandleFunc("GET /update", handlers.TimeSeriesHandler)

	log.Printf("Starting server on %s", config.PORT)
	log.Fatal(http.ListenAndServe(config.PORT, mux))
}
