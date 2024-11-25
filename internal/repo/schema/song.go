package schema

import (
	"database/sql"
	"github.com/google/uuid"
	"github.com/khostya/effective-mobile/internal/domain"
	"github.com/khostya/effective-mobile/internal/dto"
	"github.com/lib/pq"
	"time"
)

type (
	Song struct {
		ID          uuid.UUID        `db:"id"`
		Song        string           `db:"song"`
		GroupTitle  string           `db:"group_title"`
		Link        sql.Null[string] `db:"link"`
		Verses      []string         `db:"verses"`
		ReleaseDate time.Time        `db:"release_date"`
	}
)

func (f Song) SelectColumns() []string {
	return []string{"id", "song", "group_title", "link", "verses", "release_date"}
}

func (f Song) InsertValues() []any {
	return []any{f.ID, f.Song, f.GroupTitle, f.Link, pq.StringArray(f.Verses), f.ReleaseDate}
}

func (f Song) InsertColumns() []string {
	return []string{"id", "song", "group_title", "link", "verses", "release_date"}
}

func NewSong(song domain.Song) Song {
	var title string
	if song.Group != nil {
		title = song.Group.Title
	}

	return Song{
		ID:          song.ID,
		Song:        song.Song,
		GroupTitle:  title,
		Link:        nullIfDefault(song.Link),
		Verses:      song.Verses,
		ReleaseDate: song.ReleaseDate,
	}
}

func NewDomainSong(song Song) domain.Song {
	return domain.Song{
		ID:   song.ID,
		Song: song.Song,
		Group: &domain.Group{
			Title: song.GroupTitle,
		},
		Link:        song.Link.V,
		Verses:      song.Verses,
		ReleaseDate: song.ReleaseDate,
	}
}

func NewDomainSongs(songs []Song) []domain.Song {
	var res = make([]domain.Song, 0)

	for _, song := range songs {
		res = append(res, NewDomainSong(song))
	}

	return res
}

func NewSongUpdate(param dto.UpdateSongParam) Song {
	var link = sql.Null[string]{Valid: true}
	var song string
	if param.Song != nil {
		song = *param.Song
	}

	if param.Link != nil {
		link.V = *param.Link
	}

	return Song{
		ID:     param.ID,
		Song:   song,
		Link:   link,
		Verses: param.Verses,
	}
}

func (f Song) UpdateColumns() []string {
	res := []string{"link"}
	if f.Song != "" {
		res = append(res, "song")
	}
	if f.Verses != nil {
		res = append(res, "verses")
	}
	return res
}

func (f Song) UpdateValues() []any {
	res := []any{f.Link}
	if f.Song != "" {
		res = append(res, f.Song)
	}
	if f.Verses != nil {
		res = append(res, pq.StringArray(f.Verses))
	}
	return res
}
