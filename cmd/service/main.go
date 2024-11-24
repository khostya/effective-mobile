package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"
	"github.com/google/uuid"
	"github.com/khostya/effective-mobile/pkg/api"
	"github.com/khostya/effective-mobile/pkg/httpserver"
	"net/http"
	"time"
)

const port = "8079"

func main() {
	router := chi.NewRouter()
	router.Use(cors.AllowAll().Handler)
	router.Get("/info", server{}.GetInfo)

	httpserver := httpserver.New(router,
		httpserver.Port(port),
	)

	httpserver.Start()
	<-httpserver.Notify()
}

type server struct {
}

func (s server) GetInfo(w http.ResponseWriter, r *http.Request) {
	releaseDate := time.Now().Format(time.DateOnly)

	render.Status(r, http.StatusOK)
	render.JSON(w, r, api.SongDetail{
		ReleaseDate: releaseDate,
		Text:        uuid.NewString(),
		Link:        uuid.NewString(),
	})
}
