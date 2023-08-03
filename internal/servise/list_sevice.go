package servise

import (
	"github.com/tanya-mtv/metricsservice/internal/models"
	"github.com/tanya-mtv/metricsservice/internal/repository"
	"github.com/tanya-mtv/metricsservice/internal/utils"
)

type StorageService struct {
	repo repository.MetricStorage
}

func NewMetricService(repo repository.MetricStorage) *StorageService {
	return &StorageService{repo: repo}
}

func (s *StorageService) UpdateCounter(metricsName string, value int64) utils.Counter {
	return s.repo.UpdateCounter(metricsName, value)
}

func (s *StorageService) UpdateGauge(metricsName string, value float64) utils.Gauge {
	return s.repo.UpdateGauge(metricsName, value)
}

func (s *StorageService) GetAll() ([]models.Metrics, error) {
	return s.repo.GetAll()
}

func (s *StorageService) GetCounter(metricName string) (utils.Counter, bool) {
	return s.repo.GetCounter(metricName)
}

func (s *StorageService) GetGauge(metricName string) (utils.Gauge, bool) {
	return s.repo.GetGauge(metricName)
}
