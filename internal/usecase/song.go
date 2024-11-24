//go:generate mockgen -source ./mocks/song.go -destination=./mocks/song_mock.go -package=mock_usecase
package usecase

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/khostya/effective-mobile/internal/domain"
	"github.com/khostya/effective-mobile/internal/dto"
)

type (
	songRepo interface {
		Create(ctx context.Context, song domain.Song) error
		GetByID(ctx context.Context, id uuid.UUID) (domain.Song, error)
		Update(ctx context.Context, param dto.UpdateSongParam) error
		Delete(ctx context.Context, id uuid.UUID) error
		Get(ctx context.Context, param dto.GetSongsParam) ([]domain.Song, error)
	}

	groupRepo interface {
		CreateOnConflictDoNothing(ctx context.Context, group domain.Group) error
	}

	SongDeps struct {
		SongRepo  songRepo
		GroupRepo groupRepo
		InfoSong  infoSong
		Tm        transactionManager
	}

	Song struct {
		songRepo  songRepo
		groupRepo groupRepo
		infoSong  infoSong
		tm        transactionManager
	}
)

func (uc Song) DeleteByID(ctx context.Context, id uuid.UUID) error {
	return uc.songRepo.Delete(ctx, id)
}

func (uc Song) Create(ctx context.Context, param dto.CreateSongParam) error {
	info, err := uc.infoSong.GetInfo(ctx, dto.GetSongInfo{
		Song:  param.Song,
		Group: param.Group,
	})
	if err != nil {
		return err
	}

	return uc.tm.Unwrap(uc.tm.RunRepeatableRead(ctx, func(ctx context.Context) error {
		group := domain.Group{
			Title: param.Group,
		}

		err := uc.groupRepo.CreateOnConflictDoNothing(ctx, group)
		if err != nil {
			return err
		}

		return uc.songRepo.Create(ctx, domain.Song{
			ID:          uuid.New(),
			Song:        param.Song,
			Group:       &group,
			Link:        info.Link,
			Text:        domain.Text(info.Text),
			ReleaseDate: info.ReleaseDate,
		})
	}))
}

func (uc Song) Get(ctx context.Context, param dto.GetSongsParam) ([]domain.Song, error) {
	return uc.songRepo.Get(ctx, param)
}

func (uc Song) GetByVerse(ctx context.Context, param dto.GetByVerseParam) ([]string, error) {
	song, err := uc.songRepo.GetByID(ctx, param.ID)
	if err != nil {
		return nil, err
	}

	verses, err := song.Text.GetVerse(param.Page)
	if err == nil {
		return verses, nil
	}

	if errors.Is(err, domain.ErrOutOfRange) {
		return []string{}, ErrOutOfRange
	}
	return nil, err
}

func (uc Song) Update(ctx context.Context, param dto.UpdateSongParam) error {
	return uc.songRepo.Update(ctx, param)
}

func NewSongUseCase(deps SongDeps) Song {
	return Song{
		songRepo:  deps.SongRepo,
		infoSong:  deps.InfoSong,
		groupRepo: deps.GroupRepo,
		tm:        deps.Tm,
	}
}
