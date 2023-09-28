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

	// создаем буферизованный канал для отправки результатов
	results := make(chan models.Metrics, numjobs)

	// создаем и запускаем 3 воркера, это и есть пул,
	// передаем id, это для наглядности, канал задач и канал результатов
	for w := 1; w <= a.cfg.RateLimit; w++ {
		go worker(jobs, results)
	}

	// отправляем в канал задач метрики
	// задач у нас 5, а воркера 3, значит одновременно решается только 3 задачи
	for j := 1; j <= numjobs; j++ {
		fmt.Println("get metric  ", metrics[j-1])
		jobs <- metrics[j-1]
	}
	// закрываем канал на стороне отправителя
	close(jobs)

	go a.recieveChainData(ctx, results)
}
func (a *agent) recieveChainData(ctx context.Context, res chan models.Metrics) {
	defer close(res)
	for {
		val, ok := <-res
		if !ok {
			fmt.Println(val, "<-- loop broke!")
			break
		} else {
			a.metrics.PostMessage(ctx, val)

		}
	}

}

func worker(jobs <-chan models.Metrics, results chan<- models.Metrics) {
	for val := range jobs {
		// немного замедлим выполнение рабочего
		time.Sleep(1 * time.Second)
		results <- val
	}
}
