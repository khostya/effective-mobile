package http

import (
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/khostya/effective-mobile/internal/dto"
	"github.com/khostya/effective-mobile/internal/http/api"
	"net/http"
)

func (s *server) GetVerseId(w http.ResponseWriter, r *http.Request) {
	page, err := parsePage(r)
	if err != nil {
		s.error(w, r, http.StatusBadRequest, err)
		return
	}

	id := chi.URLParam(r, idParam)
	songID, err := uuid.Parse(id)
	if err != nil {
		s.error(w, r, http.StatusBadRequest, err)
		return
	}

	text, err := s.useCases.Song.GetByVerse(r.Context(), dto.GetByVerseParam{
		ID:   songID,
		Page: page,
	})
	if err != nil {
		s.error(w, r, http.StatusInternalServerError, err)
		return
	}
	s.json(w, r, http.StatusOK, api.GetVerseResponse{Text: text})
}
