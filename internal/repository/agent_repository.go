package repository

import "github.com/tanya-mtv/metricsservice/internal/utils"

type MemRepositoryAgent struct {
	gaugeData   map[string]utils.Gauge
	counterData map[string]utils.Counter
}

func NewMetricRepositoryAgent() *MemRepositoryAgent {

	return &MemRepositoryAgent{
		gaugeData:   make(map[string]utils.Gauge),
		counterData: make(map[string]utils.Counter),
	}
}

func (m *MemRepositoryAgent) SetValueGauge(metricName string, value utils.Gauge) {
	m.gaugeData[metricName] = value
}

func (m *MemRepositoryAgent) SetValueCounter(metricName string, value utils.Counter) {
	m.counterData[metricName] = value
}

func (m *MemRepositoryAgent) GetAllCounter() map[string]utils.Counter {
	return m.counterData
}

func (m *MemRepositoryAgent) GetAllGauge() map[string]utils.Gauge {
	return m.gaugeData
}
