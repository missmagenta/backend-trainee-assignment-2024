package memory

import (
	"github.com/dgraph-io/ristretto"
	"time"
)

type Repositories struct {
	Banner Banner
}

func NewRepositories(cache *ristretto.Cache, ttl time.Duration) Repositories {
	return Repositories{
		Banner: NewBanner(cache, ttl),
	}
}
