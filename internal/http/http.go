package http

import (
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/khostya/effective-mobile/internal/config"
	"github.com/khostya/effective-mobile/internal/dto"
	"github.com/khostya/effective-mobile/pkg/httpserver"
	"log"
	"log/slog"
	"net/http"
	"net/http/httputil"
	"strconv"
)

func MustRun(ctx context.Context, cfg config.HTTP, useCases UseCases) <-chan error {
	httpserver, err := newHttpServer(ctx, cfg, useCases)
	if err != nil {
		res := make(chan error, 2)
		res <- err
		return res
	}

	httpserver.Start()
	slog.Info("started")
	return httpserver.Notify()
}

func newHttpServer(ctx context.Context, cfg config.HTTP, useCases UseCases) (*httpserver.Server, error) {
	server, err := newServer(useCases)
	if err != nil {
		return nil, err
	}

	router := getRouter(server)
	httpserver := httpserver.New(
		router,
		httpserver.Port(cfg.Port),
		httpserver.IdleTimeout(cfg.IdleTimeout),
		httpserver.MaxHeaderBytes(cfg.MaxHeaderBytes),
		httpserver.WriteTimeout(cfg.WriteTimeout),
	)

	go func() {
		<-ctx.Done()
		defer slog.Info("stopped")
		if err := httpserver.Shutdown(); err != nil {
			log.Fatalf("HTTP handler Shutdown: %s", err)
		}
	}()

	return httpserver, nil
}

func getRouter(server *server) chi.Router {
	router := chi.NewRouter()
	router.Use(logging)
	router.Use(cors.AllowAll().Handler)

	router.Delete("/{id}", server.DeleteId)
	router.Post("/", server.PostCreate)
	router.Put("/{id}", server.Put)
	router.Get("/verse/{id}", server.GetVerseId)
	router.Get("/", server.Get)

	return router
}

func logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqDump, err := httputil.DumpRequest(r, true)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		slog.Debug(string(reqDump))
		next.ServeHTTP(w, r)
	})
}

const (
	pageParam           = "page"
	sizeParam           = "size"
	idParam             = "id"
	songParam           = "song"
	groupParam          = "group"
	linkParam           = "link"
	releaseDateGteParam = "release_date_gte"
	releaseDateLteParam = "release_date_lte"
)

func parsePage(r *http.Request) (dto.Page, error) {
	page := r.URL.Query().Get(pageParam)

	pageInt, err := strconv.ParseUint(page, 10, 32)
	if err != nil {
		return dto.Page{}, err
	}

	size := r.URL.Query().Get(sizeParam)
	sizeInt, err := strconv.ParseUint(size, 10, 32)
	if err != nil {
		return dto.Page{}, err
	}

	return dto.Page{
		Page: uint(pageInt),
		Size: uint(sizeInt),
	}, nil
}
