package handlers

import (
	"fmt"
	"net/http"
	"path/filepath"
	"text/template"
	"time"

	"github.com/BDrimus/go-htmx-data-visualisation/internal/config"
	"github.com/BDrimus/go-htmx-data-visualisation/internal/models/timeseries"
)

var seriesCache = make(map[string]timeseries.TimeSeries)

type StockData struct {
	Symbol string
	Series timeseries.TimeSeries
}

func StockHandler(w http.ResponseWriter, r *http.Request) {

	stockData, err := getStockData(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tmpl := template.New(config.BaseTemplate).Funcs(funcMap)
	tmpl, err = tmpl.ParseFiles(
		filepath.Join(config.TemplatesFolder, config.BaseTemplate),
		filepath.Join(config.TemplatesFolder, "stock_component.html"),
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tmpl.ExecuteTemplate(w, "stock_component", stockData); err != nil {
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
		timeInterval = "1m"
	}

	timeIntervalDuration, err := parseTimeInterval(timeInterval)
	if err != nil {
		return nil, fmt.Errorf("invalid time interval")
	}

	series := checkSeriesCache(symbol, timeIntervalDuration)

	return &StockData{
		Symbol: symbol,
		Series: series,
	}, nil
}

func checkSeriesCache(symbol string, timeIntervalDuration time.Duration) (series timeseries.TimeSeries) {
	if cachedSeries, exists := seriesCache[symbol]; exists {
		series = cachedSeries
	} else {
		series = timeseries.GenerateTimeSeries(50, 100.0, 0, 0, timeIntervalDuration)
		seriesCache = make(map[string]timeseries.TimeSeries) // Clear the cache
		seriesCache[symbol] = series
	}

	return series
}
