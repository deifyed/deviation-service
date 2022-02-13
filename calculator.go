package main

import "math"

type movingAverageCalculator struct {
	Mean     float64
	dsquared float64
	count    int
}

func (m *movingAverageCalculator) Update(newValue int64) {
	m.count++

	val := float64(newValue)

	diff := (val - m.Mean) / float64(m.count)
	newMean := m.Mean + diff

	dsquaredIncrement := (val - newMean) * (val - m.Mean)

	m.Mean = newMean
	m.dsquared += dsquaredIncrement
}

func (m *movingAverageCalculator) Variance() float64 {
	return m.dsquared / float64(m.count)
}

func (m *movingAverageCalculator) StandardDeviation() float64 {
	return math.Sqrt(m.Variance())
}
