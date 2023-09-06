package handler

import (
	"github.com/tanya-mtv/metricsservice/internal/models"
	"github.com/tanya-mtv/metricsservice/internal/repository"
)

type metricStorage interface {
	UpdateCounter(n string, v int64) repository.Counter
	UpdateGauge(n string, v float64) repository.Gauge
	GetAll() []models.Metrics
	GetCounter(metricName string) (repository.Counter, bool)
	GetGauge(metricName string) (repository.Gauge, bool)
	UpdateMetrics([]models.Metrics) error
}

type pingDB interface {
	Ping() error
}
