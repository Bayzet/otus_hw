package internalhttp

import (
	"context"
	"log"
	"net/http"
)

type Server struct {
	addr string
	mux  *http.ServeMux
}

func NewServer(addr string, mux *http.ServeMux) *Server {
	return &Server{
		addr: addr,
		mux:  mux,
	}
}

func (s *Server) Start(ctx context.Context) error {
	if err := http.ListenAndServe(s.addr, s.mux); err != nil {
		log.Fatal(err)
	}

	<-ctx.Done()
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	ctx.Done()
	return nil
}
