package usecase

import (
	"context"
	"github.com/khostya/effective-mobile/internal/repo"
	"github.com/khostya/effective-mobile/pkg/api"
	"net/http"
)

type (
	transactionManager interface {
		RunRepeatableRead(ctx context.Context, fx func(ctxTX context.Context) error) error
		Unwrap(err error) error
	}

	infoSong interface {
		GetInfo(ctx context.Context, params *api.GetInfoParams, reqEditors ...api.RequestEditorFn) (*http.Response, error)
	}

	Dependencies struct {
		Pg         repo.Repositories
		Transactor transactionManager
		Client     infoSong
	}

	UseCases struct {
		Deps Dependencies
		Song Song
	}
)

func NewUseCases(deps Dependencies) UseCases {
	pg := deps.Pg

	return UseCases{
		Deps: deps,
		Song: NewUserUseCase(pg.Song, pg.Group, deps.Client, deps.Transactor),
	}
}
