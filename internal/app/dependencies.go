package app

import (
	"github.com/khostya/effective-mobile/internal/repo"
	"github.com/khostya/effective-mobile/internal/repo/transactor"
	"github.com/khostya/effective-mobile/internal/usecase"
	"github.com/khostya/effective-mobile/internal/usecase/api"
	httpapi "github.com/khostya/effective-mobile/pkg/api"
	"github.com/khostya/effective-mobile/pkg/postgres"
)

func newDependencies(apiEndpoint string, pool *postgres.Pool) (usecase.Dependencies, error) {
	transactor := transactor.NewTransactionManager(pool)
	pgRepositories := repo.NewRepositories(transactor)

	client, err := httpapi.NewClient(apiEndpoint)
	if err != nil {
		return usecase.Dependencies{}, err
	}

	return usecase.Dependencies{
		Pg:         pgRepositories,
		Transactor: transactor,
		Client:     api.NewSongInfo(client),
	}, nil
}
