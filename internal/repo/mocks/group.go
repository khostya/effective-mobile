// DONT EDIT: Auto generated

package mock_repository

import (
	"context"

	"github.com/khostya/effective-mobile/internal/domain"
)

// groupStorage ...
type groupStorage interface {
	CreateOnConflictDoNothing(ctx context.Context, group domain.Group) error
	GetByID(ctx context.Context, title string) (*domain.Group, error)
}
