package servise

import (
	"github.com/tanya-mtv/metricsservice/internal/models"
	"github.com/tanya-mtv/metricsservice/internal/repository"
	"github.com/tanya-mtv/metricsservice/internal/utils"
)

type MetricStorage interface {
	UpdateCounter(metricsName string, value int64) utils.Counter
	UpdateGauge(metricsName string, value float64) utils.Gauge
	GetAll() ([]models.Metrics, error)
	GetCounter(metricName string) (utils.Counter, bool)
	GetGauge(metricName string) (utils.Gauge, bool)
}
type Service struct {
	MetricStorage
}

func NewServise(repo *repository.Repository) *Service {
	return &Service{
		MetricStorage: NewMetricService(repo.MetricStorage),
	}
}
