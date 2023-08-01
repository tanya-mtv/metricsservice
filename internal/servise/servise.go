package servise

import (
	"github.com/tanya-mtv/metricsservice/internal/repository"
	"github.com/tanya-mtv/metricsservice/internal/utils"
)

type MetricStorage interface {
	UpdateCounter(metricsName string, value int64) utils.Counter
	UpdateGauge(metricsName string, value float64) utils.Gauge
}
type Service struct {
	MetricStorage
}

func NewServise(repo *repository.Repository) *Service {
	return &Service{
		MetricStorage: NewMetricService(repo.MetricStorage),
	}
}

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
