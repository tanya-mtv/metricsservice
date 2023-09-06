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
	// m.countersLock.RLock()
	// defer m.countersLock.RUnlock()

	// data := m.counterData

	// for name, value := range data {
	// 	data[name] = value
	// }

	// return data
	return m.counterData

}

func (m *MetricRepositoryCollector) GetAllGauge() map[string]Gauge {
	// m.gaugesLock.RLock()
	// data := m.gaugeData
	// m.gaugesLock.RUnlock()

	// for name, value := range data {
	// 	data[name] = value
	// }
	// return data

	return m.gaugeData
}
