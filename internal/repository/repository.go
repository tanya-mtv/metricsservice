package repository

import (
	"github.com/tanya-mtv/metricsservice/internal/logger"
	"github.com/tanya-mtv/metricsservice/internal/utils"
)

type MetricStorage interface {
	UpdateCounter(n string, v int64) utils.Counter
	UpdateGauge(n string, v float64) utils.Gauge
}

type Repository struct {
	MetricStorage
}

// func NewRepository(db *sql.DB, log logger.Logger) *Repository {
func NewRepository(log logger.Logger) *Repository {
	return &Repository{
		MetricStorage: NewMetricStorage(),
	}
}

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
