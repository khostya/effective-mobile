package domain

import (
	"github.com/google/uuid"
	"github.com/khostya/effective-mobile/internal/dto"
	"time"
)

type (
	Text []string

	Song struct {
		ID          uuid.UUID `json:"id"`
		Song        string    `json:"song"`
		Group       *Group    `json:"group"`
		Link        string    `json:"link"`
		Verses      Text      `json:"verses"`
		ReleaseDate time.Time `json:"release_date"`
	}
)

func (t Text) GetVerse(page dto.Page) ([]string, error) {
	left, err := page.Offset()
	if err != nil {
		return nil, err
	}
	if left >= uint(len(t)) {
		return nil, ErrOutOfRange
	}

	right := int(left + page.Size)
	return t[left:min(len(t), right)], nil
}
