package internalhttp

import (
	"context"
	"log"
	"net/http"

	"github.com/pkg/errors"

	"github.com/Bayzet/otus_hw/hw12_13_14_15_calendar/internal/logger"
)

type Server struct {
	addr   string
	server *http.Server
}

func NewServer(router http.Handler, host, port string, logger logger.Logger) *Server {
	addr := host + ":" + port

	return &Server{
		addr: addr,
		server: &http.Server{
			Addr:    addr,
			Handler: loggingMiddleware(logger, router),
		},
	}
}

func (s *Server) Start(ctx context.Context) error {
	log.Printf("starting HTTP server on %s", s.addr)

	err := s.server.ListenAndServe()
	if err != nil {
		return errors.Wrap(err, "Ошибка старта HTTP сервера")
	}

	<-ctx.Done()
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	log.Printf("stopping HTTP server on %s", s.addr)

	if err := s.server.Close(); err != nil {
		return errors.Wrap(err, "Ошибка остановки HTTP сервера")
	}

	return nil
}
