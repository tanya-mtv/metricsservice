package metrics

import (
	"bytes"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"runtime"
	"strconv"
	"time"

	"github.com/tanya-mtv/metricsservice/internal/config"
	"github.com/tanya-mtv/metricsservice/internal/repository"
)

var pollCount int64

type ServiceMetrics struct {
	cfg               *config.ConfigAgent
	metricsRepository *repository.MetricRepositoryCollector
}

func NewServiceMetrics(cfg *config.ConfigAgent, metricsRepository *repository.MetricRepositoryCollector) *ServiceMetrics {

	return &ServiceMetrics{
		cfg:               cfg,
		metricsRepository: metricsRepository,
	}
}

func (sm *ServiceMetrics) MetricsMonitor() {

	var rtm runtime.MemStats
	interval := time.Duration(sm.cfg.PollInterval) * time.Second
	for {
		time.Sleep(interval)
		pollCount += 1

		runtime.ReadMemStats(&rtm)
		sm.metricsRepository.SetValueGauge("Alloc", repository.Gauge(rtm.Alloc))
		sm.metricsRepository.SetValueGauge("BuckHashSys", repository.Gauge(rtm.BuckHashSys))
		sm.metricsRepository.SetValueGauge("Frees", repository.Gauge(rtm.Frees))
		sm.metricsRepository.SetValueGauge("GCCPUFraction", repository.Gauge(rtm.GCCPUFraction))
		sm.metricsRepository.SetValueGauge("GCSys", repository.Gauge(rtm.GCSys))
		sm.metricsRepository.SetValueGauge("HeapAlloc", repository.Gauge(rtm.HeapAlloc))
		sm.metricsRepository.SetValueGauge("HeapIdle", repository.Gauge(rtm.HeapIdle))
		sm.metricsRepository.SetValueGauge("HeapInuse", repository.Gauge(rtm.HeapInuse))
		sm.metricsRepository.SetValueGauge("HeapObjects", repository.Gauge(rtm.HeapObjects))
		sm.metricsRepository.SetValueGauge("HeapReleased", repository.Gauge(rtm.HeapReleased))
		sm.metricsRepository.SetValueGauge("HeapSys", repository.Gauge(rtm.HeapSys))
		sm.metricsRepository.SetValueGauge("LastGC", repository.Gauge(rtm.LastGC))
		sm.metricsRepository.SetValueGauge("Lookups", repository.Gauge(rtm.Lookups))
		sm.metricsRepository.SetValueGauge("MCacheInuse", repository.Gauge(rtm.MCacheInuse))
		sm.metricsRepository.SetValueGauge("MCacheSys", repository.Gauge(rtm.MCacheSys))
		sm.metricsRepository.SetValueGauge("MSpanInuse", repository.Gauge(rtm.MSpanInuse))
		sm.metricsRepository.SetValueGauge("MSpanSys", repository.Gauge(rtm.MSpanSys))
		sm.metricsRepository.SetValueGauge("Mallocs", repository.Gauge(rtm.Mallocs))
		sm.metricsRepository.SetValueGauge("NextGC", repository.Gauge(rtm.NextGC))
		sm.metricsRepository.SetValueGauge("NumForcedGC", repository.Gauge(rtm.NumForcedGC))
		sm.metricsRepository.SetValueGauge("NumGC", repository.Gauge(rtm.NumGC))
		sm.metricsRepository.SetValueGauge("OtherSys", repository.Gauge(rtm.OtherSys))
		sm.metricsRepository.SetValueGauge("PauseTotalNs", repository.Gauge(rtm.PauseTotalNs))
		sm.metricsRepository.SetValueGauge("StackInuse", repository.Gauge(rtm.StackInuse))
		sm.metricsRepository.SetValueGauge("StackSys", repository.Gauge(rtm.StackSys))
		sm.metricsRepository.SetValueGauge("Sys", repository.Gauge(rtm.Sys))
		sm.metricsRepository.SetValueGauge("TotalAlloc", repository.Gauge(rtm.TotalAlloc))

		sm.metricsRepository.SetValueCounter("pollCount", repository.Counter(pollCount))
		sm.metricsRepository.SetValueGauge("RandomValue", repository.Gauge(float64(rand.Float64())))

	}
}

func (sm *ServiceMetrics) Post(metricsType string, metricName string, metricValue string, url string) (string, error) {
	r := bytes.NewReader([]byte{})

	resp, err := http.Post(fmt.Sprintf("%s%s/%s/%s", url, metricsType, metricName, metricValue), "text/plain", r)

	if err != nil {
		fmt.Println("Can't post message", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	return string(body), err
}

func (sm *ServiceMetrics) PostMessage() {
	addr := fmt.Sprintf("http://%s/update/", sm.cfg.Port)

	for {

		for name, value := range sm.metricsRepository.GetAllGauge() {

			_, err := sm.Post("gauge", name, strconv.FormatFloat(float64(value), 'f', -1, 64), addr)

			if err != nil {
				fmt.Println("error reading body", err)
			}

		}

		for name, value := range sm.metricsRepository.GetAllCounter() {

			_, err := sm.Post("counter", name, strconv.FormatUint(uint64(value), 10), addr)

			if err != nil {
				fmt.Println("error reading body", err)
			}

		}

		pollCount = 0
		time.Sleep(time.Duration(sm.cfg.ReportInterval) * time.Second)
	}
}
