package agent

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/tanya-mtv/metricsservice/internal/config"
	"github.com/tanya-mtv/metricsservice/internal/logger"
	"github.com/tanya-mtv/metricsservice/internal/metrics"
	"github.com/tanya-mtv/metricsservice/internal/repository"
	"github.com/tanya-mtv/metricsservice/internal/servise"
)

type agent struct {
	logger  logger.Logger
	cfg     *config.ConfigAgent
	metrics *metrics.ServiceMetrics
	service *servise.Service
}

func NewAgent(log logger.Logger, cfg *config.ConfigAgent) *agent {
	return &agent{
		logger: log,
		cfg:    cfg,
	}
}

func (a *agent) Run() error {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	repos := repository.NewRepository(a.logger)

	a.service = servise.NewServise(repos)
	a.metrics = metrics.NewServiceMetrics(a.cfg, a.service)

	time.Sleep(time.Duration(a.cfg.ReportInterval) * time.Second)

	go a.metrics.NewMonitor()

	go a.metrics.PostMessage()
	// addr := "http://" + a.cfg.Port + "/update/"
	// for {
	// 	for name, value := range a.service.valuesGauge {
	// 		fmt.Println("k", name)
	// 		fmt.Println("v", value)
	// 		if name == "pollCount" {
	// 			a.metrics.Post("counter", name, strconv.FormatUint(uint64(value), 10), addr)
	// 		} else {
	// 			a.metrics.Post("gauge", name, strconv.FormatFloat(value, 'f', -1, 64), addr)
	// 		}

	// 	}

	// 	pollCount = 0
	// 	time.Sleep(time.Duration(a.cfg.ReportInterval) * time.Second)
	// }

	<-ctx.Done()
	return nil
}
