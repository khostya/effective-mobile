package repo

import (
	"github.com/khostya/effective-mobile/internal/repo/transactor"
)

type Repositories struct {
	Song  Song
	Group Group
}

func NewRepositories(provider transactor.QueryEngineProvider) Repositories {
	return Repositories{
		Song:  NewSongRepo(provider),
		Group: NewGroupRepo(provider),
	}
}
