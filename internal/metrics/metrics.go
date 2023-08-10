package metrics

import (
	"bytes"
	"fmt"
	"math/rand"
	"net/http"
	"reflect"
	"runtime"
	"strconv"
	"time"

	"github.com/tanya-mtv/metricsservice/internal/config"
	"github.com/tanya-mtv/metricsservice/internal/servise"
	"github.com/tanya-mtv/metricsservice/internal/utils"
)

// var valuesGauge = map[string]float64{}
var pollCount uint64

var metrics = map[string]bool{
	"Alloc":         true,
	"BuckHashSys":   true,
	"Frees":         true,
	"GCCPUFraction": true,
	"GCSys":         true,
	"HeapAlloc":     true,
	"HeapIdle":      true,
	"HeapInuse":     true,
	"HeapObjects":   true,
	"HeapReleased":  true,
	"HeapSys":       true,
	"LastGC":        true,
	"Lookups":       true,
	"MCacheInuse":   true,
	"MCacheSys":     true,
	"MSpanInuse":    true,
	"MSpanSys":      true,
	"Mallocs":       true,
	"NextGC":        true,
	"NumForcedGC":   true,
	"NumGC":         true,
	"OtherSys":      true,
	"PauseTotalNs":  true,
	"StackInuse":    true,
	"StackSys":      true,
	"Sys":           true,
	"TotalAlloc":    true,
	"PollCount":     true,
	"RandomValue":   true,
}

type ServiceMetrics struct {
	cfg     *config.ConfigAgent
	metrics *servise.Service
}

func NewServiceMetrics(cfg *config.ConfigAgent, metrics *servise.Service) *ServiceMetrics {
	return &ServiceMetrics{
		cfg:     cfg,
		metrics: metrics,
	}
}

func (sm *ServiceMetrics) NewMonitor() {

	var rtm runtime.MemStats
	var interval = time.Duration(sm.cfg.PollInterval) * time.Second
	for {
		<-time.After(interval)
		pollCount += 1

		runtime.ReadMemStats(&rtm)

		v := reflect.ValueOf(rtm)
		typeOfS := v.Type()

		for i := 0; i < v.NumField(); i++ {
			metricsName := typeOfS.Field(i).Name

			if _, ok := metrics[metricsName]; ok {

				switch fmt.Sprintf("%T", v.Field(i).Interface()) {
				case "uint64":
					sm.metrics.SetGauge(metricsName, utils.Gauge(float64(v.Field(i).Interface().(uint64))))

				case "uint32":
					sm.metrics.SetGauge(metricsName, utils.Gauge(float64(v.Field(i).Interface().(uint32))))

				case "float64":
					sm.metrics.SetGauge(metricsName, utils.Gauge(v.Field(i).Interface().(float64)))

				}

			}

		}
		sm.metrics.SetGauge("pollCount", utils.Gauge(float64(pollCount)))
		sm.metrics.SetGauge("RandomValue", utils.Gauge(float64(rand.Float64())))

	}
}

func (sm *ServiceMetrics) Post(metricsType string, metricName string, metricValue string, url string) {
	r := bytes.NewReader([]byte{})
	fmt.Println("URL ", url+metricsType+"/"+metricName+"/"+metricValue)
	resp, err := http.Post(url+metricsType+"/"+metricName+"/"+metricValue, "text/plain", r)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
}

func (sm *ServiceMetrics) PostMessage() {
	addr := "http://" + sm.cfg.Port + "/update/"
	for {
		for name, value := range sm.metrics.MetricStorageAgent.GetAllGauge() {
			fmt.Println("v", value)
			sm.Post("gauge", name, strconv.FormatFloat(float64(value), 'f', -1, 64), addr)
		}

		for name, value := range sm.metrics.MetricStorageAgent.GetAllCounter() {
			fmt.Println("v", value)
			sm.Post("counter", name, strconv.FormatUint(uint64(value), 10), addr)
		}

		pollCount = 0
		time.Sleep(time.Duration(sm.cfg.ReportInterval) * time.Second)
	}
}
