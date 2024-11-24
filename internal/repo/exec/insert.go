package exec

import (
	"context"
	sq "github.com/Masterminds/squirrel"
	"github.com/khostya/effective-mobile/internal/repo/repoerr"
	"github.com/khostya/effective-mobile/internal/repo/transactor"
)

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
