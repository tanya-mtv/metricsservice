package agent

import (
	"context"
	"os"
	"os/signal"
	"syscall"

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

	go a.metrics.NewMonitor()

	go a.metrics.PostMessage()

	<-ctx.Done()
	return nil
}
