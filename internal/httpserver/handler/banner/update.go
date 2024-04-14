package banner

import (
	"backend-trainee-assignment-2024/internal/entity"
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

func (b Router) update(w http.ResponseWriter, r *http.Request) {
	request, err := b.parseBanner(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	banner := entity.Banner{
		Id:        id,
		Content:   request.Content,
		IsActive:  sql.NullBool{Valid: true, Bool: request.IsActive},
		FeatureId: request.FeatureId,
	}
	for _, tagId := range request.TagIds {
		tag := entity.Tag{TagId: tagId, FeatureId: banner.FeatureId, BannerId: id}
		banner.Tags = append(banner.Tags, tag)
	}

	updatedId, err := b.banner.Update(r.Context(), banner)
	if errors.Is(err, errors.New("not found")) {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response := struct {
		UpdatedID int `json:"updated_id"`
	}{
		UpdatedID: updatedId,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
