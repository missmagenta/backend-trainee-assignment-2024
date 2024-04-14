package cache

import (
	"backend-trainee-assignment-2024/config"
	"github.com/dgraph-io/ristretto"
)

func New(cfg config.CACHE) (*ristretto.Cache, error) {
	cache, err := ristretto.NewCache(&ristretto.Config{
		MaxCost:     cfg.MaxCost,
		NumCounters: cfg.NumCounters,
		BufferItems: 64,
	})

	if err != nil {
		panic(err)
	}

	return cache, err
}
