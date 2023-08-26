package metrics

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"runtime"
	"sync"
	"sync/atomic"

	"github.com/tanya-mtv/metricsservice/internal/logger"

	"github.com/tanya-mtv/metricsservice/internal/config"
	"github.com/tanya-mtv/metricsservice/internal/models"
	"github.com/tanya-mtv/metricsservice/internal/repository"
)

type counter struct {
	num int32
	sync.RWMutex
}

func (c *counter) inc() {
	c.Lock()
	defer c.Unlock()

	atomic.AddInt32(&c.num, 1)
}

func (c *counter) value() int32 {
	c.RLock()
	defer c.RUnlock()

	return c.num
}

func (c *counter) nulValue() {
	c.Lock()
	defer c.Unlock()

	c.num = 0
}

type ServiceMetrics struct {
	cfg               *config.ConfigAgent
	metricsRepository *repository.MetricRepositoryCollector
	counter           *counter
	httpClient        *http.Client
	buf               bytes.Buffer
	gzr               *gzip.Writer
}

func NewServiceMetrics(cfg *config.ConfigAgent, metricsRepository *repository.MetricRepositoryCollector) *ServiceMetrics {
	var bf bytes.Buffer
	gz, _ := gzip.NewWriterLevel(&bf, gzip.BestSpeed)
	return &ServiceMetrics{
		cfg:               cfg,
		metricsRepository: metricsRepository,
		counter: &counter{
			num: 0,
		},
		httpClient: &http.Client{},
		gzr:        gz,
		buf:        bf,
	}
}

func (sm *ServiceMetrics) MetricsMonitor() {

	var rtm runtime.MemStats

	sm.counter.inc()

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

	sm.metricsRepository.SetValueCounter("PollCount", repository.Counter(sm.counter.value()))
	sm.metricsRepository.SetValueGauge("RandomValue", repository.Gauge(float64(rand.Float64())))

}

func (sm *ServiceMetrics) Post(metric *models.Metrics, url string, log logger.Logger) (string, error) {

	data, err := json.Marshal(&metric)
	if err != nil {
		log.Debug("Can't post message")
		return "", err
	}

	err = sm.Compression(log, data)

	if err != nil {
		log.Info(err)
		return "", err
	}

	req, err := http.NewRequest("POST", url, bytes.NewReader(sm.buf.Bytes()))
	if err != nil {
		log.Error(err)
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Content-Encoding", "gzip")
	req.Header.Set("Accept-Encoding", "identity")
	resp, err := sm.httpClient.Do(req)

	if err != nil {
		log.Debug("Can't post message")
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	return string(body), err
}

func newMetric(metricName, metricsType string) *models.Metrics {

	return &models.Metrics{
		ID:    metricName,
		MType: metricsType,
	}
}

func (sm *ServiceMetrics) PostMessage(log logger.Logger) {
	addr := fmt.Sprintf("http://%s/update", sm.cfg.Port)

	for name, value := range sm.metricsRepository.GetAllGauge() {
		data := newMetric(name, "gauge")
		tmp := float64(value)
		data.Value = &tmp

		_, err := sm.Post(data, addr, log)

		if err != nil {
			log.Info(err)
		}

	}

	for name, value := range sm.metricsRepository.GetAllCounter() {

		data := newMetric(name, "counter")
		tmp := int64(value)
		data.Delta = &tmp

		_, err := sm.Post(data, addr, log)

		if err != nil {
			log.Info(err)
		}

	}

	sm.counter.nulValue()

}

func (sm *ServiceMetrics) Compression(log logger.Logger, b []byte) error {

	sm.buf.Reset()
	sm.gzr.Reset(&sm.buf)
	_, err := sm.gzr.Write(b)
	if err != nil {
		log.Debug(err)
		return err
	}
	sm.gzr.Close()

	return nil
}
