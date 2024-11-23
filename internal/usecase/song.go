//go:generate mockgen -source ./mocks/song.go -destination=./mocks/song_mock.go -package=mock_usecase
package usecase

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/khostya/effective-mobile/internal/domain"
	"github.com/khostya/effective-mobile/internal/dto"
	"github.com/khostya/effective-mobile/pkg/api"
	"net/http"
	"strings"
	"time"
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
	resp, err := uc.infoSong.GetInfo(ctx, &api.GetInfoParams{
		Song:  param.Song,
		Group: param.Group,
	})
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	info, err := api.ParseGetInfoResponse(resp)
	if err != nil {
		return err
	}
	if info.StatusCode() != http.StatusOK {
		return errors.New("internal error")
	}

	releaseDate, err := time.Parse(time.RFC3339, info.JSON200.ReleaseDate)
	if err != nil {
		return errors.New("invalid release date")
	}

	return uc.tm.RunRepeatableRead(ctx, func(ctx context.Context) error {
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
			Link:        info.JSON200.Link,
			Text:        info.JSON200.Text,
			ReleaseDate: releaseDate,
		})
	})
}

func (uc Song) Get(ctx context.Context, param dto.GetSongsParam) ([]domain.Song, error) {
	return uc.songRepo.Get(ctx, param)
}

func (uc Song) GetByVerse(ctx context.Context, id uuid.UUID, page dto.Page) ([]string, error) {
	song, err := uc.songRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	verses := strings.Split(song.Text, "\n\n")
	left := page.Offset()
	if left >= uint(len(verses)) {
		return nil, errors.New("out of range")
	}

	right := int(left + page.Size)
	return verses[left:min(len(verses), right)], nil
}

func (uc Song) Update(ctx context.Context, param dto.UpdateSongParam) error {
	return uc.songRepo.Update(ctx, param)
}

func NewUserUseCase(songRepo songRepo, groupRepo groupRepo, infoSong infoSong, tm transactionManager) Song {
	return Song{
		songRepo:  songRepo,
		infoSong:  infoSong,
		groupRepo: groupRepo,
		tm:        tm,
	}
}
