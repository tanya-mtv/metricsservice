package servise

import (
	"github.com/tanya-mtv/metricsservice/internal/repository"
	"github.com/tanya-mtv/metricsservice/internal/utils"
)

type StorageServiceAgent struct {
	repo repository.MetricStorageAgent
}

func NewMetricServiceAgent(repo repository.MetricStorageAgent) *StorageServiceAgent {
	return &StorageServiceAgent{repo: repo}
}

func (s *StorageServiceAgent) SetGauge(metricName string, value utils.Gauge) {
	s.repo.SetGauge(metricName, value)
}

func (s *StorageServiceAgent) GetAllCounter() map[string]utils.Counter {
	return s.repo.GetAllCounter()
}

func (s *StorageServiceAgent) GetAllGauge() map[string]utils.Gauge {
	return s.repo.GetAllGauge()
}
