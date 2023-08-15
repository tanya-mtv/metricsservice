package models

type Metrics struct {
	ID         string
	MetricType string
	CountValue int64
	GaugeValue float64
}
