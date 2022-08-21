package internalhttp

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Bayzet/otus_hw/hw12_13_14_15_calendar/internal/logger"

	"github.com/pkg/errors"
)

type respStatusWriter struct {
	http.ResponseWriter
	length int
	status int
}

func (w *respStatusWriter) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}

func (w *respStatusWriter) Write(b []byte) (int, error) {
	if w.status == 0 {
		w.status = 200
	}
	n, err := w.ResponseWriter.Write(b)
	w.length += n
	return n, errors.Wrap(err, "Ошибка записи response")
}

func loggingMiddleware(logger logger.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		writer := respStatusWriter{ResponseWriter: w}
		startRequest := time.Now()
		next.ServeHTTP(&writer, req)
		msg := fmt.Sprintf("%v [%v] %v %v %v %v %v %v %v",
			req.RemoteAddr,
			startRequest.Format(time.RFC822Z),
			req.Method,
			req.RequestURI,
			req.Proto,
			writer.status,
			writer.length,
			time.Since(startRequest),
			req.Header.Get("User-Agent"),
		)
		logger.Info(msg)
	})
}
