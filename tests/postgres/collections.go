//go:build integration

package postgres

import (
	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	"github.com/khostya/effective-mobile/internal/domain"
	"time"
)

func NewSong(songTitle string, group domain.Group) domain.Song {
	return domain.Song{
		ID:          uuid.New(),
		Song:        songTitle,
		Group:       &group,
		Link:        gofakeit.URL(),
		Text:        domain.Text(gofakeit.UUID()),
		ReleaseDate: time.Now(),
	}
}

func NewGroup(groupTitle string) domain.Group {
	return domain.Group{
		Title: groupTitle,
	}
}
