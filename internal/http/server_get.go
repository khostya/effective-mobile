package http

import (
	"github.com/khostya/effective-mobile/internal/dto"
	"net/http"
)

func (s *server) Get(w http.ResponseWriter, r *http.Request) {
	page, err := parsePage(r)
	if err != nil {
		s.error(w, r, http.StatusBadRequest, err)
		return
	}

	song := r.URL.Query().Get(songParam)
	group := r.URL.Query().Get(groupParam)
	link := r.URL.Query().Get(linkParam)

	songs, err := s.useCases.Song.Get(r.Context(), dto.GetSongsParam{
		Page:  &page,
		Song:  song,
		Group: group,
		Link:  link,
	})
	if err != nil {
		s.error(w, r, http.StatusInternalServerError, err)
		return
	}

	s.json(w, r, http.StatusOK, songs)
}
