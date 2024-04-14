package postgres

import (
	"backend-trainee-assignment-2024/pkg/postgres"
)

type Repositories struct {
	Banner Banner
}

func NewRepositories(pg *postgres.Postgres) Repositories {
	return Repositories{
		Banner: NewBanner(pg),
	}
}
