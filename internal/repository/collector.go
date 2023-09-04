package repository

import (
	"sync"
)

type MetricRepositoryCollector struct {
	gaugeData    map[string]Gauge
	counterData  map[string]Counter
	countersLock sync.Mutex
	gaugesLock   sync.Mutex
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

	return m.counterData

}

func (m *MetricRepositoryCollector) GetAllGauge() map[string]Gauge {
	return m.gaugeData
}
