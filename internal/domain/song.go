package domain

import (
	"github.com/google/uuid"
	"time"
)

type (
	Song struct {
		ID          uuid.UUID `json:"id"`
		Song        string    `json:"song"`
		Group       *Group    `json:"group"`
		Link        string    `json:"link"`
		Text        string    `json:"text"`
		ReleaseDate time.Time `json:"release_date"`
	}
)
