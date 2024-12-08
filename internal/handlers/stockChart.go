package handlers

import (
	"encoding/json"
	"html/template"
	"net/http"
	"path/filepath"

	"github.com/BDrimus/go-htmx-data-visualisation/internal/config"
)

func StockChartHandler(w http.ResponseWriter, r *http.Request) {

	// Get primary stock data
	stockData := stockCache.Primary

	var compareStockData *StockData

	// Get comparison stock if specified
	if stockCache.Compare != nil {
		compareStockData = stockCache.Compare
	}

	data := struct {
		PrimaryStock template.JS
		CompareStock template.JS
	}{
		PrimaryStock: template.JS(mustMarshal(stockData)),
	}
	if compareStockData != nil {
		data.CompareStock = template.JS(mustMarshal(compareStockData))
	}

	tmpl := template.New(config.BaseTemplate).Funcs(funcMap)
	tmpl, err := tmpl.ParseFiles(
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

func mustMarshal(v interface{}) string {
	b, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return string(b)
}
