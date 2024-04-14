package memory

import (
	"backend-trainee-assignment-2024/internal/entity"
	"errors"
	"fmt"
	"github.com/dgraph-io/ristretto"
	"time"
)

type Banner struct {
	cache *ristretto.Cache
	ttl   time.Duration
}

type Key struct {
	FeatureId int
	TagId     int
}

func (k Key) String() string {
	return fmt.Sprintf("key banner tag=%v feature=%v", k.TagId, k.FeatureId)
}

func NewBanner(cache *ristretto.Cache, ttl time.Duration) Banner {
	return Banner{cache: cache, ttl: ttl}
}

func (mem Banner) Get(key any) (entity.Banner, error) {
	obj, ok := mem.cache.Get(key)
	if !ok {
		return entity.Banner{}, errors.New("not found")
	}
	return *obj.(*entity.Banner), nil
}

func (mem Banner) Delete(id int) error {
	banner, err := mem.Get(id)
	if err != nil {
		return err
	}
	mem.cache.Del(banner.Id)
	mem.deleteKeys(banner)
	mem.cache.Del(id)
	return nil
}

func (mem Banner) deleteKeys(banner entity.Banner) {
	for _, tag := range banner.Tags {
		key := Key{FeatureId: tag.FeatureId, TagId: tag.TagId}.String()
		mem.cache.Del(key)
	}
}

func (mem Banner) Set(banner entity.Banner) {
	for _, tag := range banner.Tags {
		key := Key{FeatureId: tag.FeatureId, TagId: tag.TagId}.String()
		mem.cache.SetWithTTL(key, &banner, 1, mem.ttl)
	}
	mem.cache.SetWithTTL(banner.Id, &banner, 1, mem.ttl)
}
