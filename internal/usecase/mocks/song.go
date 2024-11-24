// DONT EDIT: Auto generated

package mock_usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/khostya/effective-mobile/internal/domain"
	"github.com/khostya/effective-mobile/internal/dto"
)

// songUseCase ...
type songUseCase interface {
	DeleteByID(ctx context.Context, id uuid.UUID) error
	Create(ctx context.Context, param dto.CreateSongParam) error
	Get(ctx context.Context, param dto.GetSongsParam) ([]domain.Song, error)
	GetByVerse(ctx context.Context, param dto.GetByVerseParam) ([]string, error)
	Update(ctx context.Context, param dto.UpdateSongParam) error
}
