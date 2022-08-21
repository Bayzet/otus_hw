package grpc

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/Bayzet/otus_hw/hw12_13_14_15_calendar/internal/logger"

	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
)

func UnaryServerRequestLoggerInterceptor(logger logger.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		startRequest := time.Now()

		hndlr, err := handler(ctx, req)

		var host string
		if p, ok := peer.FromContext(ctx); ok {
			host, _, err = net.SplitHostPort(p.Addr.String())
			if err != nil {
				return hndlr, errors.Wrap(err, "Ошибка разделения хоста и порта")
			}
		}

		msg := fmt.Sprintf("%v [%v] %v %v %v",
			host,
			startRequest.Format(time.RFC822Z),
			info.FullMethod,
			time.Since(startRequest),
			req,
		)
		logger.Info(msg)

		return hndlr, errors.Wrap(err, "Ошибка UnaryServerRequestLoggerInterceptor")
	}
}
