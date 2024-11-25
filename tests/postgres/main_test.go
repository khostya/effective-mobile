//go:build integration

package postgres

import (
	"context"
	"github.com/khostya/effective-mobile/tests/postgres/postgresql"
	"os"
	"testing"
)

var (
	db *postgresql.DBPool
)

const (
	songsTable  = "effective.songs"
	groupsTable = "effective.groups"
)

func TestMain(m *testing.M) {
	db = postgresql.NewFromEnv()

	code := m.Run()
	truncate()
	db.Close()

	os.Exit(code)
}

func truncate() {
	db.TruncateTable(context.Background(), groupsTable, songsTable)
}
