package repository

import (
	"github.com/tanya-mtv/metricsservice/internal/logger"
	"github.com/tanya-mtv/metricsservice/internal/models"
	"github.com/tanya-mtv/metricsservice/internal/utils"
)

type MetricStorage interface {
	UpdateCounter(n string, v int64) utils.Counter
	UpdateGauge(n string, v float64) utils.Gauge
	GetAll() ([]models.Metrics, error)
	GetCounter(metricName string) (utils.Counter, bool)
	GetGauge(metricName string) (utils.Gauge, bool)
}

type MetricStorageAgent interface {
	SetGauge(metricName string, value utils.Gauge)
	GetAllCounter() map[string]utils.Counter
	GetAllGauge() map[string]utils.Gauge
}

type Repository struct {
	MetricStorage
	MetricStorageAgent
}

// func NewRepository(db *sql.DB, log logger.Logger) *Repository {
func NewRepository(log logger.Logger) *Repository {
	return &Repository{
		MetricStorage:      NewMetricStorage(),
		MetricStorageAgent: NewMetricStorageAgent(),
	}
}
