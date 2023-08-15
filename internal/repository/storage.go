package repository

import (
	"sync"

	"github.com/tanya-mtv/metricsservice/internal/models"
)

type MetricRepositoryStorage struct {
	gaugeData    map[string]Gauge
	counterData  map[string]Counter
	countersLock sync.Mutex
	gaugesLock   sync.Mutex
}

func NewMetricRepository() *MetricRepositoryStorage {

	return &MetricRepositoryStorage{
		gaugeData:   make(map[string]Gauge),
		counterData: make(map[string]Counter),
	}
}

func (m *MetricRepositoryStorage) UpdateCounter(n string, v int64) Counter {
	m.countersLock.Lock()
	defer m.countersLock.Unlock()

	m.counterData[n] += Counter(v)
	return m.counterData[n]
}

func (m *MetricRepositoryStorage) UpdateGauge(n string, v float64) Gauge {
	m.gaugesLock.Lock()
	defer m.gaugesLock.Unlock()

	m.gaugeData[n] = Gauge(v)
	return m.gaugeData[n]
}

func (m *MetricRepositoryStorage) GetAll() ([]models.Metrics, error) {
	metricsSlice := make([]models.Metrics, 0)

	for name, value := range m.counterData {

		data := models.Metrics{
			ID:         name,
			MetricType: "counter",
			CountValue: int64(value),
		}
		metricsSlice = append(metricsSlice, data)
	}

	for name, value := range m.gaugeData {

		data := models.Metrics{
			ID:         name,
			MetricType: "gauge",
			GaugeValue: float64(value),
		}
		metricsSlice = append(metricsSlice, data)
	}
	return metricsSlice, nil
}

func (m *MetricRepositoryStorage) GetCounter(metricName string) (Counter, bool) {
	m.gaugesLock.Lock()
	defer m.gaugesLock.Unlock()

	res, found := m.counterData[metricName]
	return res, found
}

func (m *MetricRepositoryStorage) GetGauge(metricName string) (Gauge, bool) {
	m.gaugesLock.Lock()
	defer m.gaugesLock.Unlock()

	res, found := m.gaugeData[metricName]
	return res, found
}
