package repository

import (
	"context"

	"github.com/tanya-mtv/metricsservice/internal/config"
	"github.com/tanya-mtv/metricsservice/internal/logger"
	"github.com/tanya-mtv/metricsservice/internal/models"
)

type Gauge float64
type Counter int64

type metricStorage interface {
	UpdateCounter(n string, v int64) Counter
	UpdateGauge(n string, v float64) Gauge
	GetAll() []models.Metrics
	GetCounter(metricName string) (Counter, bool)
	GetGauge(metricName string) (Gauge, bool)
}

type metricCollector interface {
	SetValueGauge(metricName string, value Gauge)
	SetValueCounter(metricName string, value Counter)
	GetAllCounter() map[string]Counter
	GetAllGauge() map[string]Gauge
}

type metricFiles interface {
	LoadLDataFromFile()
	SaveDataToFile(ctx context.Context)
}

type Storage struct {
	metricStorage
	metricFiles
}

type Collector struct {
	metricCollector
}

func NewStorage(repository *MetricStorage, cfg *config.ConfigServer, log logger.Logger) *Storage {
	return &Storage{
		metricFiles:   NewMetricMetricFiles(repository, cfg.FileName, cfg.Interval, log),
		metricStorage: NewMetricStorage(),
	}

}

func NewCollector() *Collector {
	return &Collector{
		metricCollector: NewMetricRepositoryCollector(),
	}

}
