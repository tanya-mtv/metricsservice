package repository

import (
	"sync"

	"github.com/tanya-mtv/metricsservice/internal/models"
)

type MetricRepositoryCollector struct {
	gaugeData   map[string]Gauge
	counterData map[string]Counter
	mu          sync.RWMutex
}

func NewMetricRepositoryCollector() *MetricRepositoryCollector {

	return &MetricRepositoryCollector{
		gaugeData:   make(map[string]Gauge),
		counterData: make(map[string]Counter),
	}
}

func (m *MetricRepositoryCollector) SetValueGauge(metricName string, value Gauge) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.gaugeData[metricName] = value
}

func (m *MetricRepositoryCollector) SetValueCounter(metricName string, value Counter) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.counterData[metricName] = value
}

func (m *MetricRepositoryCollector) GetAllCounter() map[string]Counter {
	m.mu.RLock()
	defer m.mu.RUnlock()

	data := make(map[string]Counter, len(m.counterData))

	for name, value := range m.counterData {
		data[name] = value
	}

	return data
}

func (m *MetricRepositoryCollector) GetAllGauge() map[string]Gauge {
	m.mu.RLock()
	defer m.mu.RUnlock()

	data := make(map[string]Gauge, len(m.gaugeData))

	for name, value := range m.gaugeData {
		data[name] = value
	}
	return data
}

func (m *MetricRepositoryCollector) GetAllMetricsList() []models.Metrics {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var listmetrics []models.Metrics
	k := 0
	for name, value := range m.gaugeData {
		k += 1
		tmp := float64(value)
		listmetrics = append(listmetrics, models.Metrics{ID: name, MType: "gauge", Value: &tmp})
		if k == 1 {
			break
		}

	}

	for name, value := range m.counterData {
		tmp := int64(value)
		listmetrics = append(listmetrics, models.Metrics{ID: name, MType: "counter", Delta: &tmp})

	}
	return listmetrics
}
