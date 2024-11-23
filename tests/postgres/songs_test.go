//go:build integration

package postgres

import (
	"context"
	"github.com/google/uuid"
	"github.com/khostya/effective-mobile/internal/domain"
	"github.com/khostya/effective-mobile/internal/dto"
	"github.com/khostya/effective-mobile/internal/repo"
	"github.com/khostya/effective-mobile/internal/repo/repoerr"
	"github.com/khostya/effective-mobile/internal/repo/transactor"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"testing"
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
	s.T().Parallel()
	s.transactor = transactor.NewTransactionManager(db.GetPool())
	s.songRepo = repo.NewSongRepo(s.transactor)
	s.groupRepo = repo.NewGroupRepo(s.transactor)
	s.ctx = context.Background()
}

func (s *SongsTestSuite) TestCreate() {
	_ = s.create()
}

func (s *SongsTestSuite) TestGetByID() {
	song := s.create()
	song.Group = nil

	actual, err := s.songRepo.GetByID(s.ctx, song.ID)

	require.NoError(s.T(), err)
	require.EqualExportedValues(s.T(), song, actual)
}

func (s *SongsTestSuite) TestDelete() {
	song := s.create()

	err := s.songRepo.Delete(s.ctx, song.ID)
	require.NoError(s.T(), err)

	_, err = s.songRepo.GetByID(s.ctx, song.ID)
	require.ErrorIs(s.T(), err, repoerr.ErrNotFound)
}

func (s *SongsTestSuite) TestUpdate() {
	song := s.create()
	song.Group = nil

	randomString := uuid.New().String()
	err := s.songRepo.Update(s.ctx, dto.UpdateSongParam{
		ID:   song.ID,
		Song: &randomString,
		Text: &randomString,
		Link: &randomString,
	})
	require.NoError(s.T(), err)

	song.Song = randomString
	song.Link = randomString
	song.Text = randomString

	actual, err := s.songRepo.GetByID(s.ctx, song.ID)
	require.NoError(s.T(), err)
	require.EqualExportedValues(s.T(), song, actual)

	song = actual

	var empty string
	err = s.songRepo.Update(s.ctx, dto.UpdateSongParam{
		ID:   song.ID,
		Song: &empty,
		Text: &empty,
		Link: &empty,
	})
	require.NoError(s.T(), err)

	song.Link = empty
	song.Text = empty

	actual, err = s.songRepo.GetByID(s.ctx, song.ID)
	require.NoError(s.T(), err)
	require.EqualExportedValues(s.T(), song, actual)
}

func (s *SongsTestSuite) create() domain.Song {
	group := s.createGroup()

	song := NewSong(group)

	err := s.songRepo.Create(s.ctx, song)
	require.NoError(s.T(), err)

	return song
}

func (s *SongsTestSuite) createGroup() domain.Group {
	group := NewGroup()

	err := s.groupRepo.CreateOnConflictDoNothing(s.ctx, group)
	require.NoError(s.T(), err)

	return group
}
