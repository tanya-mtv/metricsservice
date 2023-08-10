package main

import (
	"github.com/tanya-mtv/metricsservice/internal/agent"
	"github.com/tanya-mtv/metricsservice/internal/config"
	"github.com/tanya-mtv/metricsservice/internal/logger"
)

func main() {
	cfg, err := config.InitConfigAgent()
	if err != nil {
		// sugarLogger.Error("error initialazing config", zap.String("initConfig", "fail"), err)
		panic("error initialazing config")
	}

	appLogger := logger.NewAppLogger()
	ag := agent.NewAgent(appLogger, cfg)
	// appLogger.Fatal(srv.Run())
	if err := ag.Run(); err != nil {
		panic(err)
	}

	// go NewMonitor(cfg.PollInterval)

}

// func NewMonitor(duration int) {

// 	var rtm runtime.MemStats
// 	var interval = time.Duration(duration) * time.Second
// 	for {
// 		<-time.After(interval)
// 		pollCount += 1
// 		// Read full mem stats
// 		runtime.ReadMemStats(&rtm)

// 		v := reflect.ValueOf(rtm)
// 		typeOfS := v.Type()

// 		for i := 0; i < v.NumField(); i++ {
// 			metricsName := typeOfS.Field(i).Name

// 			if _, ok := metrics[metricsName]; ok {
// 				switch fmt.Sprintf("%T", v.Field(i).Interface()) {
// 				case "uint64":
// 					valuesGauge[metricsName] = float64(v.Field(i).Interface().(uint64))
// 				case "uint32":
// 					valuesGauge[metricsName] = float64(v.Field(i).Interface().(uint32))
// 				case "float64":
// 					valuesGauge[metricsName] = v.Field(i).Interface().(float64)

// 				}

// 			}

// 		}
// 		valuesGauge["pollCount"] = float64(pollCount)
// 		valuesGauge["RandomValue"] = rand.Float64()

// 	}
// }

// func post(metricsType string, metricName string, metricValue string, url string) {
// 	r := bytes.NewReader([]byte{})
// 	fmt.Println("URL ", url+metricsType+"/"+metricName+"/"+metricValue)
// 	resp, err := http.Post(url+metricsType+"/"+metricName+"/"+metricValue, "text/plain", r)
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer resp.Body.Close()
// }
