package http

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestServer(t *testing.T) {
	_, err := newServer(UseCases{})
	require.NoError(t, err)
}
