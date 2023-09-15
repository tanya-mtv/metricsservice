package repository

import (
	"sync"

	"github.com/tanya-mtv/metricsservice/internal/models"
)

type MetricRepositoryCollector struct {
	gaugeData    map[string]Gauge
	counterData  map[string]Counter
	countersLock sync.RWMutex
	gaugesLock   sync.RWMutex
}

func NewMetricRepositoryCollector() *MetricRepositoryCollector {

	return &MetricRepositoryCollector{
		gaugeData:   make(map[string]Gauge),
		counterData: make(map[string]Counter),
	}
}

func (m *MetricRepositoryCollector) SetValueGauge(metricName string, value Gauge) {
	m.gaugesLock.Lock()
	defer m.gaugesLock.Unlock()

	m.gaugeData[metricName] = value
}

func (m *MetricRepositoryCollector) SetValueCounter(metricName string, value Counter) {
	m.countersLock.Lock()
	defer m.countersLock.Unlock()

	m.counterData[metricName] = value
}

func (m *MetricRepositoryCollector) GetAllCounter() map[string]Counter {
	m.countersLock.RLock()
	defer m.countersLock.RUnlock()

	data := make(map[string]Counter, len(m.counterData))

	for name, value := range m.counterData {
		data[name] = value
	}

	return data
}

func (m *MetricRepositoryCollector) GetAllGauge() map[string]Gauge {
	m.gaugesLock.RLock()
	defer m.gaugesLock.Unlock()

	data := make(map[string]Gauge, len(m.gaugeData))

	for name, value := range m.gaugeData {
		data[name] = value
	}
	return data
}

func (m *MetricRepositoryCollector) GetAllMetricsList() []models.Metrics {
	m.gaugesLock.RLock()
	defer m.gaugesLock.RUnlock()

	var listmetrics []models.Metrics
	for name, value := range m.gaugeData {
		tmp := float64(value)
		listmetrics = append(listmetrics, models.Metrics{ID: name, MType: "gauge", Value: &tmp})

	}

	m.countersLock.RLock()
	defer m.countersLock.RUnlock()

	for name, value := range m.counterData {
		tmp := int64(value)
		listmetrics = append(listmetrics, models.Metrics{ID: name, MType: "counter", Delta: &tmp})

	}
	return listmetrics
}
