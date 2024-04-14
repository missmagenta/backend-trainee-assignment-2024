package handlers

import (
	"backend-trainee-assignment-2024/internal/httpserver/handler/banner"
	"backend-trainee-assignment-2024/internal/usecase"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

type Router struct {
	mux      *chi.Mux
	useCases usecase.UseCases
}

func NewRouter(mux *chi.Mux, useCases usecase.UseCases) *Router {
	return &Router{
		mux:      mux,
		useCases: useCases,
	}
}

func (r *Router) SetupMiddlewares() {
	r.mux.Use(cors.AllowAll().Handler)
	r.mux.Use(middleware.Recoverer)
	r.mux.Use(middleware.RequestID)
}

func (r *Router) SetupRoutes() {
	r.mux.Group(func(router chi.Router) {
		banner.New(router, r.useCases.Banner)
	})
}
