package domain

import (
	"github.com/google/uuid"
	"github.com/khostya/effective-mobile/internal/dto"
	"strings"
	"time"
)

type (
	Text string

	Song struct {
		ID          uuid.UUID `json:"id"`
		Song        string    `json:"song"`
		Group       *Group    `json:"group"`
		Link        string    `json:"link"`
		Text        Text      `json:"text"`
		ReleaseDate time.Time `json:"release_date"`
	}
)

func (t Text) GetVerse(page dto.Page) ([]string, error) {
	verses := strings.Split(string(t), "\n\n")
	left, err := page.Offset()
	if err != nil {
		return nil, err
	}
	if left >= uint(len(verses)) {
		return nil, ErrOutOfRange
	}

	right := int(left + page.Size)
	return verses[left:min(len(verses), right)], nil
}
