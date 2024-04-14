package app

import (
	"backend-trainee-assignment-2024/config"
	"backend-trainee-assignment-2024/internal/httpserver/handler"
	"backend-trainee-assignment-2024/internal/usecase"
	"github.com/go-chi/chi/v5"
	"net"
	"time"
)

func newHttpServer(http config.HTTP, dependencies usecase.Dependencies) *Server {
	r := chi.NewRouter()

	useCases := usecase.NewUseCases(dependencies)
	router := handlers.NewRouter(r, useCases)
	router.SetupMiddlewares()
	router.SetupRoutes()

	return New(r,
		Port(http.Port),
		MaxHeaderBytes(http.MaxHeaderBytes),
		IdleTimeout(http.IdleTimeout),
		WriteTimeout(http.WriteTimeout),
		ReadTimeout(http.ReadTimeout),
	)
}

type Option func(server *Server)

func Port(port string) Option {
	return func(s *Server) {
		s.Server.Addr = net.JoinHostPort("", port)
	}
}

func ReadTimeout(timeout time.Duration) Option {
	return func(s *Server) {
		s.Server.ReadTimeout = timeout
	}
}

func WriteTimeout(timeout time.Duration) Option {
	return func(s *Server) {
		s.Server.WriteTimeout = timeout
	}
}

func IdleTimeout(timeout time.Duration) Option {
	return func(s *Server) {
		s.Server.IdleTimeout = timeout
	}
}

func MaxHeaderBytes(max int) Option {
	return func(s *Server) {
		s.Server.MaxHeaderBytes = max
	}
}

func ShutdownTimeout(timeout time.Duration) Option {
	return func(s *Server) {
		s.ShutdownTimeout = timeout
	}
}
