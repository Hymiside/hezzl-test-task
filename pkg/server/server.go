package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type ConfigServer struct {
	Host string
	Port string
}

type Server struct {
	httpServer *http.Server
}

func (s *Server) RunServer(handler *chi.Mux, c ConfigServer) error {

	s.httpServer = &http.Server{
		Addr:    fmt.Sprintf("%s:%s", c.Host, c.Port),
		Handler: handler,
	}
	return s.httpServer.ListenAndServe()
}

func (s *Server) ShutdownServer(ctx context.Context) error {
	err := s.httpServer.Shutdown(ctx)
	if err != nil {
		return fmt.Errorf("microservice error stop: %w", err)
	}
	return nil
}
