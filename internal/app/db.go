package app

import (
	"backend-trainee-assignment-2024/config"
	"backend-trainee-assignment-2024/pkg/postgres"
	"github.com/uptrace/bun"
	"time"
)

func MaxOpenConnsDB(max int) bun.DBOption {
	return func(db *bun.DB) {
		db.SetMaxOpenConns(max)
	}
}

func MaxIdleConnsDB(max int) bun.DBOption {
	return func(db *bun.DB) {
		db.SetMaxIdleConns(max)
	}
}

func ConnMaxLifetimeDB(max time.Duration) bun.DBOption {
	return func(db *bun.DB) {
		db.SetConnMaxLifetime(max)
	}
}

func ConnMaxIdleTimeDB(max time.Duration) bun.DBOption {
	return func(db *bun.DB) {
		db.SetConnMaxIdleTime(max)
	}
}

func openDB(cfg config.PG) (*postgres.Postgres, error) {
	db, err := postgres.New(cfg.URL,
		MaxOpenConnsDB(cfg.MaxOpenConns),
		ConnMaxIdleTimeDB(cfg.ConnMaxIdleTime),
		MaxIdleConnsDB(cfg.MaxIdleConns),
		ConnMaxLifetimeDB(cfg.ConnMaxLifetime),
	)
	if err != nil {
		return nil, err
	}
	return db, nil
}
