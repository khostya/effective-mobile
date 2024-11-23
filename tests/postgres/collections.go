//go:build integration

package postgres

import (
	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	"github.com/khostya/effective-mobile/internal/domain"
	"time"
)

func NewSong(group domain.Group) domain.Song {
	return domain.Song{
		ID:          uuid.New(),
		Song:        gofakeit.UUID(),
		Group:       &group,
		Link:        gofakeit.URL(),
		Text:        gofakeit.UUID(),
		ReleaseDate: time.Now(),
	}
}

func NewGroup() domain.Group {
	return domain.Group{
		Title: gofakeit.UUID(),
	}
}
