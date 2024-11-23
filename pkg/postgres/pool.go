package postgres

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
	"os"
)

type Pool = pgxpool.Pool

func NewPool(ctx context.Context, url string) (*Pool, error) {
	pool, err := pgxpool.New(ctx, url)
	if err != nil {
		return nil, errors.Wrap(err, "Unable to create connection pool")
	}

	return pool, nil
}

func NewPoolFromEnv(ctx context.Context, key string) (*Pool, error) {
	url := os.Getenv(key)
	if url == "" {
		return nil, errors.New(fmt.Sprintf("Unable to parse %s", key))
	}

	return NewPool(ctx, url)
}
