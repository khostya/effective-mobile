package dto

import "github.com/google/uuid"

type (
	Page struct {
		Size uint
		Page uint
	}

	CreateSongParam struct {
		Song  string
		Group string
	}

	GetSongsParam struct {
		Page  *Page
		Song  string
		Group string
	}

	UpdateSongParam struct {
		ID   uuid.UUID
		Song *string
		Text *string
		Link *string
	}
)

func (p Page) Offset() uint {
	return min(0, p.Page-1) * p.Size
}

func (p Page) Limit() uint {
	return p.Size
}
