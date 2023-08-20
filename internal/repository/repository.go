package repository

import (
	"github.com/tanya-mtv/metricsservice/internal/models"
)

type Gauge float64
type Counter int64

type metricRepositoryStorage interface {
	UpdateCounter(n string, v int64) Counter
	UpdateGauge(n string, v float64) Gauge
	GetAll() ([]models.Metrics, error)
	GetCounter(metricName string) (Counter, bool)
	GetGauge(metricName string) (Gauge, bool)
}

type metricRepositoryCollector interface {
	SetValueGauge(metricName string, value Gauge)
	SetValueCounter(metricName string, value Counter)
	GetAllCounter() map[string]Counter
	GetAllGauge() map[string]Gauge
}

type RepositoryStorage struct {
	metricRepositoryStorage
}

type RepositoryCollector struct {
	metricRepositoryCollector
}
