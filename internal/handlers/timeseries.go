package handlers

import (
	"net/http"
	"path/filepath"
	"text/template"

	"github.com/BDrimus/go-htmx-data-visualisation/internal/config"
	"github.com/BDrimus/go-htmx-data-visualisation/internal/models/timeseries"
)

// Template functions
var funcMap = template.FuncMap{
	"add":     func(a, b float64) float64 { return a + b },
	"mul":     func(a, b float64) float64 { return a * b },
	"sub":     func(a, b float64) float64 { return a - b },
	"div":     func(a, b float64) float64 { return a / b },
	"float64": func(x int) float64 { return float64(x) },
}

func initTemplate() (*template.Template, error) {
	// Create template with functions
	tmpl := template.New(config.BaseTemplate).Funcs(funcMap)
	return tmpl.ParseFiles(
		filepath.Join(config.TemplatesFolder, config.BaseTemplate),
		filepath.Join(config.TemplatesFolder, "timeseries.html"),
		filepath.Join(config.TemplatesFolder, "timeseries_info_panel.html"),
		filepath.Join(config.TemplatesFolder, "svg_chart.html"),
	)
}

func TimeSeriesHandler(w http.ResponseWriter, r *http.Request) {
	series := timeseries.GenerateTimeSeries(50, 100.0, 0, 0)
	tmpl, err := initTemplate()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	target := r.Header.Get("HX-Target")
	var templateName string

	switch target {
	case "info-panel-container":
		templateName = "timeseries_info_panel"
	case "chart-container":
		templateName = "svg_chart"
	case "data-container":
		templateName = "data_container"
	default:
		templateName = config.BaseTemplate
	}

	if err := tmpl.ExecuteTemplate(w, templateName, series); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
