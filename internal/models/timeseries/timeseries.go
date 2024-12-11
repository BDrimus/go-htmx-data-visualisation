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
	Points                         []DataPoint
	StartValue                     float64
	EndValue                       float64
	Delta                          float64
	MaxValue                       float64
	MinValue                       float64
	PricePercentageChangeFromStart float64
}

// GenerateTimeSeries generates a random time series with the given number of points and initial value.
func GenerateTimeSeries(numOfPoints int, initialValue, volatility, trend float64, timeInterval time.Duration) TimeSeries {

	// Assign default value if not provided
	if timeInterval == 0 {
		timeInterval = time.Minute
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
	series.MinValue = series.min()
	series.PricePercentageChangeFromStart = series.pricePercentageChangeFromStart()

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

func (ts TimeSeries) min() float64 {

	var smallestValue float64 = ts.Points[0].Value

	for _, point := range ts.Points {
		if point.Value < smallestValue {
			smallestValue = point.Value
		}
	}

	return smallestValue
}

// pricePercentageChangeFromStart returns the percentage change from the start value to the end value
func (ts TimeSeries) pricePercentageChangeFromStart() float64 {
	if ts.EndValue > ts.StartValue {
		return (ts.EndValue - ts.StartValue) / (ts.MaxValue - ts.StartValue) * 100
	} else {
		return (ts.StartValue - ts.EndValue) / (ts.StartValue - ts.MinValue) * -100
	}
}
