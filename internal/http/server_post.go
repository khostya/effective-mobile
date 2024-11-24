package http

import (
	"github.com/go-chi/render"
	"github.com/khostya/effective-mobile/internal/dto"
	"github.com/khostya/effective-mobile/internal/http/api"
	"net/http"
)

func (s *server) PostCreate(w http.ResponseWriter, r *http.Request) {
	var req api.PostJSONBody
	err := render.DecodeJSON(r.Body, &req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if req.Group == "" || req.Song == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = s.useCases.Song.Create(r.Context(), dto.CreateSongParam{
		Song:  req.Song,
		Group: req.Group,
	})
	if err != nil {
		s.error(w, r, http.StatusInternalServerError, err)
		return
	}
}
