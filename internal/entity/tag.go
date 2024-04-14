package entity

import "github.com/uptrace/bun"

type Tag struct {
	bun.BaseModel `bun:"table:banners.tags"`
	BannerId      int `json:"banner_id"`
	FeatureId     int `json:"feature_id" bun:",pk"`
	TagId         int `json:"tag_id" bun:",pk"`
}
