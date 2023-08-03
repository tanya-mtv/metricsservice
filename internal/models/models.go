package models

type Metrics struct {
	ID         string  `json:"id"`                   // имя метрики
	MetricType string  `json:"type"`                 // параметр, принимающий значение gauge или counter
	CountValue int64   `json:"countvalue,omitempty"` // значение метрики в случае передачи counter
	GaugeValue float64 `json:"gaugevalue,omitempty"` // значение метрики в случае передачи gauge
}
