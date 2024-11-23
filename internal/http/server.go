package http

import (
	"context"
	"github.com/google/uuid"
	"github.com/khostya/effective-mobile/internal/domain"
	"github.com/khostya/effective-mobile/internal/dto"
	"github.com/khostya/effective-mobile/internal/usecase"
	"github.com/khostya/effective-mobile/pkg/validator"
)

var (
	_ songUseCase = (*usecase.Song)(nil)
)

type (
	songUseCase interface {
		DeleteByID(ctx context.Context, id uuid.UUID) error
		Create(ctx context.Context, param dto.CreateSongParam) error
		Get(ctx context.Context, param dto.GetSongsParam) ([]domain.Song, error)
		GetByVerse(ctx context.Context, id uuid.UUID, page dto.Page) ([]string, error)
		Update(ctx context.Context, param dto.UpdateSongParam) error
	}

	UseCases struct {
		Song songUseCase
	}

	server struct {
		useCases  UseCases
		validator *validator.Validator
	}
)

func newServer(useCases UseCases) (*server, error) {
	validator, err := validator.NewValidate()
	if err != nil {
		return nil, err
	}

	return &server{
		useCases:  useCases,
		validator: validator,
	}, nil
}
