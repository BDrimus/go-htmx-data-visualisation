package handlers

import (
	"fmt"
	"net/http"
	"path/filepath"
	"text/template"

	"github.com/BDrimus/go-htmx-data-visualisation/internal/config"
	"github.com/BDrimus/go-htmx-data-visualisation/internal/models/timeseries"
)

type StockPair struct {
	Primary *StockData
	Compare *StockData
}

var stockCache = &StockPair{
	Primary: nil,
	Compare: nil,
}

type StockData struct {
	Symbol string
	Series *timeseries.TimeSeries
}

// Handler for the primary stock data
func HandlePrimaryStock(w http.ResponseWriter, r *http.Request) {

	stockData, err := getStockData(r, false)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tmpl := initStockTemplate()

	if err := tmpl.ExecuteTemplate(w, "stock_component", stockData); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// Handler for the comparison stock data
func HandleCompareStock(w http.ResponseWriter, r *http.Request) {

	stockData, err := getStockData(r, true)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tmpl := initStockTemplate()

	if err := tmpl.ExecuteTemplate(w, "stock_info_panel", stockData); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

func getStockData(r *http.Request, isCompare bool) (*StockData, error) {
	symbol := r.URL.Query().Get("symbol")
	if symbol == "" {
		return nil, fmt.Errorf("stock symbol is required")
	}

	timeInterval := r.URL.Query().Get("interval")
	if timeInterval == "" {
		timeInterval = "1m"
	}

	timeIntervalDuration, err := parseTimeInterval(timeInterval)
	if err != nil {
		return nil, fmt.Errorf("invalid time interval")
	}

	series := timeseries.GenerateTimeSeries(50, 100.0, 0, 0, timeIntervalDuration)
	stockData := &StockData{
		Symbol: symbol,
		Series: &series,
	}

	if isCompare {
		stockCache.Compare = stockData
	} else {
		stockCache.Primary = stockData
	}

	return stockData, nil
}

func initStockTemplate() *template.Template {
	tmpl := template.New(config.BaseTemplate).Funcs(funcMap)
	tmpl, err := tmpl.ParseFiles(
		filepath.Join(config.TemplatesFolder, config.BaseTemplate),
		filepath.Join(config.TemplatesFolder, "stock_component.html"),
	)
	if err != nil {
		panic(err)
	}

	return tmpl
}
