package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/tanya-mtv/metricsservice/internal/repository"
	"github.com/tanya-mtv/metricsservice/internal/servise"

	"github.com/tanya-mtv/metricsservice/internal/config"
	"github.com/tanya-mtv/metricsservice/internal/handler"
	"github.com/tanya-mtv/metricsservice/internal/logger"
)

type server struct {
	logger     logger.Logger
	httpServer *http.Server
	cfg        *config.ConfigServer
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

	// db, err := repository.NewPostgresDB(*s.cfg.Postgresql)

	// if err != nil {
	// 	s.logger.Error("Failed to initialaze db:", zap.String("string connection", "postgres"), err)
	// } else {
	// 	s.logger.Info("successful connection to database", zap.String("Connection", "Ok"))
	// }

	// redis_db, err := repository.NewRedisDB(s.cfg.Redis, ctx)

	// if err != nil {
	// 	s.logger.Error("Failed to initialaze redis db:", zap.String("string connection", "redis"), err)
	// } else {
	// 	s.logger.Info("successful connection to Redis", zap.String("Connection", "Ok"))
	// }

	// repos := repository.NewRepository(db, redis_db, s.logger)
	// serv := servise.NewServise(repos)

	// s.ps = serv
	repos := repository.NewRepository(s.logger)
	serv := servise.NewServise(repos)

	handl := handler.NewHandler(serv, s.logger, s.cfg)

	// amen := amenitie.NewAmenitie(s.ps, s.logger, s.cfg)
	httpServer := &http.Server{
		Addr:           ":" + s.cfg.Port,
		Handler:        handl.InitRoutes(),
		MaxHeaderBytes: 1 << 20,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
	}

	s.httpServer = httpServer

	// s.logger.Info("Starting Reader Kafka consumers")

	go func() {
		// s.logger.Info("Writer microservice connectlistening on port: %s", s.httpServer.Addr)
		if err := s.httpServer.ListenAndServe(); err != nil {
			// s.logger.WarnMsg("ListenAndServe", err)
			fmt.Println("ListenAndServe")
		}
	}()

	<-ctx.Done()
	return nil
}
