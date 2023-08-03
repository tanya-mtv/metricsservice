package repository

import (
	"github.com/tanya-mtv/metricsservice/internal/models"
	"github.com/tanya-mtv/metricsservice/internal/utils"
)

type MemStorage struct {
	// db *sqlx.DB
	gaugeData   map[string]utils.Gauge
	counterData map[string]utils.Counter
}

func NewMetricStorage() *MemStorage {
	// return &AuthPostgres{db: db}
	return &MemStorage{
		gaugeData:   make(map[string]utils.Gauge),
		counterData: make(map[string]utils.Counter),
	}
}

func (m *MemStorage) UpdateCounter(n string, v int64) utils.Counter {

	m.counterData[n] += utils.Counter(v)
	return m.counterData[n]

}

func (m *MemStorage) UpdateGauge(n string, v float64) utils.Gauge {
	m.gaugeData[n] = utils.Gauge(v)
	return m.gaugeData[n]
}

func (m *MemStorage) GetAll() ([]models.Metrics, error) {
	var metricsSlice []models.Metrics

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

func (m *MemStorage) GetCounter(metricName string) (utils.Counter, bool) {
	res, found := m.counterData[metricName]
	return res, found
}

func (m *MemStorage) GetGauge(metricName string) (utils.Gauge, bool) {

	res, found := m.gaugeData[metricName]
	return res, found
}
