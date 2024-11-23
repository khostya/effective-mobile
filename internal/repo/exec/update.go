package exec

import (
	"context"
	sq "github.com/Masterminds/squirrel"
	"github.com/khostya/effective-mobile/internal/repo/repoerr"
	"github.com/khostya/effective-mobile/internal/repo/transactor"
)

func Update(ctx context.Context, query sq.UpdateBuilder, db transactor.QueryEngine) error {
	rawQuery, args, err := query.ToSql()
	if err != nil {
		return err
	}

	tag, err := db.Exec(ctx, rawQuery, args...)
	if err == nil && tag.RowsAffected() == 0 {
		return repoerr.ErrNotFound
	}
	return err
}
