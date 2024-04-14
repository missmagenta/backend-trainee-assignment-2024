package banner

import (
	"errors"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

func (b Router) deleteById(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	bannerId, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = b.banner.DeleteById(r.Context(), bannerId)
	if errors.Is(err, errors.New("not found")) {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
