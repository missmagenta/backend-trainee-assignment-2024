package usecase

import (
	"backend-trainee-assignment-2024/internal/entity"
	"backend-trainee-assignment-2024/internal/model"
	"backend-trainee-assignment-2024/internal/usecase/repo/memory"
	"backend-trainee-assignment-2024/internal/usecase/repo/postgres"
	"context"
	"math/rand"
)

type Banner struct {
	pg  postgres.Banner
	mem memory.Banner
}

func NewBannerUseCase(pg postgres.Banner, memory memory.Banner) Banner {
	return Banner{pg, memory}
}

func (uc Banner) Create(ctx context.Context, banner entity.Banner) (int, error) {
	id, err := uc.pg.Create(ctx, banner)
	if err != nil {
		return 0, err
	}
	banner.Id = id
	uc.mem.Set(banner)
	return id, err
}

func (uc Banner) Update(ctx context.Context, banner entity.Banner) (int, error) {
	return uc.pg.Update(ctx, banner)
}

func (uc Banner) DeleteById(ctx context.Context, id int) (int, error) {
	deletedId, err := uc.pg.DeleteById(ctx, id)
	if err != nil {
		return deletedId, err
	}
	uc.mem.Delete(id)
	return deletedId, nil
}

func (uc Banner) GetUserBanner(ctx context.Context, filter model.Filter, useLastRevision bool, isAdmin bool) (entity.Banner, error) {
	if useLastRevision || rand.Intn(10) == 0 {
		banner, err := uc.pg.GetForUser(ctx, filter, isAdmin)
		if err != nil {
			return banner, err
		}
		uc.mem.Set(banner)
		return banner, nil
	}

	key := memory.Key{TagId: int(filter.TagId.Int32), FeatureId: int(filter.FeatureId.Int32)}.String()
	banner, err := uc.mem.Get(key)
	if err == nil {
		return banner, nil
	}

	banner, err = uc.pg.GetForUser(ctx, filter, isAdmin)
	if err != nil {
		return entity.Banner{}, err
	}
	uc.mem.Set(banner)
	return banner, nil
}

func (uc Banner) Get(ctx context.Context, filter model.Filter, page model.Page, isAdmin bool) ([]entity.Banner, error) {
	return uc.pg.Get(ctx, filter, page, isAdmin)
}
