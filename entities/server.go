package entities

import (
	"context"
	"net/http"

	"github.com/Util787/user-manager-api/internal/config"
)

type Server struct {
	httpServer *http.Server
}

func (s *Server) CreateAndRun(config *config.ServerConfig, handler http.Handler) error {
	s.httpServer = &http.Server{
		Addr:              ":" + config.Port,
		Handler:           handler,
		MaxHeaderBytes:    1 << 20, // 1 MB
		ReadHeaderTimeout: config.ReadHeaderTimeout,
		WriteTimeout:      config.WriteTimeout,
		ReadTimeout:       config.ReadTimeout,
	}

	if err := s.httpServer.ListenAndServe(); err != nil {
		return err
	}
	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
