package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"net/http"
	"reflect"
	"runtime"
	"strconv"

	"time"

	"github.com/tanya-mtv/metricsservice/internal/utils"
)

type Monitor1 struct {
	Alloc, HeapAlloc, HeapIdle, HeapInuse, HeapObjects, HeapReleased, HeapSys,
	TotalAlloc, LastGC, Lookups, MCacheInuse, MCacheSys, MSpanInuse, MSpanSys, NextGC, StackInuse, StackSys,
	Sys,
	Mallocs,
	Frees,
	PauseTotalNs,
	BuckHashSys,
	OtherSys,
	GCSys uint64
	NumGC,
	NumForcedGC uint32
	GCCPUFraction float64
	PollCount     utils.Counter
	RandomValue   utils.Gauge
}

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

const pollInterval = 2
const reportInterval = 10
const URL = "http://127.0.0.1:8080/update/"

var valuesGauge = map[string]float64{}
var pollCount uint64

func main() {
	go NewMonitor(pollInterval)

	time.Sleep(reportInterval * time.Second)

	for {
		for name, value := range valuesGauge {
			fmt.Println("k", name)
			fmt.Println("v", value)
			if name == "pollCount" {
				post("counter", name, strconv.FormatUint(uint64(value), 10))
			} else {
				post("gauge", name, strconv.FormatFloat(value, 'f', -1, 64))
			}

		}

		pollCount = 0
		time.Sleep(reportInterval * time.Second)
	}
}

func NewMonitor(duration int) {

	var rtm runtime.MemStats
	var interval = time.Duration(duration) * time.Second
	for {
		<-time.After(interval)
		pollCount += 1
		// Read full mem stats
		runtime.ReadMemStats(&rtm)

		v := reflect.ValueOf(rtm)
		typeOfS := v.Type()

		for i := 0; i < v.NumField(); i++ {
			metricsName := typeOfS.Field(i).Name

			if _, ok := metrics[metricsName]; ok {
				// fmt.Printf("Field: %s\tValue: %v\n", typeOfS.Field(i).Name, v.Field(i).Interface())

				switch fmt.Sprintf("%T", v.Field(i).Interface()) {
				case "uint64":
					valuesGauge[metricsName] = float64(v.Field(i).Interface().(uint64))
				case "uint32":
					valuesGauge[metricsName] = float64(v.Field(i).Interface().(uint32))
				case "float64":
					valuesGauge[metricsName] = v.Field(i).Interface().(float64)
					// default:

				}

			}

		}
		valuesGauge["pollCount"] = float64(pollCount)
		valuesGauge["RandomValue"] = rand.Float64()

	}
}

func post(metricsType string, metricName string, metricValue string) {
	r := bytes.NewReader([]byte{})
	resp, err := http.Post(URL+metricsType+"/"+metricName+"/"+metricValue, "text/plain", r)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
}
