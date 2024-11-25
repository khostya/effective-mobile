//go:build integration

package postgres

import (
	"context"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	"github.com/khostya/effective-mobile/internal/domain"
	"github.com/khostya/effective-mobile/internal/dto"
	"github.com/khostya/effective-mobile/internal/repo"
	"github.com/khostya/effective-mobile/internal/repo/repoerr"
	"github.com/khostya/effective-mobile/internal/repo/transactor"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

type SongsTestSuite struct {
	suite.Suite
	ctx        context.Context
	songRepo   repo.Song
	groupRepo  repo.Group
	transactor *transactor.TransactionManager
}

func TestSongs(t *testing.T) {
	suite.Run(t, new(SongsTestSuite))
}

func (s *SongsTestSuite) SetupSuite() {
	s.transactor = transactor.NewTransactionManager(db.GetPool())
	s.songRepo = repo.NewSongRepo(s.transactor)
	s.groupRepo = repo.NewGroupRepo(s.transactor)
	s.ctx = context.Background()
}

func (s *SongsTestSuite) TestCreate() {
	_ = s.create(gofakeit.UUID(), gofakeit.UUID(), time.Now())
}

func (s *SongsTestSuite) TestCreateDuplicateError() {
	song := s.create(gofakeit.UUID(), gofakeit.UUID(), time.Now())

	err := s.songRepo.Create(s.ctx, song)
	require.ErrorIs(s.T(), err, repoerr.ErrDuplicate)
}

func (s *SongsTestSuite) TestGetByID() {
	song := s.create(gofakeit.UUID(), gofakeit.UUID(), time.Now())

	actual, err := s.songRepo.GetByID(s.ctx, song.ID)

	require.NoError(s.T(), err)
	require.EqualExportedValues(s.T(), song, actual)
}

func (s *SongsTestSuite) TestDelete() {
	song := s.create(gofakeit.UUID(), gofakeit.UUID(), time.Now())

	err := s.songRepo.Delete(s.ctx, song.ID)
	require.NoError(s.T(), err)

	_, err = s.songRepo.GetByID(s.ctx, song.ID)
	require.ErrorIs(s.T(), err, repoerr.ErrNotFound)
}

func (s *SongsTestSuite) TestDeleteNotFound() {
	err := s.songRepo.Delete(s.ctx, uuid.New())
	require.ErrorIs(s.T(), err, repoerr.ErrNotFound)
}

func (s *SongsTestSuite) TestUpdate() {
	song := s.create(gofakeit.UUID(), gofakeit.UUID(), time.Now())

	randomString := uuid.New().String()
	err := s.songRepo.Update(s.ctx, dto.UpdateSongParam{
		ID:     song.ID,
		Song:   &randomString,
		Verses: []string{randomString},
		Link:   &randomString,
	})
	require.NoError(s.T(), err)

	song.Song = randomString
	song.Link = randomString
	song.Verses = []string{randomString}

	actual, err := s.songRepo.GetByID(s.ctx, song.ID)
	require.NoError(s.T(), err)
	require.EqualExportedValues(s.T(), song, actual)

	song = actual

	var empty string
	err = s.songRepo.Update(s.ctx, dto.UpdateSongParam{
		ID:   song.ID,
		Song: &empty,
		Link: &empty,
	})
	require.NoError(s.T(), err)

	song.Link = empty

	actual, err = s.songRepo.GetByID(s.ctx, song.ID)
	require.NoError(s.T(), err)
	require.EqualExportedValues(s.T(), song, actual)
}

func (s *SongsTestSuite) TestUpdateNotFound() {
	randomString := uuid.New().String()
	err := s.songRepo.Update(s.ctx, dto.UpdateSongParam{
		ID:     uuid.New(),
		Song:   &randomString,
		Verses: []string{randomString},
		Link:   &randomString,
	})
	require.ErrorIs(s.T(), err, repoerr.ErrNotFound)
}

func (s *SongsTestSuite) TestGet() {
	truncate()

	_ = s.create("song1", "group", time.Now())
	song2 := s.create("song2", "group", time.Now().Add(time.Hour))

	songs, err := s.songRepo.Get(s.ctx, dto.GetSongsParam{
		Group: "group",
		Page: &dto.Page{
			Page: 1,
			Size: 10,
		},
		Link: url,
	})
	require.NoError(s.T(), err)
	require.Len(s.T(), songs, 2)

	songs, err = s.songRepo.Get(s.ctx, dto.GetSongsParam{
		Song: "song1",
		Page: &dto.Page{
			Page: 1,
			Size: 10,
		},
	})
	require.NoError(s.T(), err)
	require.Len(s.T(), songs, 1)

	t := time.Now().Add(time.Hour * 24 * 365)
	songs, err = s.songRepo.Get(s.ctx, dto.GetSongsParam{
		ReleaseDateGte: &t,
	})
	require.NoError(s.T(), err)
	require.Len(s.T(), songs, 0)

	songs, err = s.songRepo.Get(s.ctx, dto.GetSongsParam{
		ReleaseDateLte: &t,
	})
	require.NoError(s.T(), err)
	require.Len(s.T(), songs, 2)

	songs, err = s.songRepo.Get(s.ctx, dto.GetSongsParam{
		ReleaseDateGte: &song2.ReleaseDate,
	})
	require.NoError(s.T(), err)
	require.Len(s.T(), songs, 1)
}

func (s *SongsTestSuite) create(songTitle string, groupTitle string, releaseDate time.Time) domain.Song {
	group := s.createGroup(groupTitle)

	song := NewSong(songTitle, group, releaseDate)

	err := s.songRepo.Create(s.ctx, song)
	require.NoError(s.T(), err)

	return song
}

func (s *SongsTestSuite) createGroup(groupTitle string) domain.Group {
	group := NewGroup(groupTitle)

	err := s.groupRepo.CreateOnConflictDoNothing(s.ctx, group)
	require.NoError(s.T(), err)

	return group
}
