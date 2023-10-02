package agent

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/tanya-mtv/metricsservice/internal/models"

	"github.com/tanya-mtv/metricsservice/internal/config"
	"github.com/tanya-mtv/metricsservice/internal/logger"
	"github.com/tanya-mtv/metricsservice/internal/metrics"
	"github.com/tanya-mtv/metricsservice/internal/repository"
)

type agent struct {
	cfg     *config.ConfigAgent
	metrics *metrics.ServiceMetrics
	log     logger.Logger
}

func NewAgent(cfg *config.ConfigAgent, log logger.Logger) *agent {
	return &agent{
		cfg: cfg,
		log: log,
	}
}

func (a *agent) Run() error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	collector := repository.NewMetricRepositoryCollector()
	a.metrics = metrics.NewServiceMetrics(collector, a.cfg, a.log)

	pollIntervalTicker := time.NewTicker(time.Duration(a.cfg.PollInterval) * time.Second)
	defer pollIntervalTicker.Stop()

	reportIntervalTicker := time.NewTicker(time.Duration(a.cfg.ReportInterval) * time.Second)
	defer reportIntervalTicker.Stop()

	for {
		select {
		case <-ctx.Done():
			stop()

			return nil
		case <-pollIntervalTicker.C:
			go a.metrics.MetricsMonitor()
			go a.metrics.MetricsMonitorGopsutil(ctx)
		case <-reportIntervalTicker.C:
			go a.createWorkerPool(ctx)
		}
	}

}

func (a *agent) createWorkerPool(ctx context.Context) {
	metrics := a.metrics.GetAllMetricList()
	numjobs := len(metrics)
	jobs := make(chan models.Metrics, numjobs)

	for w := 1; w <= a.cfg.RateLimit; w++ {
		go func(w int) {
			a.worker(ctx, jobs, w)
		}(w)
	}

	for j := 1; j <= numjobs; j++ {
		fmt.Println("get metric  ", metrics[j-1])
		jobs <- metrics[j-1]
	}
	close(jobs)
}
func (a *agent) recieveChainData(ctx context.Context, jobs <-chan models.Metrics, w int) {
	for {
		select {
		case <-ctx.Done():
			return
		case val, ok := <-jobs:
			if !ok {
				fmt.Println("<-- loop broke!")
				return
			} else {
				fmt.Println("worker ", w, "send metric", val)
				go a.metrics.PostMessage(ctx, val)
			}
		}
	}

}

func (a *agent) worker(ctx context.Context, jobs <-chan models.Metrics, w int) {
	a.recieveChainData(ctx, jobs, w)
}
