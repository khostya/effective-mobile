package http

import (
	"errors"
	"github.com/go-chi/render"
	"github.com/khostya/effective-mobile/internal/repo/repoerr"
	"net/http"
)

func (s *server) json(w http.ResponseWriter, r *http.Request, status int, resp interface{}) {
	render.Status(r, status)
	render.JSON(w, r, resp)
}

func (s *server) error(w http.ResponseWriter, r *http.Request, status int, err error) {
	switch {
	case status == http.StatusBadRequest:
		w.WriteHeader(status)
	case errors.Is(err, repoerr.ErrNotFound):
		w.WriteHeader(http.StatusNotFound)
	default:
		s.internalServerError(w, r, err)
	}
}

const (
	retryAfterInSec = "30"
)

func (s *server) internalServerError(w http.ResponseWriter, _ *http.Request, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Header().Set("Retry-After", retryAfterInSec)
	if err != nil {
		w.Write([]byte(err.Error()))
	}
}
