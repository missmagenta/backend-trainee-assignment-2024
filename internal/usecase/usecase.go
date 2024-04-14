package usecase

import (
	"backend-trainee-assignment-2024/internal/usecase/repo/memory"
	"backend-trainee-assignment-2024/internal/usecase/repo/postgres"
)

type Dependencies struct {
	Pg     postgres.Repositories
	Memory memory.Repositories
}

type UseCases struct {
	Banner Banner
	Deps   Dependencies
}

func NewUseCases(deps Dependencies) UseCases {
	pg := deps.Pg
	memory := deps.Memory

	return UseCases{
		Deps:   deps,
		Banner: NewBannerUseCase(pg.Banner, memory.Banner),
	}
}
