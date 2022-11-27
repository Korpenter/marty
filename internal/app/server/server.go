package server

import (
	"context"
	"fmt"
	"github.com/Mldlr/marty/internal/app/config"
	"github.com/Mldlr/marty/internal/app/router"
	"github.com/samber/do"
	"go.uber.org/zap"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Server struct {
	cfg              *config.Config
	log              *zap.Logger
	srv              http.Server
	shutdownFinished chan struct{}
}

func NewServer(i *do.Injector) *Server {
	cfg := do.MustInvoke[*config.Config](i)
	log := do.MustInvoke[*zap.Logger](i)

	return &Server{
		cfg: cfg,
		log: log,
		srv: http.Server{
			Handler: router.NewRouter(i),
			Addr:    fmt.Sprintf(cfg.ServiceAddress),
		},
		shutdownFinished: make(chan struct{}),
	}
}

func (s *Server) RunHTTP() {
	if s.shutdownFinished == nil {
		s.shutdownFinished = make(chan struct{})
	}

	err := s.srv.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		s.log.Fatal("unexpected error from ListenAndServe", zap.Error(err))
	}

	s.log.Info("waiting for shutdown finishing...")
	<-s.shutdownFinished
	s.log.Info("shutdown finished")
}

func (s *Server) WaitForExitingSignal(timeout time.Duration) {
	var waiter = make(chan os.Signal, 1)
	signal.Notify(waiter, syscall.SIGTERM, syscall.SIGINT)

	<-waiter

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	err := s.srv.Shutdown(ctx)
	if err != nil {
		s.log.Info("shutting down: " + err.Error())
	} else {
		s.log.Info("shutdown processed successfully")
		close(s.shutdownFinished)
	}
}
