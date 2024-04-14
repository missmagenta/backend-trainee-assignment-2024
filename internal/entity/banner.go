package entity

import (
	"context"
	"database/sql"
	"github.com/uptrace/bun"
	"time"
)

type Banner struct {
	bun.BaseModel `bun:"table:banners.banners"`

	Id        int          `json:"id" bun:",pk,autoincrement"`
	Tags      []Tag        `json:"tags" bun:"rel:has-many,join:id=banner_id"`
	FeatureId int          `json:"feature_id"`
	Content   string       `json:"content"`
	IsActive  sql.NullBool `json:"is_active"`

	CreatedAt *time.Time `json:"created_at" bun:",nullzero"`
	UpdatedAt *time.Time `json:"updated_at" bun:",nullzero"`
}

func (b *Banner) BeforeAppendModel(ctx context.Context, query bun.Query) error {
	now := time.Now()

	if _, ok := query.(*bun.InsertQuery); ok {
		if b.CreatedAt == nil {
			b.CreatedAt = &now
		}
		if b.UpdatedAt == nil {
			b.UpdatedAt = &now
		}
	} else if _, ok := query.(*bun.UpdateQuery); ok {
		if b.UpdatedAt == nil {
			b.UpdatedAt = &now
		}
	}

	return nil
}
