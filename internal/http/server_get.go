package http

import (
	"github.com/khostya/effective-mobile/internal/dto"
	"net/http"
)

func (s *server) Get(w http.ResponseWriter, r *http.Request) {
	page, err := s.parsePage(r)
	if err != nil {
		s.error(w, r, http.StatusBadRequest, err)
		return
	}

	var (
		song  string
		group string
	)

	if r.URL.Query().Has("song") {
		song = r.URL.Query().Get("song")
	}
	if r.URL.Query().Has("group") {
		group = r.URL.Query().Get("group")
	}

	songs, err := s.useCases.Song.Get(r.Context(), dto.GetSongsParam{
		Page:  &page,
		Song:  song,
		Group: group,
	})
	if err != nil {
		s.error(w, r, http.StatusInternalServerError, err)
		return
	}

	s.json(w, r, http.StatusOK, songs)
}
