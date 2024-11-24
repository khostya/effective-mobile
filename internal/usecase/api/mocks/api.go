// DONT EDIT: Auto generated

package mock_api

import (
	"context"

	"github.com/khostya/effective-mobile/internal/dto"
)

// songApi ...
type songApi interface {
	GetInfo(ctx context.Context, param dto.GetSongInfo) (*dto.SongDetail, error)
}
