package internalhttp

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/Bayzet/otus_hw/hw12_13_14_15_calendar/internal/consts"

	"github.com/Bayzet/otus_hw/hw12_13_14_15_calendar/internal/logger"
)

func LoggingMiddleware(logger *logger.Logger) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		rTime := ctx.Value(consts.StartRequestTime).(time.Time)
		ip := strings.Split(r.RemoteAddr, ":")

		msg := fmt.Sprintf("%v [%v] %v %v %v %v %v %v",
			ip[0],
			rTime.Format(time.RFC822Z),
			r.Method,
			r.RequestURI,
			r.Proto,
			r.Context().Value(consts.ResponseStatusCode),
			time.Since(rTime),
			r.Header.Get("User-Agent"),
		)
		logger.Info(msg)
	})
}
