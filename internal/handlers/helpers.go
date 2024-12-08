package handlers

import (
	"strconv"
	"strings"
	"text/template"
	"time"
)

// Template functions
var funcMap = template.FuncMap{
	"add":     func(a, b float64) float64 { return a + b },
	"mul":     func(a, b float64) float64 { return a * b },
	"sub":     func(a, b float64) float64 { return a - b },
	"div":     func(a, b float64) float64 { return a / b },
	"float64": func(x int) float64 { return float64(x) },
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
