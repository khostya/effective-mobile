//go:generate mockgen -source ./mocks/song.go -destination=./mocks/song_mock.go -package=mock_repository
package repo

import (
	"context"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/khostya/effective-mobile/internal/domain"
	"github.com/khostya/effective-mobile/internal/dto"
	"github.com/khostya/effective-mobile/internal/repo/exec"
	"github.com/khostya/effective-mobile/internal/repo/repoerr"
	"github.com/khostya/effective-mobile/internal/repo/schema"
	"github.com/khostya/effective-mobile/internal/repo/transactor"
)

const (
	songsTable = "effective.songs"
)

type (
	Song struct {
		queryEngineProvider transactor.QueryEngineProvider
	}
)

func (s Song) Create(ctx context.Context, song domain.Song) error {
	db := s.queryEngineProvider.GetQueryEngine(ctx)

	record := schema.NewSong(song)
	query := sq.Insert(songsTable).
		Columns(record.InsertColumns()...).
		Values(record.InsertValues()...).
		PlaceholderFormat(sq.Dollar)

	err := exec.Insert(ctx, query, db)
	if err != nil && exec.IsDuplicateKeyError(err) {
		return repoerr.ErrDuplicate
	}

	return err
}

func (s Song) Get(ctx context.Context, param dto.GetSongsParam) ([]domain.Song, error) {
	db := s.queryEngineProvider.GetQueryEngine(ctx)

	query := sq.Select(schema.Song{}.SelectColumns()...).
		From(songsTable).
		PlaceholderFormat(sq.Dollar)

	var n uint
	if param.Group != "" {
		n += 1
		query = query.Where(fmt.Sprintf("songs.group = $%v", n), param.Group)
	}

	if param.Song != "" {
		n += 1
		query = query.Where(fmt.Sprintf("song = $%v", n), param.Song)
	}

	if param.Page != nil {
		offset, err := param.Page.Offset()
		if err != nil {
			return nil, err
		}
		query = query.
			Offset(uint64(offset)).
			Limit(uint64(param.Page.Limit()))
	}

	songs, err := exec.ScanALL[schema.Song](ctx, query, db)
	if err != nil {
		return []domain.Song{}, err
	}

	return schema.NewDomainSongs(songs), nil
}

func (s Song) GetByID(ctx context.Context, id uuid.UUID) (domain.Song, error) {
	db := s.queryEngineProvider.GetQueryEngine(ctx)

	query := sq.Select(schema.Song{}.SelectColumns()...).
		From(songsTable).
		Where("id = $1", id).
		PlaceholderFormat(sq.Dollar)

	song, err := exec.ScanOne[schema.Song](ctx, query, db)
	if err != nil {
		return domain.Song{}, err
	}

	return schema.NewDomainSong(song), nil
}

func (s Song) Update(ctx context.Context, param dto.UpdateSongParam) error {
	db := s.queryEngineProvider.GetQueryEngine(ctx)

	record := schema.NewSongUpdate(param)
	query := sq.Update(songsTable).
		PlaceholderFormat(sq.Dollar)

	values := record.UpdateValues()
	columns := record.UpdateColumns()

	var n int = 1
	for i, v := range values {
		query = query.Set(columns[i], v)
		n += 1
	}

	query = query.
		Where(fmt.Sprintf("id = $%v", n), param.ID.String())

	return exec.Update(ctx, query, db)
}

func (s Song) Delete(ctx context.Context, id uuid.UUID) error {
	db := s.queryEngineProvider.GetQueryEngine(ctx)

	query := sq.Delete(songsTable).
		Where("id = $1", id).
		PlaceholderFormat(sq.Dollar)

	return exec.Delete(ctx, query, db)
}

func NewSongRepo(provider transactor.QueryEngineProvider) Song {
	return Song{provider}
}
