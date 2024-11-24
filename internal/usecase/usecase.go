//go:generate mockgen -source ./mocks/api.go -destination=./mocks/song_mock.go -package=mock_repository
package usecase

import (
	"context"
	"github.com/khostya/effective-mobile/internal/dto"
	"github.com/khostya/effective-mobile/internal/repo"
)

type (
	transactionManager interface {
		RunRepeatableRead(ctx context.Context, fx func(ctxTX context.Context) error) error
		Unwrap(err error) error
	}

	infoSong interface {
		GetInfo(ctx context.Context, param dto.GetSongInfo) (*dto.SongDetail, error)
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
		Song: NewSongUseCase(SongDeps{
			SongRepo:  pg.Song,
			GroupRepo: pg.Group,
			Tm:        deps.Transactor,
			InfoSong:  deps.Client,
		}),
	}
}
