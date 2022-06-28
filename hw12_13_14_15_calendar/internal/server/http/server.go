package internalhttp

import (
	"context"
	"io"
	"net/http"
	"os"
)

type Server struct {
	Host string
	Port string
}

type Application interface { // TODO
}

type Logger interface {
	Info(string)
	Warn(string)
	Error(string)
	Debug(string)
}

func helloHandler(res http.ResponseWriter, req *http.Request) {
	io.WriteString(res, "Hello world!")
}

func NewServer(host, port string, logger Logger, app Application) *Server {
	http.Handle("/", loggingMiddleware(logger, http.HandlerFunc(helloHandler)))

	return &Server{
		Host: host,
		Port: port,
	}
}

func (s *Server) Start(ctx context.Context) error {
	http.ListenAndServe(s.Host+":"+s.Port, nil)

	<-ctx.Done()
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	os.Exit(1)
	return nil
}
