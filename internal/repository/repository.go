package repository

type Gauge float64
type Counter int64

type metricCollector interface {
	SetValueGauge(metricName string, value Gauge)
	SetValueCounter(metricName string, value Counter)
	GetAllCounter() map[string]Counter
	GetAllGauge() map[string]Gauge
}

type Collector struct {
	metricCollector
}

func NewCollector() *Collector {
	return &Collector{
		metricCollector: NewMetricRepositoryCollector(),
	}

}
