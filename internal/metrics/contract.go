package metrics

import (
	"github.com/tanya-mtv/metricsservice/internal/models"
	"github.com/tanya-mtv/metricsservice/internal/repository"
)

type metricCollector interface {
	SetValueGauge(metricName string, value repository.Gauge)
	SetValueCounter(metricName string, value repository.Counter)
	GetAllCounter() map[string]repository.Counter
	GetAllGauge() map[string]repository.Gauge
	GetAllMetrics() []models.Metrics
}
