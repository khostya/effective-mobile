package http

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/google/uuid"
	"github.com/khostya/effective-mobile/internal/dto"
	"github.com/khostya/effective-mobile/internal/http/api"
	"net/http"
	"strings"
)

func (s *server) Put(w http.ResponseWriter, r *http.Request) {
	var req api.PutIdJSONBody
	err := render.DecodeJSON(r.Body, &req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	songID, err := uuid.Parse(chi.URLParam(r, idParam))
	if err != nil {
		s.error(w, r, http.StatusBadRequest, err)
		return
	}

	var text = make([]string, 0)
	if req.Text != nil {
		text = strings.Split(*req.Text, "\n\n")
	}

	err = s.useCases.Song.Update(r.Context(), dto.UpdateSongParam{
		ID:     songID,
		Song:   req.Song,
		Link:   req.Link,
		Verses: text,
	})
	if err != nil {
		s.error(w, r, http.StatusInternalServerError, err)
		return
	}
}
