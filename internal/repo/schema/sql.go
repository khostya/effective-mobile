package schema

import "database/sql"

func nullIfDefault[T comparable](v T) sql.Null[T] {
	var def T
	return sql.Null[T]{V: v, Valid: def != v}
}
