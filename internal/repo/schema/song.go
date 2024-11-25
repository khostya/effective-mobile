package schema

import (
	"database/sql"
	"github.com/google/uuid"
	"github.com/khostya/effective-mobile/internal/domain"
	"github.com/khostya/effective-mobile/internal/dto"
	"time"
)

type (
	Song struct {
		ID          uuid.UUID        `db:"id"`
		Song        string           `db:"song"`
		GroupTitle  string           `db:"group_title"`
		Link        sql.Null[string] `db:"link"`
		Text        sql.Null[string] `db:"text"`
		ReleaseDate time.Time        `db:"release_date"`
	}
)

func (f Song) SelectColumns() []string {
	return []string{"id", "song", "group_title", "link", "text", "release_date"}
}

func (f Song) InsertValues() []any {
	return []any{f.ID, f.Song, f.GroupTitle, f.Link, f.Text, f.ReleaseDate}
}

func (f Song) InsertColumns() []string {
	return []string{"id", "song", "group_title", "link", "text", "release_date"}
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
		Text:        nullIfDefault(string(song.Text)),
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
		Text:        domain.Text(song.Text.V),
		ReleaseDate: song.ReleaseDate,
	}
}

func NewDomainSongs(songs []Song) []domain.Song {
	var res []domain.Song = make([]domain.Song, 0)

	for _, song := range songs {
		res = append(res, NewDomainSong(song))
	}

	return res
}

func NewSongUpdate(param dto.UpdateSongParam) Song {
	var text = sql.Null[string]{Valid: true}
	var link = sql.Null[string]{Valid: true}
	var song string
	if param.Song != nil {
		song = *param.Song
	}

	if param.Link != nil {
		link.V = *param.Link
	}
	if param.Text != nil {
		text.V = *param.Text
	}

	return Song{
		ID:   param.ID,
		Song: song,
		Link: link,
		Text: text,
	}
}

func (f Song) UpdateColumns() []string {
	res := []string{"link", "text"}
	if f.Song != "" {
		res = append(res, "song")
	}
	return res
}

func (f Song) UpdateValues() []any {
	res := []any{f.Link, f.Text}
	if f.Song != "" {
		res = append(res, f.Song)
	}
	return res
}
