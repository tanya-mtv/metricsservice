package server

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/lib/pq"

	"github.com/tanya-mtv/metricsservice/internal/models"
	"github.com/tanya-mtv/metricsservice/internal/repository"

	"github.com/gin-gonic/gin"
	"github.com/tanya-mtv/metricsservice/internal/config"
	"github.com/tanya-mtv/metricsservice/internal/logger"
)

type metricStorage interface {
	UpdateCounter(n string, v int64) repository.Counter
	UpdateGauge(n string, v float64) repository.Gauge
	GetAll() []models.Metrics
	GetCounter(metricName string) (repository.Counter, bool)
	GetGauge(metricName string) (repository.Gauge, bool)
	UpdateMetrics([]*models.Metrics) ([]*models.Metrics, error)
}

type server struct {
	cfg    *config.ConfigServer
	router *gin.Engine
	log    logger.Logger
	stor   metricStorage
}

func NewServer(cfg *config.ConfigServer, log logger.Logger) *server {
	return &server{
		cfg: cfg,
		log: log,
	}
}

func (s *server) Run() error {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	db, err := repository.NewPostgresDB(s.cfg.DSN)

	if err != nil {
		s.log.Info("Failed to initialaze db: %s", err.Error())
	} else {
		s.log.Info("Success connection to db")
		defer db.Close()
	}
	s.openStorage(ctx, db)

	s.router = s.NewRouter(db)

	go func() {
		s.log.Info("Connect listening on port: %s", s.cfg.Port)
		if err := s.router.Run(s.cfg.Port); err != nil {

			s.log.Fatal("Can't ListenAndServe on port", s.cfg.Port)
		}
	}()

	<-ctx.Done()
	return nil
}
