package server

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/tanya-mtv/metricsservice/internal/fileservice"

	"github.com/tanya-mtv/metricsservice/internal/repository"

	"github.com/gin-gonic/gin"
	"github.com/tanya-mtv/metricsservice/internal/config"
	"github.com/tanya-mtv/metricsservice/internal/logger"
)

type server struct {
	cfg    *config.ConfigServer
	router *gin.Engine
	log    logger.Logger
	cron   fileservice.DataOper
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

	stor := repository.NewMetricStorage()

	s.router = NewRouter(stor, s.cfg, s.log)

	s.cron = openStorage(ctx, stor, s.cfg.FileName, s.cfg.Interval, s.cfg.Restore, s.log)

	go func() {
		s.log.Info("Connect listening on port: %s", s.cfg.Port)
		if err := s.router.Run(s.cfg.Port); err != nil {

			s.log.Fatal("Can't ListenAndServe on port", s.cfg.Port)
		}
	}()

	<-ctx.Done()
	return nil
}
