package timeseries

import (
	"testing"
	"time"
)

func TestGenerateTimeSeries(t *testing.T) {
	t.Run("generates correct number of points", func(t *testing.T) {
		numPoints := 10
		series := GenerateTimeSeries(numPoints, 100.0, 0.1, 0.02, time.Minute)

		if len(series.Points) != numPoints {
			t.Errorf("Expected %d points, got %d", numPoints, len(series.Points))
		}
	})

	t.Run("timestamps are sequential with correct interval", func(t *testing.T) {
		interval := time.Minute
		series := GenerateTimeSeries(5, 100.0, 0.1, 0.02, interval)

		for i := 1; i < len(series.Points); i++ {
			expected := series.Points[i-1].Time.Add(interval)
			if !series.Points[i].Time.Equal(expected) {
				t.Errorf("Expected timestamp %v, got %v", expected, series.Points[i].Time)
			}
		}
	})

	t.Run("metadata is correctly calculated", func(t *testing.T) {
		series := GenerateTimeSeries(10, 100.0, 0.1, 0.02, time.Minute)

		if series.StartValue != series.Points[0].Value {
			t.Errorf("StartValue doesn't match first point value")
		}

		if series.EndValue != series.Points[len(series.Points)-1].Value {
			t.Errorf("EndValue doesn't match last point value")
		}

		if series.Delta != (series.EndValue - series.StartValue) {
			t.Errorf("Delta calculation incorrect")
		}
	})

	t.Run("min and max values are correct", func(t *testing.T) {
		series := GenerateTimeSeries(10, 100.0, 0.1, 0.02, time.Minute)

		// Check if MaxValue is actually the maximum
		for _, point := range series.Points {
			if point.Value > series.MaxValue {
				t.Errorf("Found value %f greater than MaxValue %f", point.Value, series.MaxValue)
			}
		}

		// Check if MinValue is actually the minimum
		for _, point := range series.Points {
			if point.Value < series.MinValue {
				t.Errorf("Found value %f less than MinValue %f", point.Value, series.MinValue)
			}
		}
	})
}

func TestPricePercentageChangeFromStart(t *testing.T) {
	t.Run("percentage change calculation is within bounds", func(t *testing.T) {
		series := GenerateTimeSeries(10, 100.0, 0.1, 0.02, time.Minute)

		if series.PricePercentageChangeFromStart < -100 || series.PricePercentageChangeFromStart > 100 {
			t.Errorf("Percentage change %f is outside valid range [-100, 100]", series.PricePercentageChangeFromStart)
		}
	})
}
