package exec

import (
	"context"
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/khostya/effective-mobile/internal/repo/repoerr"
	"github.com/khostya/effective-mobile/internal/repo/transactor"
)

func InsertWithReturningID(ctx context.Context, query sq.InsertBuilder, db transactor.QueryEngine) (uint, error) {
	row, err := InsertWithRow(ctx, query, db)
	if err != nil {
		return 0, err
	}

	var id uint
	err = row.Scan(&id)

	if err != nil && IsDuplicateKeyError(err) {
		return 0, repoerr.ErrDuplicate
	}

	return id, err
}

func InsertWithRow(ctx context.Context, query sq.InsertBuilder, db transactor.QueryEngine) (pgx.Row, error) {
	rawQuery, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	row := db.QueryRow(ctx, rawQuery, args...)
	return row, nil
}

func Insert(ctx context.Context, query sq.InsertBuilder, db transactor.QueryEngine) error {
	rawQuery, args, err := query.ToSql()
	if err != nil {
		return err
	}

	_, err = db.Exec(ctx, rawQuery, args...)
	if err == nil {
		return nil
	}

	if IsDuplicateKeyError(err) {
		return repoerr.ErrDuplicate
	}
	return err
}
