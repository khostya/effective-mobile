package schema

import (
	"github.com/khostya/effective-mobile/internal/domain"
)

type (
	Group struct {
		Title string `db:"title"`
	}
)

func (g Group) InsertValues() []any {
	return []any{g.Title}
}

func (g Group) InsertColumns() []string {
	return []string{"title"}
}

func (g Group) SelectColumns() []string {
	return []string{"title"}
}

func NewGroup(song domain.Group) Group {
	return Group{
		Title: song.Title,
	}
}

func NewDomainGroup(song *Group) *domain.Group {
	if song == nil {
		return nil
	}
	return &domain.Group{
		Title: song.Title,
	}
}

func NewDomainGroups(songs []*Group) []*domain.Group {
	if songs == nil {
		return nil
	}

	var res []*domain.Group
	for _, song := range songs {
		res = append(res, NewDomainGroup(song))
	}

	return res
}
