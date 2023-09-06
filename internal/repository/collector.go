package repository

import (
	"sync"
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
	data := m.counterData
	m.countersLock.RUnlock()

	for name, value := range data {
		data[name] = value
	}

	return data

}

func (m *MetricRepositoryCollector) GetAllGauge() map[string]Gauge {
	m.gaugesLock.RLock()
	data := m.gaugeData
	m.gaugesLock.RUnlock()

	for name, value := range data {
		data[name] = value
	}

	return m.gaugeData
}
