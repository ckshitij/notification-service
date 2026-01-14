package server

import (
	"net/http"
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

			log.Info(r.Context(), "http request completed",
				logger.String("method", r.Method),
				logger.String("path", r.URL.Path),
				logger.Int("status", rw.status),
				logger.Int("bytes", rw.bytes),
				logger.Int("duration_micro_sec", int(duration.Microseconds())),
				logger.String("remote_ip", r.RemoteAddr),
			)
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
