package server

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/tanya-mtv/metricsservice/internal/repository"

	"github.com/tanya-mtv/metricsservice/internal/fileservice"

	"github.com/gin-gonic/gin"
	"github.com/tanya-mtv/metricsservice/internal/config"
	"github.com/tanya-mtv/metricsservice/internal/handler"
	"github.com/tanya-mtv/metricsservice/internal/logger"
)

type server struct {
	cfg    *config.ConfigServer
	router *gin.Engine
	log    logger.Logger
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

	s.router = gin.New()

	repo := repository.NewStorage()

	h := handler.NewHandler(repo, s.cfg, s.log)

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

	fs := fileservice.NewMetricFiles(repo, s.cfg.FileName, s.cfg.Interval, s.log)
	if s.cfg.FileName != "" {
		if s.cfg.Restore {
			go fs.LoadLDataFromFile()
		}
		if s.cfg.Interval != 0 {
			go fs.SaveDataToFile(ctx)
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
