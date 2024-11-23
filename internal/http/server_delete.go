package http

import (
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"net/http"
)

func (s *server) DeleteId(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	songID, err := uuid.Parse(id)
	if err != nil {
		s.error(w, r, http.StatusBadRequest, err)
		return
	}

	err = s.useCases.Song.DeleteByID(r.Context(), songID)
	if err != nil {
		s.error(w, r, http.StatusInternalServerError, err)
		return
	}
	w.WriteHeader(http.StatusOK)
}
