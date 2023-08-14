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

	"github.com/tanya-mtv/metricsservice/internal/utils"
)

var pollCount int64

type ServiceMetrics struct {
	cfg               *config.ConfigAgent
	metricsRepository *repository.Repository
}

func NewServiceMetrics(cfg *config.ConfigAgent, metricsRepository *repository.Repository) *ServiceMetrics {

	return &ServiceMetrics{
		cfg:               cfg,
		metricsRepository: metricsRepository,
	}
}

func (sm *ServiceMetrics) NewMonitor() {

	var rtm runtime.MemStats
	interval := time.Duration(sm.cfg.PollInterval) * time.Second
	for {
		time.Sleep(interval)
		pollCount += 1

		runtime.ReadMemStats(&rtm)
		sm.metricsRepository.SetValueGauge("Alloc", utils.Gauge(rtm.Alloc))
		sm.metricsRepository.SetValueGauge("BuckHashSys", utils.Gauge(rtm.BuckHashSys))
		sm.metricsRepository.SetValueGauge("Frees", utils.Gauge(rtm.Frees))
		sm.metricsRepository.SetValueGauge("GCCPUFraction", utils.Gauge(rtm.GCCPUFraction))
		sm.metricsRepository.SetValueGauge("GCSys", utils.Gauge(rtm.GCSys))
		sm.metricsRepository.SetValueGauge("HeapAlloc", utils.Gauge(rtm.HeapAlloc))
		sm.metricsRepository.SetValueGauge("HeapIdle", utils.Gauge(rtm.HeapIdle))
		sm.metricsRepository.SetValueGauge("HeapInuse", utils.Gauge(rtm.HeapInuse))
		sm.metricsRepository.SetValueGauge("HeapObjects", utils.Gauge(rtm.HeapObjects))
		sm.metricsRepository.SetValueGauge("HeapReleased", utils.Gauge(rtm.HeapReleased))
		sm.metricsRepository.SetValueGauge("HeapSys", utils.Gauge(rtm.HeapSys))
		sm.metricsRepository.SetValueGauge("LastGC", utils.Gauge(rtm.LastGC))
		sm.metricsRepository.SetValueGauge("Lookups", utils.Gauge(rtm.Lookups))
		sm.metricsRepository.SetValueGauge("MCacheInuse", utils.Gauge(rtm.MCacheInuse))
		sm.metricsRepository.SetValueGauge("MCacheSys", utils.Gauge(rtm.MCacheSys))
		sm.metricsRepository.SetValueGauge("MSpanInuse", utils.Gauge(rtm.MSpanInuse))
		sm.metricsRepository.SetValueGauge("MSpanSys", utils.Gauge(rtm.MSpanSys))
		sm.metricsRepository.SetValueGauge("Mallocs", utils.Gauge(rtm.Mallocs))
		sm.metricsRepository.SetValueGauge("NextGC", utils.Gauge(rtm.NextGC))
		sm.metricsRepository.SetValueGauge("NumForcedGC", utils.Gauge(rtm.NumForcedGC))
		sm.metricsRepository.SetValueGauge("NumGC", utils.Gauge(rtm.NumGC))
		sm.metricsRepository.SetValueGauge("OtherSys", utils.Gauge(rtm.OtherSys))
		sm.metricsRepository.SetValueGauge("PauseTotalNs", utils.Gauge(rtm.PauseTotalNs))
		sm.metricsRepository.SetValueGauge("StackInuse", utils.Gauge(rtm.StackInuse))
		sm.metricsRepository.SetValueGauge("StackSys", utils.Gauge(rtm.StackSys))
		sm.metricsRepository.SetValueGauge("Sys", utils.Gauge(rtm.Sys))
		sm.metricsRepository.SetValueGauge("TotalAlloc", utils.Gauge(rtm.TotalAlloc))

		sm.metricsRepository.SetValueCounter("pollCount", utils.Counter(pollCount))
		sm.metricsRepository.SetValueGauge("RandomValue", utils.Gauge(float64(rand.Float64())))

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
