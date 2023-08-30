package server

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/tanya-mtv/metricsservice/internal/config"
	"github.com/tanya-mtv/metricsservice/internal/handler"
	"github.com/tanya-mtv/metricsservice/internal/logger"
	"github.com/tanya-mtv/metricsservice/internal/repository"
)

type server struct {
	cfg     *config.ConfigServer
	router  *gin.Engine
	log     logger.Logger
	storage *repository.Storage
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

	repos := repository.NewMetricStorage()
	s.router = gin.New()
	s.storage = repository.NewStorage(repos, s.cfg, s.log)

	h := handler.NewHandler(s.storage, s.cfg, s.log)

	s.router.Use(h.GzipMiddleware())
	s.router.Use(h.WithLogging())

	s.router.GET("/", h.GetAllMetrics())

	s.router.POST("/update", h.PostMetricsUpdateJSON())
	s.router.POST("/update/:metricType/:metricName/:metricValue", h.PostMetrics())

	value := s.router.Group("/value")
	{
		value.POST("/", h.PostMetricsValueJSON())
		value.GET("/counter/:metricName", h.GetMethodCounter())
		value.GET("/gauge/:metricName", h.GetMethodGauge())
	}

	if s.cfg.FileName != "" {
		if s.cfg.Restore {
			go s.storage.LoadLDataFromFile()
		}
		if s.cfg.Interval != 0 {
			go s.storage.SaveDataToFile(ctx)
		}
	}

	go func() {
		s.log.Info("Connect listening on port: %s", s.cfg.Port)
		if err := s.router.Run(s.cfg.Port); err != nil {

			s.log.Fatal("Can't ListenAndServe on port", s.cfg.Port)
		}
	}()

	<-ctx.Done()
	return nil
}
