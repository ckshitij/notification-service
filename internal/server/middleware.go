package server

import (
	"net/http"
	"strings"
	"time"

	"github.com/ckshitij/notify-srv/internal/logger"
	"github.com/google/uuid"
)

func AccessLogMiddleware(log logger.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			rw := newResponseWriter(w)
			next.ServeHTTP(rw, r)

			duration := time.Since(start)

			fields := []logger.Field{
				logger.String("method", r.Method),
				logger.String("path", r.URL.Path),
				logger.Int("status", rw.status),
				logger.Int("bytes", rw.bytes),
				logger.Int("duration_micro_sec", int(duration.Microseconds())),
				logger.String("remote_ip", r.RemoteAddr),
			}

			// Internal endpoints: default to Debug, but still respect errors
			if strings.Contains(r.URL.Path, "/internal/") {
				if rw.status >= http.StatusInternalServerError {
					log.Error(r.Context(), "internal request failed", fields...)
				} else {
					log.Debug(r.Context(), "internal request completed", fields...)
				}
				return
			}

			// External APIs
			switch {
			case rw.status >= http.StatusInternalServerError:
				log.Error(r.Context(), "http request failed", fields...)

			case rw.status >= http.StatusBadRequest:
				log.Warn(r.Context(), "http request client error", fields...)

			default:
				log.Info(r.Context(), "http request completed", fields...)
			}
		})
	}
}

func RequestIDMiddleware(log logger.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			reqID := r.Header.Get("X-Request-ID")
			if reqID == "" {
				reqID = uuid.NewString()
			}

			ctx := logger.WithRequestID(r.Context(), reqID)
			w.Header().Set("X-Request-ID", reqID)

			log.Info(ctx, "incoming request",
				logger.String("method", r.Method),
				logger.String("path", r.URL.Path),
			)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
