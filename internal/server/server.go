package server

import (
	"context"
	"net/http"
	"time"

	"github.com/ckshitij/notify-srv/internal/logger"
)

type Server struct {
	httpServer *http.Server
	log        logger.Logger
}

func New(addr string, log logger.Logger, handler http.Handler) *Server {
	return &Server{
		log: log,
		httpServer: &http.Server{
			Addr:              addr,
			Handler:           handler,
			ReadTimeout:       10 * time.Second,
			WriteTimeout:      10 * time.Second,
			IdleTimeout:       60 * time.Second,
			ReadHeaderTimeout: 5 * time.Second,
		},
	}
}

func (s *Server) Start(ctx context.Context) error {
	s.log.Info(ctx, "http server starting",
		logger.String("addr", s.httpServer.Addr),
	)

	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	s.log.Info(ctx, "http server shutting down")
	return s.httpServer.Shutdown(ctx)
}
