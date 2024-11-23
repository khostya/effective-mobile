package exec

import (
	"context"
	sq "github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/khostya/effective-mobile/internal/repo/repoerr"
	"github.com/khostya/effective-mobile/internal/repo/transactor"
)

func ScanOne[T any](ctx context.Context, query sq.SelectBuilder, db transactor.QueryEngine) (T, error) {
	var defaultT T

	records, err := ScanALL[T](ctx, query, db)
	if err != nil {
		return defaultT, err
	}

	if len(records) == 0 {
		return defaultT, repoerr.ErrNotFound
	}

	return records[0], nil
}

func ScanALL[T any](ctx context.Context, query sq.SelectBuilder, db transactor.QueryEngine) ([]T, error) {
	rawQuery, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	var records []T
	if err := pgxscan.Select(ctx, db, &records, rawQuery, args...); err != nil {
		return nil, err
	}

	return records, nil
}
