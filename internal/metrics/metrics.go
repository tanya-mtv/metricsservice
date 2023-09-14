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
	"time"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/tanya-mtv/metricsservice/internal/constants"
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
	cfg        *config.ConfigAgent
	collector  metricCollector
	counter    *counter
	httpClient *retryablehttp.Client
	buf        bytes.Buffer
	gzr        *gzip.Writer
	log        logger.Logger
}

func NewServiceMetrics(collector *repository.MetricRepositoryCollector, cfg *config.ConfigAgent, log logger.Logger) *ServiceMetrics {
	var bf bytes.Buffer
	gz, _ := gzip.NewWriterLevel(&bf, gzip.BestSpeed)

	retryClient := retryablehttp.NewClient()
	retryClient.RetryMax = constants.RetryMax
	retryClient.RetryWaitMin = constants.RetryWaitMin
	retryClient.RetryWaitMax = constants.RetryWaitMax
	retryClient.Backoff = backoff

	return &ServiceMetrics{
		cfg:       cfg,
		collector: collector,
		counter: &counter{
			num: 0,
		},
		httpClient: retryClient,
		gzr:        gz,
		buf:        bf,
		log:        log,
	}
}

func (sm *ServiceMetrics) MetricsMonitor() {

	var rtm runtime.MemStats

	sm.counter.inc()

	runtime.ReadMemStats(&rtm)
	sm.collector.SetValueGauge("Alloc", repository.Gauge(rtm.Alloc))
	sm.collector.SetValueGauge("BuckHashSys", repository.Gauge(rtm.BuckHashSys))
	sm.collector.SetValueGauge("Frees", repository.Gauge(rtm.Frees))
	sm.collector.SetValueGauge("GCCPUFraction", repository.Gauge(rtm.GCCPUFraction))
	sm.collector.SetValueGauge("GCSys", repository.Gauge(rtm.GCSys))
	sm.collector.SetValueGauge("HeapAlloc", repository.Gauge(rtm.HeapAlloc))
	sm.collector.SetValueGauge("HeapIdle", repository.Gauge(rtm.HeapIdle))
	sm.collector.SetValueGauge("HeapInuse", repository.Gauge(rtm.HeapInuse))
	sm.collector.SetValueGauge("HeapObjects", repository.Gauge(rtm.HeapObjects))
	sm.collector.SetValueGauge("HeapReleased", repository.Gauge(rtm.HeapReleased))
	sm.collector.SetValueGauge("HeapSys", repository.Gauge(rtm.HeapSys))
	sm.collector.SetValueGauge("LastGC", repository.Gauge(rtm.LastGC))
	sm.collector.SetValueGauge("Lookups", repository.Gauge(rtm.Lookups))
	sm.collector.SetValueGauge("MCacheInuse", repository.Gauge(rtm.MCacheInuse))
	sm.collector.SetValueGauge("MCacheSys", repository.Gauge(rtm.MCacheSys))
	sm.collector.SetValueGauge("MSpanInuse", repository.Gauge(rtm.MSpanInuse))
	sm.collector.SetValueGauge("MSpanSys", repository.Gauge(rtm.MSpanSys))
	sm.collector.SetValueGauge("Mallocs", repository.Gauge(rtm.Mallocs))
	sm.collector.SetValueGauge("NextGC", repository.Gauge(rtm.NextGC))
	sm.collector.SetValueGauge("NumForcedGC", repository.Gauge(rtm.NumForcedGC))
	sm.collector.SetValueGauge("NumGC", repository.Gauge(rtm.NumGC))
	sm.collector.SetValueGauge("OtherSys", repository.Gauge(rtm.OtherSys))
	sm.collector.SetValueGauge("PauseTotalNs", repository.Gauge(rtm.PauseTotalNs))
	sm.collector.SetValueGauge("StackInuse", repository.Gauge(rtm.StackInuse))
	sm.collector.SetValueGauge("StackSys", repository.Gauge(rtm.StackSys))
	sm.collector.SetValueGauge("Sys", repository.Gauge(rtm.Sys))
	sm.collector.SetValueGauge("TotalAlloc", repository.Gauge(rtm.TotalAlloc))

	sm.collector.SetValueCounter("PollCount", repository.Counter(sm.counter.value()))
	sm.collector.SetValueGauge("RandomValue", repository.Gauge(float64(rand.Float64())))
}

