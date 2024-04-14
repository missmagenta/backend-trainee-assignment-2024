package banner

import (
	"backend-trainee-assignment-2024/internal/model"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
)

func (b Router) getUserBanner(w http.ResponseWriter, r *http.Request) {
	isAdmin := false
	token := r.Context().Value("token")
	if token != nil && token == "admin" {
		isAdmin = true
	}
	filter := b.parseFilter(r)
	if !filter.TagId.Valid || !filter.FeatureId.Valid {
		http.Error(w, errors.New("feature_id and tag_id is required").Error(), http.StatusBadRequest)
		return
	}

	useLastRevision, _ := strconv.ParseBool(r.URL.Query().Get("use_last_revision"))

	banner, err := b.banner.GetUserBanner(r.Context(), filter, useLastRevision, isAdmin)
	if errors.Is(err, errors.New("not found")) {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(banner.Content); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func (b Router) get(w http.ResponseWriter, r *http.Request) {
	isAdmin := false
	token := r.Context().Value("token")
	if token != nil && token == "admin" {
		isAdmin = true
	}
	filter := b.parseFilter(r)
	page := b.parsePage(r)

	banners, err := b.banner.Get(r.Context(), filter, page, isAdmin)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	bannerModels := model.NewBanners(banners)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(bannerModels); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func (b Router) parseFilter(r *http.Request) model.Filter {
	featureId, featureErr := strconv.Atoi(r.URL.Query().Get("feature_id"))

	tagId, tagErr := strconv.Atoi(r.URL.Query().Get("tag_id"))

	filter := model.Filter{
		TagId:     sql.NullInt32{Valid: tagErr == nil, Int32: int32(tagId)},
		FeatureId: sql.NullInt32{Valid: featureErr == nil, Int32: int32(featureId)},
	}

	return filter
}
