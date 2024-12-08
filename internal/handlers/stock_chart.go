package handlers

import (
	"encoding/json"
	"html/template"
	"net/http"
	"path/filepath"

	"github.com/BDrimus/go-htmx-data-visualisation/internal/config"
)

func StockChartHandler(w http.ResponseWriter, r *http.Request) {

	stockData, err := getStockData(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// In StockChartHandler, before template execution:
	jsonData, err := json.Marshal(stockData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Pass the JSON string to template
	data := struct {
		StockDataJSON template.JS
	}{
		StockDataJSON: template.JS(jsonData),
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
