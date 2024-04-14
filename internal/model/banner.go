package model

import (
	"backend-trainee-assignment-2024/internal/entity"
	"time"
)

type Banner struct {
	Id        int    `json:"id"`
	Tags      []int  `json:"tags"`
	FeatureId int    `json:"feature_id"`
	Content   string `json:"content"`
	IsActive  bool   `json:"is_active"`

	CreatedAt time.Time `json:"created_at" bun:",nullzero"`
	UpdatedAt time.Time `json:"updated_at" bun:",nullzero"`
}

func NewBanners(banners []entity.Banner) []Banner {
	converted := make([]Banner, len(banners))
	for i, banner := range banners {
		converted[i] = Banner{
			Id:        banner.Id,
			Tags:      extractTagIds(banner.Tags),
			FeatureId: banner.FeatureId,
			Content:   banner.Content,
			IsActive:  banner.IsActive.Bool,
			CreatedAt: *banner.CreatedAt,
			UpdatedAt: *banner.UpdatedAt,
		}
	}
	return converted
}

func extractTagIds(tags []entity.Tag) []int {
	ids := make([]int, len(tags))
	for i, tag := range tags {
		ids[i] = tag.TagId
	}
	return ids
}
