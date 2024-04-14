package app

import (
	"context"
	"net/http"
	"time"
)

const (
	defaultReadTimeout     = 5 * time.Second
	defaultWriteTimeout    = 5 * time.Second
	defaultIdleTimout      = 5 * time.Second
	defaultAddr            = ":80"
	defaultShutdownTimeout = 3 * time.Second
	defaultMaxHeaderBytes  = 2 << 18
)

type Server struct {
	Server          *http.Server
	notify          chan error
	ShutdownTimeout time.Duration
}

func New(handler http.Handler, opts ...Option) *Server {
	httpServer := &http.Server{
		Handler:        handler,
		ReadTimeout:    defaultReadTimeout,
		WriteTimeout:   defaultWriteTimeout,
		IdleTimeout:    defaultIdleTimout,
		MaxHeaderBytes: defaultMaxHeaderBytes,
		Addr:           defaultAddr,
	}

	s := &Server{
		Server:          httpServer,
		notify:          make(chan error, 1),
		ShutdownTimeout: defaultShutdownTimeout,
	}

	for _, opt := range opts {
		opt(s)
	}

	return s
}

func (s *Server) Start() {
	go func() {
		s.notify <- s.Server.ListenAndServe()
		close(s.notify)
	}()
}

func (s *Server) Notify() <-chan error {
	return s.notify
}

func (s *Server) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), s.ShutdownTimeout)
	defer cancel()

	return s.Server.Shutdown(ctx)
}
