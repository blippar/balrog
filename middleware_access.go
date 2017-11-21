package browser

import (
	"fmt"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/go-chi/chi/middleware"
	log "github.com/sirupsen/logrus"
)

// AccessLogMiddleware defines a default request logger middleware based on the service logger
func (s *Server) AccessLogMiddleware(next http.Handler) http.Handler {

	fn := func(w http.ResponseWriter, r *http.Request) {

		start := time.Now()
		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
		entry := log.WithFields(map[string]interface{}{
			"method": r.Method,
			"path":   r.URL.String(),
			"remote": r.RemoteAddr,
		})

		defer func() {

			end := time.Now()

			// Log stacktrace if request panicked
			if rec := recover(); rec != nil {
				cause := fmt.Sprintf("panic: %+v\n", rec)
				fmt.Printf("\n%s\n%s\n", cause, debug.Stack())
				s.errorHandler(ww, r, http.StatusInternalServerError, cause)
			}

			entry.WithFields(map[string]interface{}{
				"status":  ww.Status(),
				"latency": end.Sub(start),
			}).Info("requestServed")

		}()

		next.ServeHTTP(ww, r)
	}
	return http.HandlerFunc(fn)
}
