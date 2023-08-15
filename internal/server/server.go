package server

import (
	"context"
	"fmt"
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
	logger logger.Logger
	cfg    *config.ConfigServer
	router *gin.Engine
}

func NewServer(log logger.Logger, cfg *config.ConfigServer) *server {
	return &server{
		logger: log,
		cfg:    cfg,
	}
}

func (s *server) Run() error {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	repos := repository.NewRepositoryStorage(s.logger)

	s.router = gin.New()

	h := handler.NewHandler(repos, s.logger, s.cfg)

	s.router.GET("/", h.GetAllMetrics())

	s.router.POST("/update/:metricType/:metricName/:metricValue", h.PostMetrics())

	value := s.router.Group("/value")
	{
		value.GET("/counter/:metricName", h.GetMethodCounter())
		value.GET("/gauge/:metricName", h.GetMethodGauge())
	}

	go func() {

		if err := s.router.Run(s.cfg.Port); err != nil {

			fmt.Println("Can't ListenAndServe on port", s.cfg.Port)
		}
	}()

	<-ctx.Done()
	return nil
}
