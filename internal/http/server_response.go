package http

import (
	"errors"
	"github.com/go-chi/render"
	"github.com/khostya/effective-mobile/internal/repo/repoerr"
	"github.com/khostya/effective-mobile/pkg/log/sl"
	"log/slog"
	"net/http"
)

func (s *server) json(w http.ResponseWriter, r *http.Request, status int, resp interface{}) {
	render.Status(r, status)
	render.JSON(w, r, resp)
	slog.Debug("json response", sl.URL(r.URL.String()), sl.Code(uint64(status)))
}

func (s *server) error(w http.ResponseWriter, r *http.Request, status int, err error) {
	switch {
	case status == http.StatusBadRequest:
		slog.Debug("bad request", sl.Err(err), sl.URL(r.URL.String()), sl.Code(http.StatusBadRequest))
		w.WriteHeader(status)
	case errors.Is(err, repoerr.ErrNotFound):
		slog.Debug("not found", sl.Err(err), sl.URL(r.URL.String()), sl.Code(http.StatusNotFound))
		w.WriteHeader(http.StatusNotFound)
	default:
		s.internalServerError(w, r, err)
	}
}

func (s *server) internalServerError(w http.ResponseWriter, r *http.Request, err error) {
	slog.Error("internal server error", sl.Err(err), sl.URL(r.URL.String()), sl.Code(http.StatusInternalServerError))
	w.WriteHeader(http.StatusInternalServerError)
	if err != nil {
		w.Write([]byte(err.Error()))
	}
}
