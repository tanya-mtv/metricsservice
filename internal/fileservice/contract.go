package fileservice

import (
	"github.com/tanya-mtv/metricsservice/internal/models"
	"github.com/tanya-mtv/metricsservice/internal/repository"
)

type fileStorage interface {
	UpdateCounter(n string, v int64) repository.Counter
	UpdateGauge(n string, v float64) repository.Gauge
	GetAll() []models.Metrics
}
