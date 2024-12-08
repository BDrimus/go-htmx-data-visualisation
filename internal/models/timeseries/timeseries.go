package timeseries

import (
	"math"
	"math/rand"
	"time"
)

type DataPoint struct {
	Value float64
	Time  time.Time
}

type TimeSeries struct {
	Points     []DataPoint
	StartValue float64
	EndValue   float64
	Delta      float64
	MaxValue   float64
}

// GenerateTimeSeries generates a random time series with the given number of points and initial value.
func GenerateTimeSeries(numOfPoints int, initialValue, volatility, trend float64, timeInterval time.Duration) TimeSeries {

	// Assign default values if not provided
	if volatility == 0.0 {
		volatility = 0.1 // 10% max change
	}

	if trend == 0.0 {
		trend = 0.02 // 2% trend
	}

	// Random source with current time as seed
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	points := make([]DataPoint, numOfPoints)
	currentValue := initialValue
	startTime := time.Now()

	for i := 0; i < numOfPoints; i++ {

		// Random change
		change := (rng.Float64()*2 - 1) * volatility * currentValue

		// Trend
		trendFactor := trend * float64(i)
		currentValue += change + trendFactor

		timestamp := startTime.Add(timeInterval * time.Duration(i))

		points[i] = DataPoint{
			Value: math.Round(currentValue*100) / 100, // Round to 2 decimal places
			Time:  timestamp,
		}
	}

	series := TimeSeries{Points: points}

	// Calculate additional metadata
	series.StartValue = points[0].Value
	series.EndValue = points[len(points)-1].Value
	series.Delta = series.EndValue - series.StartValue
	series.MaxValue = series.max()

	return series
}

func (ts TimeSeries) max() float64 {
	if ts.MaxValue == 0 {
		for _, point := range ts.Points {
			if point.Value > ts.MaxValue {
				ts.MaxValue = point.Value
			}
		}
	}
	return ts.MaxValue
}
