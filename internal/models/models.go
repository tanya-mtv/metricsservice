package models

type Metrics struct {
	ID         string  `json:"id"`
	MetricType string  `json:"type"`
	CountValue int64   `json:"countvalue,omitempty"`
	GaugeValue float64 `json:"gaugevalue,omitempty"`
}
