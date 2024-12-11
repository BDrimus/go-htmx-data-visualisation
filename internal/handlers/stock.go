package handlers

import (
	"fmt"
	"net/http"
	"path/filepath"
	"text/template"

	"github.com/BDrimus/go-htmx-data-visualisation/internal/config"
	"github.com/BDrimus/go-htmx-data-visualisation/internal/models/timeseries"
)

type StockData struct {
	Symbol string
	Series *timeseries.TimeSeries
}

type StockTemplateData struct {
	StockData *StockData
	Type      string
}

// Handler for the primary stock data
func HandlePrimaryStock(w http.ResponseWriter, r *http.Request) {

	stockData, err := getStockData(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tmpl := initStockTemplate()

	data := StockTemplateData{
		StockData: stockData,
		Type:      "primary",
	}

	if err := tmpl.ExecuteTemplate(w, "stock_component", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// Handler for the comparison stock data
func HandleCompareStock(w http.ResponseWriter, r *http.Request) {

	stockData, err := getStockData(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tmpl := initStockTemplate()

	data := StockTemplateData{
		StockData: stockData,
		Type:      "compare",
	}

	if err := tmpl.ExecuteTemplate(w, "stock_info_panel", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

func getStockData(r *http.Request) (*StockData, error) {
	symbol := r.URL.Query().Get("symbol")
	if symbol == "" {
		return nil, fmt.Errorf("stock symbol is required")
	}

	timeInterval := r.URL.Query().Get("interval")
	if timeInterval == "" {
		timeInterval = config.DefaultTimeInterval
	}

	timeIntervalDuration, err := parseTimeInterval(timeInterval)
	if err != nil {
		return nil, fmt.Errorf("invalid time interval")
	}

	series := timeseries.GenerateTimeSeries(config.MaxDataPoints, config.InitialValue, config.Volatility, config.Trend, timeIntervalDuration)
	stockData := &StockData{
		Symbol: symbol,
		Series: &series,
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
