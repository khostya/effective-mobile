// DONT EDIT: Auto generated

package mock_repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/khostya/effective-mobile/internal/domain"
	"github.com/khostya/effective-mobile/internal/dto"
)

// songRepo ...
type songRepo interface {
	Create(ctx context.Context, song domain.Song) error
	Get(ctx context.Context, param dto.GetSongsParam) ([]domain.Song, error)
	GetByID(ctx context.Context, id uuid.UUID) (domain.Song, error)
	Update(ctx context.Context, param dto.UpdateSongParam) error
	Delete(ctx context.Context, id uuid.UUID) error
}
