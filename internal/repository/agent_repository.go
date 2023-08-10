package repository

import "github.com/tanya-mtv/metricsservice/internal/utils"

type MemStorageAgent struct {
	// db *sqlx.DB
	gaugeData   map[string]utils.Gauge
	counterData map[string]utils.Counter
}

func NewMetricStorageAgent() *MemStorageAgent {
	// return &AuthPostgres{db: db}
	return &MemStorageAgent{
		gaugeData:   make(map[string]utils.Gauge),
		counterData: make(map[string]utils.Counter),
	}
}

func (m *MemStorageAgent) SetGauge(metricName string, value utils.Gauge) {
	m.gaugeData[metricName] = value
}

func (m *MemStorageAgent) GetAllCounter() map[string]utils.Counter {
	return m.counterData
}

func (m *MemStorageAgent) GetAllGauge() map[string]utils.Gauge {
	return m.gaugeData
}
