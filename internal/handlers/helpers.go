package handlers

import (
	"encoding/json"
	"html/template"
	"strconv"
	"strings"
	"time"
)

// Template functions
var funcMap = template.FuncMap{
	"add":     func(a, b float64) float64 { return a + b },
	"mul":     func(a, b float64) float64 { return a * b },
	"sub":     func(a, b float64) float64 { return a - b },
	"div":     func(a, b float64) float64 { return a / b },
	"float64": func(x int) float64 { return float64(x) },
	"toJSON": func(v interface{}) template.HTMLAttr {
		jsonBytes, err := json.Marshal(v)
		if err != nil {
			return template.HTMLAttr("{}")
		}
		// Escape the JSON string properly
		jsonStr := string(jsonBytes)
		// Replace any double quotes in the JSON with &quot; to prevent HTML attribute issues
		jsonStr = strings.ReplaceAll(jsonStr, `"`, "&quot;")
		return template.HTMLAttr(jsonStr)
	},
}

func parseTimeInterval(interval string) (time.Duration, error) {
	// Handle day format
	if strings.HasSuffix(interval, "d") {
		days, err := strconv.Atoi(strings.TrimSuffix(interval, "d"))
		if err != nil {
			return 0, err
		}
		return time.Duration(days) * 24 * time.Hour, nil
	}

	return time.ParseDuration(interval)
}

func mustMarshal(v interface{}) string {
	b, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return string(b)
}
