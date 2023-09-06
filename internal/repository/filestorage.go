package repository

import (
	"sync"

	"github.com/tanya-mtv/metricsservice/internal/models"
)

type FileStorage struct {
	gaugeData    map[string]Gauge
	counterData  map[string]Counter
	countersLock sync.Mutex
	gaugesLock   sync.Mutex
}

func NewMetricFiles() *FileStorage {

	return &FileStorage{
		gaugeData:   make(map[string]Gauge),
		counterData: make(map[string]Counter),
	}
}

func (m *FileStorage) UpdateCounter(n string, v int64) Counter {
	m.countersLock.Lock()
	defer m.countersLock.Unlock()

	m.counterData[n] += Counter(v)
	return m.counterData[n]
}

func (m *FileStorage) UpdateGauge(n string, v float64) Gauge {
	m.gaugesLock.Lock()
	defer m.gaugesLock.Unlock()

	m.gaugeData[n] = Gauge(v)
	return m.gaugeData[n]
}

func (m *FileStorage) GetAll() []models.Metrics {
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

func (m *FileStorage) GetCounter(metricName string) (Counter, bool) {
	m.gaugesLock.Lock()
	defer m.gaugesLock.Unlock()

	res, found := m.counterData[metricName]
	return res, found
}

func (m *FileStorage) GetGauge(metricName string) (Gauge, bool) {
	m.gaugesLock.Lock()
	defer m.gaugesLock.Unlock()

	res, found := m.gaugeData[metricName]
	return res, found
}

func (m *FileStorage) UpdateMetrics(metrics []models.Metrics) error {
	for _, value := range metrics {
		switch value.MType {
		case "counter":
			m.countersLock.Lock()
			defer m.countersLock.Unlock()
			tmp := *value.Delta

			m.counterData[value.ID] += Counter(tmp)
		case "gauge":
			m.gaugesLock.Lock()
			defer m.gaugesLock.Unlock()

			tmp := *value.Value

			m.gaugeData[value.ID] = Gauge(tmp)
		}
	}
	return nil
}
