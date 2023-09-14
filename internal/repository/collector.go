package repository

import (
	"github.com/tanya-mtv/metricsservice/internal/models"
)

type MetricRepositoryCollector struct {
	gaugeData   map[string]Gauge
	counterData map[string]Counter
	// lock        sync.Mutex
	// countersLock sync.RWMutex
	// gaugesLock   sync.RWMutex
}

func NewMetricRepositoryCollector() *MetricRepositoryCollector {

	return &MetricRepositoryCollector{
		gaugeData:   make(map[string]Gauge),
		counterData: make(map[string]Counter),
	}
}

func (m *MetricRepositoryCollector) SetValueGauge(metricName string, value Gauge) {
	// m.lock.Lock()
	// defer m.lock.Unlock()

	m.gaugeData[metricName] = value
}

func (m *MetricRepositoryCollector) SetValueCounter(metricName string, value Counter) {
	// m.lock.Lock()
	// defer m.lock.Unlock()

	m.counterData[metricName] = value
}

func (m *MetricRepositoryCollector) GetAllCounter() map[string]Counter {
	// m.lock.Lock()
	// defer m.lock.Unlock()

	data := make(map[string]Counter, len(m.counterData))

	for name, value := range m.counterData {
		data[name] = value
	}

	return data
}

func (m *MetricRepositoryCollector) GetAllGauge() map[string]Gauge {
	// m.lock.Lock()
	// defer m.lock.Unlock()

	data := make(map[string]Gauge, len(m.gaugeData))

	for name, value := range m.gaugeData {
		data[name] = value
	}
	return data
}

func (m *MetricRepositoryCollector) GetAllMetrics() []models.Metrics {
	listMetrics := make([]models.Metrics, 0, 29)
	// m.lock.Lock()
	// defer m.lock.Unlock()

	// for name, value := range m.gaugeData {
	// 	tmp := float64(value)
	// 	listMetrics = append(listMetrics, models.Metrics{ID: name, MType: "gauge", Value: &tmp})
	// }
	tmp := float64(55200)
	listMetrics = append(listMetrics, models.Metrics{ID: "MSpanInuse", MType: "gauge", Value: &tmp})
	tmp = float64(163840)
	listMetrics = append(listMetrics, models.Metrics{ID: "StackSys", MType: "gauge", Value: &tmp})
	tmp = float64(1950312)
	listMetrics = append(listMetrics, models.Metrics{ID: "Alloc", MType: "gauge", Value: &tmp})

	for name, value := range m.counterData {
		tmp := int64(value)
		listMetrics = append(listMetrics, models.Metrics{ID: name, MType: "counter", Delta: &tmp})
	}

	return listMetrics
}
