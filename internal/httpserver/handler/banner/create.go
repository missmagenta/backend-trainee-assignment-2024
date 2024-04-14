package banner

import (
	"backend-trainee-assignment-2024/internal/entity"
	"database/sql"
	"encoding/json"
	"net/http"
)

func (b Router) create(w http.ResponseWriter, r *http.Request) {
	request, err := b.parseBanner(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	banner := entity.Banner{
		IsActive:  sql.NullBool{Valid: true, Bool: request.IsActive},
		FeatureId: request.FeatureId,
		Content:   request.Content,
	}

	for _, tagId := range request.TagIds {
		tag := entity.Tag{TagId: tagId, FeatureId: request.FeatureId}
		banner.Tags = append(banner.Tags, tag)
	}

	id, err := b.banner.Create(r.Context(), banner)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := struct {
		BannerId int `json:"banner_id"`
	}{BannerId: id}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}
