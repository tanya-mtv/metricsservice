package agent

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/tanya-mtv/metricsservice/internal/config"
	"github.com/tanya-mtv/metricsservice/internal/metrics"
	"github.com/tanya-mtv/metricsservice/internal/repository"
)

type agent struct {
	cfg     *config.ConfigAgent
	metrics *metrics.ServiceMetrics
}

func NewAgent(cfg *config.ConfigAgent) *agent {
	return &agent{
		cfg: cfg,
	}
}

func (a *agent) Run() error {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	repos := repository.NewMetricRepositoryCollector()

	a.metrics = metrics.NewServiceMetrics(a.cfg, repos)

	go a.metrics.MetricsMonitor()

	go a.metrics.PostMessage()

	<-ctx.Done()
	return nil
}