func newMetric(metricName, metricsType string) *models.Metrics {

	return &models.Metrics{
		ID:    metricName,
		MType: metricsType,
	}
}
func (sm *ServiceMetrics) PostJSON(metrics []models.Metrics, url string) (string, error) {

	data, err := json.Marshal(&metrics)
	if err != nil {
		sm.log.Debug("Can't post message")
		return "", err
	}

	// err = sm.Compression(data)

	// if err != nil {
	// 	sm.log.Info(err)
	// 	return "", err
	// }

	// req, err := retryablehttp.NewRequest("POST", url, bytes.NewReader(sm.buf.Bytes()))

	req, err := retryablehttp.NewRequest("POST", url, bytes.NewReader(data))
	if err != nil {
		sm.log.Error(err)
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	// req.Header.Set("Content-Encoding", "gzip")
	req.Header.Set("Accept-Encoding", "identity")
	resp, err := sm.httpClient.Do(req)

	if err != nil {
		sm.log.Debug("Can't post message")
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	return string(body), err
}

func (sm *ServiceMetrics) PostMessageJSON() {
	addr := fmt.Sprintf("http://%s/updates/", sm.cfg.Port)
	// listMetrics := sm.collector.GetAllMetrics()

	listMetrics := make([]models.Metrics, 0, 29)

	for name, value := range sm.collector.GetAllGauge() {
		data := newMetric(name, "gauge")
		tmp := float64(value)
		data.Value = &tmp

		listMetrics = append(listMetrics, *data)

		fmt.Printf("listMetrics %+v\n", data)

	}

	for name, value := range sm.collector.GetAllCounter() {

		data := newMetric(name, "counter")

		tmp := int64(value)
		data.Delta = &tmp

		listMetrics = append(listMetrics, *data)

	}

	if len(listMetrics) > 0 {
		_, err := sm.PostJSON(listMetrics, addr)
		if err != nil {
			sm.log.Info(err)
		}
	}

	sm.counter.nulValue()

}

func (sm *ServiceMetrics) Compression(b []byte) error {

	sm.buf.Reset()
	sm.gzr.Reset(&sm.buf)
	_, err := sm.gzr.Write(b)
	if err != nil {
		sm.log.Debug(err)
		return err
	}
	sm.gzr.Close()

	return nil
}

func (sm *ServiceMetrics) Post(metric *models.Metrics, url string) (string, error) {

	data, err := json.Marshal(&metric)
	if err != nil {
		sm.log.Debug("Can't post message")
		return "", err
	}

	err = sm.Compression(data)

	if err != nil {
		sm.log.Info(err)
		return "", err
	}

	req, err := retryablehttp.NewRequest("POST", url, bytes.NewReader(sm.buf.Bytes()))
	if err != nil {
		sm.log.Error(err)
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Content-Encoding", "gzip")
	req.Header.Set("Accept-Encoding", "identity")
	resp, err := sm.httpClient.Do(req)

	if err != nil {
		sm.log.Debug("Can't post message")
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	return string(body), err
}

func (sm *ServiceMetrics) PostMessage() {
	addr := fmt.Sprintf("http://%s/update", sm.cfg.Port)

	for name, value := range sm.collector.GetAllGauge() {
		data := newMetric(name, "gauge")
		tmp := float64(value)
		data.Value = &tmp

		_, err := sm.Post(data, addr)

		if err != nil {
			sm.log.Info(err)
		}
	}

	for name, value := range sm.collector.GetAllCounter() {

		data := newMetric(name, "counter")
		tmp := int64(value)
		data.Delta = &tmp

		_, err := sm.Post(data, addr)

		if err != nil {
			sm.log.Info(err)
		}
	}
	sm.counter.nulValue()

}

func backoff(min, max time.Duration, attemptNum int, resp *http.Response) time.Duration {
	sleepTime := min + min*time.Duration(2*attemptNum)
	return sleepTime
}
