package http

import (
	"github.com/khostya/effective-mobile/internal/dto"
	"net/http"
	"time"
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

	releaseDateGte, err := parseReleaseDate(r.URL.Query().Get(releaseDateGteParam))
	if err != nil {
		s.error(w, r, http.StatusBadRequest, err)
		return
	}

	releaseDateLte, err := parseReleaseDate(r.URL.Query().Get(releaseDateLteParam))
	if err != nil {
		s.error(w, r, http.StatusBadRequest, err)
		return
	}

	songs, err := s.useCases.Song.Get(r.Context(), dto.GetSongsParam{
		ReleaseDateLte: releaseDateLte,
		ReleaseDateGte: releaseDateGte,
		Page:           &page,
		Song:           song,
		Group:          group,
		Link:           link,
	})
	if err != nil {
		s.error(w, r, http.StatusInternalServerError, err)
		return
	}

	s.json(w, r, http.StatusOK, songs)
}

func parseReleaseDate(q string) (*time.Time, error) {
	if q == "" {
		return nil, nil
	}

	var date *time.Time
	v, err := time.Parse(time.DateOnly, q)
	if err != nil {
		return date, err
	}
	date = &v
	return date, nil
}
