package internalhttp

import (
	"fmt"
	"net/http"
	"time"
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
	return n, err
}

func loggingMiddleware(logger Logger, next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		writer := respStatusWriter{ResponseWriter: w}
		startRequest := time.Now()
		next.ServeHTTP(&writer, r)
		msg := fmt.Sprintf("%v [%v] %v %v %v %v %v %v %v",
			r.RemoteAddr,
			startRequest.Format(time.RFC822Z),
			r.Method,
			r.RequestURI,
			r.Proto,
			writer.status,
			writer.length,
			time.Since(startRequest),
			r.Header.Get("User-Agent"),
		)
		logger.Info(msg)
	})
}
