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
			metrics := a.metrics.GetAllMetricList()
			numjobs := len(metrics)
			jobs := make(chan models.Metrics, numjobs)

			// создаем буферизованный канал для отправки результатов
			results := make(chan models.Metrics, numjobs)
			// создаем и запускаем 3 воркера, это и есть пул,
			// передаем id, это для наглядности, канал задач и канал результатов
			for w := 1; w <= a.cfg.RateLimit; w++ {
				go worker(w, jobs, results)
			}

			// отправляем в канал задач какие-то данные
			// задач у нас 5, а воркера 3, значит одновременно решается только 3 задачи
			for j := 1; j <= numjobs; j++ {
				jobs <- metrics[j-1]
			}
			// закрываем канал на стороне отправителя
			close(jobs)

			// забираем из канала результатов результаты ;)
			for b := 1; b <= numjobs; b++ {
				go a.metrics.PostMessage(ctx, metrics[b-1])
				<-results
			}
		}

	}

}

func worker(id int, jobs <-chan models.Metrics, results chan<- models.Metrics) {
	for val := range jobs {
		// для наглядности будем выводить какой рабочий начал работу и кго задачу
		fmt.Println("рабочий", id, "запущен задача", val)
		// немного замедлим выполнение рабочего
		time.Sleep(3 * time.Second)
		// для наглядности выводим какой рабочий завершил какую задачу
		fmt.Println("рабочий", id, "закончил задача", val)
		// отправляем результат в канал результатов
		results <- val
	}
}
