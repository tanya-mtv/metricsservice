package models

type Metrics struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type" `           // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}

type MetricsP struct {
	ID    string   `json:"id" db:"name"`    // имя метрики
	MType string   `json:"Mtype" `          // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}
