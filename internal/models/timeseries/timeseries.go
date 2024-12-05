package timeseries

import (
	"math"
	"math/rand"
	"time"
)

type DataPoint struct {
	Value float64
}

type TimeSeries struct {
	Points     []DataPoint
	StartValue float64
	EndValue   float64
	Delta      float64
	MaxValue   float64
}

// GenerateTimeSeries generates a random time series with the given number of points and initial value.
func GenerateTimeSeries(numOfPoints int, initialValue float64) TimeSeries {

	// Random source with current time as seed
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	points := make([]DataPoint, numOfPoints)
	currentValue := initialValue

	for i := 0; i < numOfPoints; i++ {

		const volatility = 0.1 /// 10% max change
		const trend = 0.02     // 2% trend

		// Random change
		change := (rng.Float64()*2 - 1) * volatility * currentValue

		// Trend
		trendFactor := trend * float64(i)
		currentValue += change + trendFactor

		points[i] = DataPoint{
			Value: math.Round(currentValue*100) / 100, // Round to 2 decimal places
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
