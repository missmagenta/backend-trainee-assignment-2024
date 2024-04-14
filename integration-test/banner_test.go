package main

import (
	"backend-trainee-assignment-2024/internal/httpserver/handler/banner"
	"bytes"
	"encoding/json"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/gavv/httpexpect/v2"
	"net/http"
	"testing"
)

type Id struct {
	Id int `json:"id"`
}

func TestBanner(t *testing.T) {
	e := httpexpect.Default(t, basePath)

	banners, _ := createBanner(e)
	if t.Failed() {
		return
	}

	getBanner(e, t, banners)
}

func createBanner(e *httpexpect.Expect) (banner.Create, Id) {
	json := banner.Create{
		IsActive: true,
		TagIds: []int{
			gofakeit.Number(0, 4),
			gofakeit.Number(5, 10),
		},
		Content:   map[string]string{"user": "aboba"},
		FeatureId: gofakeit.Number(0, 1000000000),
	}

	id := new(Id)
	e.POST("/banner").
		WithJSON(json).
		WithHeader("token", "admin").
		Expect().
		Status(http.StatusCreated).
		JSON().
		Decode(id)

	return json, *id
}

func getBanner(e *httpexpect.Expect, t *testing.T, banner banner.Create) {
	content := &bytes.Buffer{}
	enc := json.NewEncoder(content)
	if err := enc.Encode(banner.Content); err != nil {
		t.Fail()
		return
	}

	e.GET("/user_banner").
		WithQuery("tag_id", banner.TagIds[0]).
		WithQuery("feature_id", banner.FeatureId).
		WithHeader("token", "admin").
		Expect().
		Status(http.StatusOK).
		JSON().
		String().
		IsEqual(content.String())
}
