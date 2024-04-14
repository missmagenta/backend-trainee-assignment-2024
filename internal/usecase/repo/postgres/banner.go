package postgres

import (
	"backend-trainee-assignment-2024/internal/entity"
	"backend-trainee-assignment-2024/internal/model"
	"backend-trainee-assignment-2024/pkg/postgres"
	"context"
	"database/sql"
	"errors"
	"github.com/uptrace/bun"
	"strconv"
)

type Banner struct {
	db *bun.DB
}

func NewBanner(pg *postgres.Postgres) Banner {
	return Banner{db: pg.DB}
}

func (r Banner) Create(ctx context.Context, banner entity.Banner) (int, error) {
	err := r.db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		_, err := tx.NewInsert().Model(&banner).Exec(ctx)
		if err != nil {
			return err
		}

		for i := range banner.Tags {
			banner.Tags[i].BannerId = banner.Id
		}

		_, err = tx.NewInsert().Model(&banner.Tags).Exec(ctx)
		return err
	})

	return banner.Id, err
}

func (r Banner) Update(ctx context.Context, banner entity.Banner) (int, error) {
	var updatedId int
	err := r.db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		exec, err := tx.NewUpdate().
			Model(&banner).
			WherePK().
			OmitZero().
			Exec(ctx)
		if err != nil {
			return err
		}
		if n, err := exec.RowsAffected(); err == nil && n == 0 {
			return errors.New("not found")
		}

		_, err = tx.NewDelete().
			Model(&entity.Tag{}).
			Where("banner_id = ?", banner.Id).
			Exec(ctx)
		if err != nil {
			return err
		}

		_, err = tx.NewInsert().Model(&banner.Tags).Exec(ctx)
		if err != nil {
			return err
		}

		updatedId = banner.Id

		return nil
	})

	return updatedId, err
}

func (r Banner) DeleteById(ctx context.Context, id int) (int, error) {
	deletedId := id
	_, err := r.db.NewDelete().Model(&entity.Banner{}).Where("id = ?", id).Exec(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return deletedId, errors.New("not found")
		}
		return deletedId, err
	}

	return deletedId, nil
}

func (r Banner) Get(ctx context.Context, filter model.Filter, page model.Page, all bool) ([]entity.Banner, error) {
	where := ""
	if filter.TagId.Valid {
		where = "tag_id = " + strconv.FormatInt(int64(filter.TagId.Int32), 10)
	}
	if filter.FeatureId.Valid {
		if len(where) != 0 {
			where += " and "
		}
		where += "feature_id = " + strconv.FormatInt(int64(filter.FeatureId.Int32), 10)
	}

	bannersIds := r.db.NewSelect().
		Model(&entity.Tag{}).
		Column("banner_id").
		Distinct()
	if where != "" {
		bannersIds = bannersIds.Where(where)
	}

	var banners []entity.Banner
	bannersQuery := r.db.NewSelect().
		Model(&banners).
		Where("id in (?)", bannersIds).
		Relation("Tags")

	if page.Limit.Valid {
		bannersQuery = bannersQuery.Limit(int(page.Limit.Int32))
	}
	if page.Offset.Valid {
		bannersQuery = bannersQuery.Offset(int(page.Offset.Int32))
	}

	if !all {
		bannersQuery = bannersQuery.Where("is_active = true")
	}

	err := bannersQuery.Scan(ctx)
	return banners, err
}

func (r Banner) getTagsQuery(featureId int) *bun.SelectQuery {
	return r.db.NewSelect().
		Column("tag_id").
		Table("banners.tags").
		Where("feature_id = ?", featureId)
}

func (r Banner) GetForUser(ctx context.Context, filter model.Filter, all bool) (entity.Banner, error) {
	var (
		banner    = new(entity.Banner)
		tagsQuery = r.getTagsQuery(int(filter.FeatureId.Int32))
	)

	bannerQuery := r.db.NewSelect().
		Model(banner).
		Where("feature_id = ? and ? in (?)", filter.FeatureId, filter.TagId, tagsQuery)
	if !all {
		bannerQuery = bannerQuery.Where("is_active = true")
	}

	err := bannerQuery.Relation("Tags").Scan(ctx)
	if errors.Is(err, sql.ErrNoRows) {
		return entity.Banner{}, errors.New("not found")
	}
	if err != nil {
		return entity.Banner{}, err
	}

	return *banner, nil
}
