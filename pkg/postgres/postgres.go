package postgres

import (
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/extra/bundebug"
)

type Postgres struct {
	DB *bun.DB
}

func New(url string, opts ...bun.DBOption) (*Postgres, error) {
	sqldb, err := open(url)
	if err != nil {
		return nil, err
	}

	db := bun.NewDB(sqldb, pgdialect.New(), opts...)

	db.AddQueryHook(bundebug.NewQueryHook(bundebug.WithVerbose(true)))

	initDB(db)
	err = migrate(db)
	if err != nil {
		return nil, err
	}

	pg := &Postgres{
		DB: db,
	}

	return pg, nil
}

func (p Postgres) Close() error {
	return p.DB.Close()
}
