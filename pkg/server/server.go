package server

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type ConfigServer struct {
	Host string
	Port string
}

type Server struct{}

func (s *Server) RunServer(ctx context.Context, handler *chi.Mux, c ConfigServer) error {
	httpServer := &http.Server{
		Addr:    fmt.Sprintf("%s:%s", c.Host, c.Port),
		Handler: handler,
	}

	go func(ctx context.Context) {
		<-ctx.Done()
		httpServer.Shutdown(ctx)
	}(ctx)

	log.Printf("authentication microservice launched on http://%s:%s/", c.Host, c.Port)
	return httpServer.ListenAndServe()
}
