package handlers

import (
	"log"
	"net/http"
	"path/filepath"
	"text/template"

	"github.com/BDrimus/go-htmx-data-visualisation/internal/config"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	baseFile := filepath.Join(config.TemplatesFolder, config.BaseTemplate)
	homeFile := filepath.Join(config.TemplatesFolder, "home.html")

	// Debug path resolution
	log.Printf("Loading templates from: %s and %s", baseFile, homeFile)

	tmpl := template.New(config.BaseTemplate).Funcs(funcMap)
	tmpl, err := tmpl.ParseFiles(baseFile, homeFile)
	if err != nil {
		log.Printf("Template parsing error: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tmpl.ExecuteTemplate(w, config.BaseTemplate, nil); err != nil {
		log.Printf("Template execution error: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
