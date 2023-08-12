package metrics

import (
	"bytes"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"reflect"
	"runtime"
	"strconv"
	"time"

	"github.com/tanya-mtv/metricsservice/internal/config"
	"github.com/tanya-mtv/metricsservice/internal/repository"

	"github.com/tanya-mtv/metricsservice/internal/utils"
)

var pollCount int64

var reqmetrics = map[string]bool{
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
	cfg               *config.ConfigAgent
	metricsRepository *repository.Repository
	reqmetrics        map[string]bool
}

func NewServiceMetrics(cfg *config.ConfigAgent, metricsRepository *repository.Repository) *ServiceMetrics {
	reqmetrics := make(map[string]bool, 29)
	reqmetrics["Alloc"] = true
	reqmetrics["BuckHashSys"] = true
	reqmetrics["Frees"] = true
	reqmetrics["GCCPUFraction"] = true
	reqmetrics["GCSys"] = true
	reqmetrics["HeapAlloc"] = true
	reqmetrics["HeapIdle"] = true
	reqmetrics["HeapInuse"] = true
	reqmetrics["HeapObjects"] = true
	reqmetrics["HeapReleased"] = true
	reqmetrics["HeapSys"] = true
	reqmetrics["LastGC"] = true
	reqmetrics["Lookups"] = true
	reqmetrics["MCacheInuse"] = true
	reqmetrics["MCacheSys"] = true
	reqmetrics["MSpanInuse"] = true
	reqmetrics["MSpanSys"] = true
	reqmetrics["Mallocs"] = true
	reqmetrics["NextGC"] = true
	reqmetrics["NumForcedGC"] = true
	reqmetrics["NumGC"] = true
	reqmetrics["OtherSys"] = true
	reqmetrics["PauseTotalNs"] = true
	reqmetrics["StackInuse"] = true
	reqmetrics["StackSys"] = true
	reqmetrics["Sys"] = true
	reqmetrics["TotalAlloc"] = true
	reqmetrics["PollCount"] = true
	reqmetrics["RandomValue"] = true

	return &ServiceMetrics{
		cfg:               cfg,
		metricsRepository: metricsRepository,
		reqmetrics:        reqmetrics,
	}
}

func (sm *ServiceMetrics) NewMonitor() {

	var rtm runtime.MemStats
	interval := time.Duration(sm.cfg.PollInterval) * time.Second
	for {
		// <-time.After(interval)
		time.Sleep(interval)
		pollCount += 1

		runtime.ReadMemStats(&rtm)

		v := reflect.ValueOf(rtm)
		typeOfS := v.Type()

		for i := 0; i < v.NumField(); i++ {
			metricsName := typeOfS.Field(i).Name

			if _, ok := reqmetrics[metricsName]; ok {

				switch fmt.Sprintf("%T", v.Field(i).Interface()) {
				case "uint64":
					sm.metricsRepository.SetValueGauge(metricsName, utils.Gauge(float64(v.Field(i).Interface().(uint64))))

				case "uint32":
					sm.metricsRepository.SetValueGauge(metricsName, utils.Gauge(float64(v.Field(i).Interface().(uint32))))

				case "float64":
					sm.metricsRepository.SetValueGauge(metricsName, utils.Gauge(v.Field(i).Interface().(float64)))

				}

			}

		}
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
			fmt.Println("v", value)
			body, err := sm.Post("gauge", name, strconv.FormatFloat(float64(value), 'f', -1, 64), addr)

			if err != nil {
				panic("error reading body")
			}

			fmt.Println("data was sent successfuly", body)
		}

		for name, value := range sm.metricsRepository.GetAllCounter() {
			fmt.Println("v", value)
			body, err := sm.Post("counter", name, strconv.FormatUint(uint64(value), 10), addr)

			if err != nil {
				panic("error reading body")
			}

			fmt.Println("data was sent successfuly", body)
		}

		pollCount = 0
		time.Sleep(time.Duration(sm.cfg.ReportInterval) * time.Second)
	}
}
