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
			a.metrics.MetricsMonitor()
		case <-reportIntervalTicker.C:
			// a.metrics.PostMessageJSON(ctx)
			a.metrics.PostMessage(ctx)
		}
	}

}
