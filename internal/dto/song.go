package dto

import (
	"errors"
	"github.com/google/uuid"
	"time"
)

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

	GetSongInfo struct {
		Song, Group string
	}

	SongDetail struct {
		ReleaseDate time.Time
		Link        string
		Text        string
	}

	GetByVerseParam struct {
		ID   uuid.UUID
		Page Page
	}
)

func (p Page) Offset() (uint, error) {
	if p.Page <= 0 {
		return 0, errors.New("page out of range")
	}
	return (p.Page - 1) * p.Size, nil
}

func (p Page) Limit() uint {
	return p.Size
}
