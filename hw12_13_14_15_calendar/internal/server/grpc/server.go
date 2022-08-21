package grpc

import (
	"context"
	"log"
	"net"

	"github.com/Bayzet/otus_hw/hw12_13_14_15_calendar/internal/logger"

	"google.golang.org/grpc"

	"github.com/Bayzet/otus_hw/hw12_13_14_15_calendar/internal/app"

	"github.com/Bayzet/otus_hw/hw12_13_14_15_calendar/internal/gen/pb/calendarpb"
)

type Server struct {
	calendarpb.UnimplementedCalendarServer
	app    app.Application
	lsn    net.Listener
	logger logger.Logger
	server *grpc.Server
}

func NewServer(app app.Application, host, port string, logger logger.Logger) *Server {
	lsn, err := net.Listen("tcp", host+":"+port)
	if err != nil {
		log.Fatal(err)
	}

	return &Server{
		app:    app,
		lsn:    lsn,
		logger: logger,
	}
}

func (s *Server) Run(ctx context.Context) error {
	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			UnaryServerRequestLoggerInterceptor(s.logger),
		),
	)
	calendarpb.RegisterCalendarServer(grpcServer, s)

	s.server = grpcServer

	log.Printf("starting gRPC server on %s", s.lsn.Addr().String())
	if err := grpcServer.Serve(s.lsn); err != nil {
		log.Fatal(err)
	}

	return nil
}

func (s *Server) Stop(ctx context.Context) {
	log.Printf("stopping gRPC server on %s", s.lsn.Addr().String())

	s.server.Stop()
}
