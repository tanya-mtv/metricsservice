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

	repos := repository.NewMetricRepositoryStorage()
	s.router = gin.New()

	h := handler.NewHandler(repos, s.cfg)

	s.router.Use(h.WithLogging(s.log))

	s.router.GET("/", h.GetAllMetrics(s.log))

	s.router.POST("/update/", h.PostMetricsJSON(s.log))
	// s.router.POST("/update/:metricType/:metricName/:metricValue", h.PostMetrics())

	value := s.router.Group("/value")
	{
		value.GET("/counter/:metricName", h.GetMethodCounter())
		value.GET("/gauge/:metricName", h.GetMethodGauge())
	}

	go func() {
		s.log.Info("Connectlistening on port: %s", s.cfg.Port)
		if err := s.router.Run(s.cfg.Port); err != nil {

			s.log.Fatal("Can't ListenAndServe on port", s.cfg.Port)
		}
	}()

	<-ctx.Done()
	return nil
}
