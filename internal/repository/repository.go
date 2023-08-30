package repository

import (
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

type Storage struct {
	metricStorage
}

type Collector struct {
	metricCollector
}

func NewStorage() *Storage {
	return &Storage{
		metricStorage: NewMetricStorage(),
	}
}

func NewCollector() *Collector {
	return &Collector{
		metricCollector: NewMetricRepositoryCollector(),
	}

}
