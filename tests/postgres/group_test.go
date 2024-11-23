//go:build integration

package postgres

import (
	"context"
	"github.com/khostya/effective-mobile/internal/repo"
	"github.com/khostya/effective-mobile/internal/repo/transactor"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"testing"
)

type GroupTestSuite struct {
	suite.Suite
	ctx        context.Context
	groupRepo  repo.Group
	transactor *transactor.TransactionManager
}

func TestGroup(t *testing.T) {
	suite.Run(t, new(GroupTestSuite))
}

func (s *GroupTestSuite) SetupSuite() {
	s.T().Parallel()
	s.transactor = transactor.NewTransactionManager(db.GetPool())
	s.groupRepo = repo.NewGroupRepo(s.transactor)
	s.ctx = context.Background()
}

func (s *GroupTestSuite) TestCreateOnConflictDoNothing() {
	group := NewGroup()

	err := s.groupRepo.CreateOnConflictDoNothing(s.ctx, group)
	require.NoError(s.T(), err)

	err = s.groupRepo.CreateOnConflictDoNothing(s.ctx, group)
	require.NoError(s.T(), err)
}

func (s *GroupTestSuite) TestGetByID() {
	group := NewGroup()

	err := s.groupRepo.CreateOnConflictDoNothing(s.ctx, group)
	require.NoError(s.T(), err)

	actual, err := s.groupRepo.GetByID(s.ctx, group.Title)
	require.NoError(s.T(), err)
	require.EqualExportedValues(s.T(), group, *actual)
}
