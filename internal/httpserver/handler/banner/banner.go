package banner

import (
	"backend-trainee-assignment-2024/internal/httpserver/middleware/auth"
	"backend-trainee-assignment-2024/internal/model"
	"backend-trainee-assignment-2024/internal/usecase"
	"bytes"
	"database/sql"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"io"
	"net/http"
	"strconv"
)

type Router struct {
	banner usecase.Banner
}

func New(r chi.Router, useCase usecase.Banner) {
	banner := Router{banner: useCase}

	r.Use(auth.NewRequiredMiddleware())
	r.Get("/user_banner", banner.getUserBanner)
	r.Get("/banner", banner.get)

	r.Group(func(r chi.Router) {
		r.Use(auth.RoleRequired("admin"))

		r.Post("/banner", banner.create)
		r.Patch("/banner/{id}", banner.update)
		r.Delete("/banner/{id}", banner.deleteById)
	})
}

type Banner struct {
	TagIds    []int  `json:"tag_ids"`
	FeatureId int    `json:"feature_id"`
	Content   string `json:"content" validate:"json"`
	IsActive  bool   `json:"is_active"`
}

type Create struct {
	TagIds    []int `json:"tag_ids"`
	FeatureId int   `json:"feature_id"`
	Content   any   `json:"content"`
	IsActive  bool  `json:"is_active"`
}

func (b Router) parseBanner(reader io.ReadCloser) (Banner, error) {
	request := new(Create)
	err := render.DecodeJSON(reader, request)
	if err != nil {
		return Banner{}, err
	}

	content := &bytes.Buffer{}
	enc := json.NewEncoder(content)
	if err := enc.Encode(request.Content); err != nil {
		return Banner{}, err
	}

	banner := Banner{
		IsActive:  request.IsActive,
		Content:   content.String(),
		TagIds:    request.TagIds,
		FeatureId: request.FeatureId,
	}
	return banner, nil
}

func (b Router) parsePage(r *http.Request) model.Page {
	page := model.Page{}

	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil {
			page.Limit = sql.NullInt32{Int32: int32(limit), Valid: true}
		}
	}

	if offsetStr := r.URL.Query().Get("offset"); offsetStr != "" {
		if offset, err := strconv.Atoi(offsetStr); err == nil {
			page.Offset = sql.NullInt32{Int32: int32(offset), Valid: true}
		}
	}

	return page
}
