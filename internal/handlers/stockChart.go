package handlers

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/BDrimus/go-htmx-data-visualisation/internal/config"
)

func StockChartHandler(w http.ResponseWriter, r *http.Request) {

	// Extract primary stock data from request
	primaryStockDataJSON := r.URL.Query().Get("primaryStockData")

	primaryStockData := StockData{}

	err := json.Unmarshal([]byte(primaryStockDataJSON), &primaryStockData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Extract secondary stock data from request
	compareStockDataJSON := r.URL.Query().Get("compareStockData")

	compareStockData := StockData{}

	// Check if single stock or compare as well
	if compareStockDataJSON != "" {
		err = json.Unmarshal([]byte(compareStockDataJSON), &compareStockData)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		log.Println("No compare stock data")
	}

	data := struct {
		PrimaryStock template.JS
		CompareStock template.JS
	}{
		PrimaryStock: template.JS(mustMarshal(primaryStockData)),
	}
	if compareStockData != (StockData{}) {
		data.CompareStock = template.JS(mustMarshal(compareStockData))
	}

	tmpl := template.New(config.BaseTemplate).Funcs(funcMap)
	tmpl, err = tmpl.ParseFiles(
		filepath.Join(config.TemplatesFolder, config.BaseTemplate),
		filepath.Join(config.TemplatesFolder, "stock_chart.html"),
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tmpl.ExecuteTemplate(w, "stock_chart", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
