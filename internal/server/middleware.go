package server

import (
	"context"
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

			log.Info(r.Context(), "http request completed", fields...)
		})
	}
}

func RequestIDMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Generate a new request ID or retrieve it from the request context
			requestID := r.Header.Get("X-Request-ID")
			ctx := r.Context()
			if requestID == "" {
				u, err := uuid.NewRandom()
				if err != nil {
					// Handle the error, e.g., log it or return a default request ID
					requestID = "default-request-id"
				} else {
					requestID = u.String()
				}
				r.Header.Set("X-Request-ID", requestID)
			}

			// Set the request ID in the response header
			w.Header().Set("X-Request-ID", requestID)
			ctx = context.WithValue(ctx, "req-id", requestID)

			// Call the next handler in the chain
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
