//go:generate mockgen -source ./mocks/api.go -destination=./mocks/api_mock.go -package=mock_api
package api

import (
	"context"
	"errors"
	"github.com/khostya/effective-mobile/internal/dto"
	"github.com/khostya/effective-mobile/pkg/api"
	"net/http"
	"time"
)

type SongInfo struct {
	server api.ClientInterface
}

func NewSongInfo(server api.ClientInterface) *SongInfo {
	return &SongInfo{
		server: server,
	}
}

func (s SongInfo) GetInfo(ctx context.Context, param dto.GetSongInfo) (*dto.SongDetail, error) {
	resp, err := s.server.GetInfo(ctx, &api.GetInfoParams{
		Song:  param.Song,
		Group: param.Group,
	})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	info, err := api.ParseGetInfoResponse(resp)
	if err != nil {
		return nil, err
	}
	if info.StatusCode() != http.StatusOK {
		return nil, errors.New("internal error")
	}

	releaseDate, err := time.Parse(time.DateOnly, info.JSON200.ReleaseDate)
	if err != nil {
		return nil, errors.New("invalid release date")
	}

	return &dto.SongDetail{
		ReleaseDate: releaseDate,
		Link:        info.JSON200.Link,
		Text:        info.JSON200.Text,
	}, nil
}
