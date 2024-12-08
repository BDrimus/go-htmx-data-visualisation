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
	fs := http.FileServer(http.Dir("static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	mux.HandleFunc("/", handlers.HomeHandler)
	mux.HandleFunc("/stock", handlers.StockHandler)
	mux.HandleFunc("/stock_chart", handlers.StockChartHandler)

	log.Printf("Starting server on http://localhost%s/", config.PORT)
	log.Fatal(http.ListenAndServe(config.PORT, mux))
}
