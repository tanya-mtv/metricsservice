package repository

import (
	"errors"
	"sync"

	"github.com/tanya-mtv/metricsservice/internal/models"
)

type MetricRepositoryStorage struct {
	gaugeData    map[string]Gauge
	counterData  map[string]Counter
	countersLock sync.Mutex
	gaugesLock   sync.Mutex
}

func NewMetricRepositoryStorage() *MetricRepositoryStorage {

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
	metricsSlice := make([]models.Metrics, 0, 29)
	if len(m.counterData) == 0 && len(m.gaugeData) == 0 {
		return metricsSlice, errors.New("Storage is emty")
	}
	for name, value := range m.counterData {
		tmp := int64(value)
		data := models.Metrics{
			ID:    name,
			MType: "counter",
			Delta: &tmp,
		}
		metricsSlice = append(metricsSlice, data)
	}

	for name, value := range m.gaugeData {
		tmp := float64(value)
		data := models.Metrics{
			ID:    name,
			MType: "gauge",
			Value: &tmp,
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
