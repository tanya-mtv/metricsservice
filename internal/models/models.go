package models

type Metrics struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type" `           // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}

type MetricsP struct {
	ID    string   `json:"ID"`                 // имя метрики
	MType string   `json:"type",json:"MType" ` // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"Delta,omitempty"`    // значение метрики в случае передачи counter
	Value *float64 `json:"Value,omitempty"`    // значение метрики в случае передачи gauge
}
