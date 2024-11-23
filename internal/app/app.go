package app

import (
	"context"
	"github.com/khostya/effective-mobile/internal/config"
	"github.com/khostya/effective-mobile/internal/http"
	"github.com/khostya/effective-mobile/internal/usecase"
	"github.com/khostya/effective-mobile/pkg/postgres"
)

func Run(ctx context.Context, cfg config.Config) error {
	pool, err := postgres.NewPool(ctx, cfg.PG.URL)
	if err != nil {
		return err
	}
	defer pool.Close()

	deps, err := newDependencies(cfg.APIEndpoint, pool)
	if err != nil {
		return err
	}

	useCases := usecase.NewUseCases(deps)

	return <-http.MustRun(
		ctx,
		cfg.HTTP,
		http.UseCases{
			Song: useCases.Song,
		},
	)
}
