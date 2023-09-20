package repository

import (
	"sync"

	"github.com/tanya-mtv/metricsservice/internal/models"
)

type MetricStorage struct {
	gaugeData    map[string]Gauge
	counterData  map[string]Counter
	countersLock sync.Mutex
	gaugesLock   sync.Mutex
}

func NewMetricStorage() *MetricStorage {

	return &MetricStorage{
		gaugeData:   make(map[string]Gauge),
		counterData: make(map[string]Counter),
	}
}

func (m *MetricStorage) UpdateCounter(n string, v int64) Counter {
	m.countersLock.Lock()
	defer m.countersLock.Unlock()

	m.counterData[n] += Counter(v)
	return m.counterData[n]
}

func (m *MetricStorage) UpdateGauge(n string, v float64) Gauge {
	m.gaugesLock.Lock()
	defer m.gaugesLock.Unlock()

	m.gaugeData[n] = Gauge(v)
	return m.gaugeData[n]
}

func (m *MetricStorage) GetAll() []models.Metrics {
	metricsSlice := make([]models.Metrics, 0, 29)

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
	return metricsSlice
}

func (m *MetricStorage) GetCounter(metricName string) (Counter, bool) {
	m.gaugesLock.Lock()
	defer m.gaugesLock.Unlock()

	res, found := m.counterData[metricName]
	return res, found
}

func (m *MetricStorage) GetGauge(metricName string) (Gauge, bool) {
	m.gaugesLock.Lock()
	defer m.gaugesLock.Unlock()

	res, found := m.gaugeData[metricName]
	return res, found
}

func (m *MetricStorage) UpdateMetrics(metrics []*models.Metrics) ([]*models.Metrics, error) {
	for _, value := range metrics {
		switch value.MType {
		case "counter":
			tmp := *value.Delta

			m.UpdateCounter(value.ID, tmp)

		case "gauge":

			tmp := *value.Value
			m.UpdateGauge(value.ID, tmp)

		}
	}
	return metrics, nil
}
