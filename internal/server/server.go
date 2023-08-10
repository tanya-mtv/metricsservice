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

	repos := repository.NewRepository(s.logger)
	serv := servise.NewServise(repos)

	handl := handler.NewHandler(serv, s.logger, s.cfg)

	httpServer := &http.Server{
		Addr:           s.cfg.Port,
		Handler:        handl.InitRoutes(),
		MaxHeaderBytes: 1 << 20,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
	}

	s.httpServer = httpServer

	go func() {

		if err := s.httpServer.ListenAndServe(); err != nil {

			fmt.Println("ListenAndServe")
		}
	}()

	<-ctx.Done()
	return nil
}
