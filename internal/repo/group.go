//go:generate mockgen -source ./mocks/group.go -destination=./mocks/group_mock.go -package=mock_repository
package repo

import (
	"context"
	sq "github.com/Masterminds/squirrel"
	"github.com/khostya/effective-mobile/internal/domain"
	"github.com/khostya/effective-mobile/internal/repo/exec"
	"github.com/khostya/effective-mobile/internal/repo/repoerr"
	"github.com/khostya/effective-mobile/internal/repo/schema"
	"github.com/khostya/effective-mobile/internal/repo/transactor"
)

const (
	groupsTable = "effective.groups"
)

type (
	Group struct {
		queryEngineProvider transactor.QueryEngineProvider
	}
)

func (g Group) CreateOnConflictDoNothing(ctx context.Context, group domain.Group) error {
	db := g.queryEngineProvider.GetQueryEngine(ctx)

	record := schema.NewGroup(group)
	query := sq.Insert(groupsTable).
		Columns(record.InsertColumns()...).
		Values(record.InsertValues()...).
		PlaceholderFormat(sq.Dollar).
		Suffix("ON CONFLICT DO NOTHING")

	err := exec.Insert(ctx, query, db)
	if err != nil && exec.IsDuplicateKeyError(err) {
		return repoerr.ErrDuplicate
	}

	return err
}

func (g Group) GetByID(ctx context.Context, title string) (*domain.Group, error) {
	db := g.queryEngineProvider.GetQueryEngine(ctx)

	query := sq.Select(schema.Group{}.SelectColumns()...).
		From(groupsTable).
		PlaceholderFormat(sq.Dollar).
		Where("groups.title = ?", title)

	group, err := exec.ScanOne[schema.Group](ctx, query, db)
	if err != nil {
		return &domain.Group{}, err
	}

	return schema.NewDomainGroup(&group), nil
}

func NewGroupRepo(provider transactor.QueryEngineProvider) Group {
	return Group{provider}
}
