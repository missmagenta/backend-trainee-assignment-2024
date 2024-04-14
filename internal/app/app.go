package app

import (
	"backend-trainee-assignment-2024/internal/usecase"
	"backend-trainee-assignment-2024/internal/usecase/repo/memory"
	pg "backend-trainee-assignment-2024/internal/usecase/repo/postgres"
	"backend-trainee-assignment-2024/pkg/cache"
	"backend-trainee-assignment-2024/pkg/postgres"
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"backend-trainee-assignment-2024/config"
)

type App struct {
	httpServer *Server
	db         *postgres.Postgres
}

func NewApp(cfg *config.Config) (App, error) {
	db, err := openDB(cfg.PG)
	if err != nil {
		return App{}, err
	}

	cache, err := cache.New(cfg.CACHE)
	if err != nil {
		return App{}, err
	}

	pgRepositories := pg.NewRepositories(db)
	memoryRepositories := memory.NewRepositories(cache, cfg.CACHE.TTL)

	httpServer := newHttpServer(cfg.HTTP, usecase.Dependencies{
		Pg:     pgRepositories,
		Memory: memoryRepositories,
	})

	return App{httpServer: httpServer, db: db}, nil
}

func (app App) Run() error {
	app.httpServer.Start()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	var err error
	select {
	case s := <-interrupt:
		err = errors.New("app - Run - signal: " + s.String())
	case err = <-app.httpServer.Notify():
		err = fmt.Errorf("app - Run - httpServer.Notify: %w", err)
	}

	return err
}

func (app App) Shutdown() error {
	httpErr := app.httpServer.Shutdown()
	if httpErr != nil {
		log.Println(fmt.Errorf("app - Run - httpServer.Shutdown: %w", httpErr))
	}
	dbErr := app.db.Close()
	if httpErr != nil {
		log.Println(fmt.Errorf("app - Run - db.Close: %w", dbErr))
	}
	return errors.Join(httpErr, dbErr)
}
