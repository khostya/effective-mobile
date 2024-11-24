package http

import (
	"errors"
	mock_usecase "github.com/khostya/effective-mobile/internal/usecase/mocks"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"net/http"
	"net/http/httptest"
	"testing"
)

var (
	errMock = errors.New("mock error")
)

type (
	mocks struct {
		song *mock_usecase.MocksongUseCase
	}
)

func newMocks(t *testing.T) mocks {
	ctrl := gomock.NewController(t)
	return mocks{
		song: mock_usecase.NewMocksongUseCase(ctrl),
	}
}

func TestLogging(t *testing.T) {
	log := logging(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {

	}))

	r := httptest.NewRequest("GET", "http://localhost/", nil)
	w := httptest.NewRecorder()
	log.ServeHTTP(w, r)
}

func TestParsePage(t *testing.T) {
	r := httptest.NewRequest("GET", "http://localhost/", nil)

	t.Run("ok", func(t *testing.T) {
		q := r.URL.Query()
		q.Set(pageParam, "1")
		q.Set(sizeParam, "2")
		r.URL.RawQuery = q.Encode()

		page, err := parsePage(r)
		require.NoError(t, err)
		require.Equal(t, uint(1), page.Page)
		require.Equal(t, uint(2), page.Size)
	})

	t.Run("error page invalid", func(t *testing.T) {
		q := r.URL.Query()
		q.Set(pageParam, "sdfsf")
		q.Set(sizeParam, "1")
		r.URL.RawQuery = q.Encode()

		_, err := parsePage(r)
		require.Error(t, err)
	})

	t.Run("error size invalid", func(t *testing.T) {
		q := r.URL.Query()
		q.Set(pageParam, "1")
		q.Set(sizeParam, "fdsfas")
		r.URL.RawQuery = q.Encode()

		_, err := parsePage(r)
		require.Error(t, err)
	})
}
