package http

import (
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/khostya/effective-mobile/internal/repo/repoerr"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestServer_DeleteId(t *testing.T) {
	t.Parallel()

	type test struct {
		name   string
		input  string
		code   int
		mockFn func(m mocks)
	}

	tests := []test{
		{
			name:  "ok",
			input: uuid.New().String(),
			code:  http.StatusOK,
			mockFn: func(m mocks) {
				m.song.EXPECT().
					DeleteByID(gomock.Any(), gomock.Any()).
					Return(nil).
					Times(1)
			},
		},
		{
			name:  "error invalid id",
			input: "G13d1",
			code:  http.StatusBadRequest,
			mockFn: func(m mocks) {
			},
		},
		{
			name:  "error not found",
			input: uuid.New().String(),
			code:  http.StatusNotFound,
			mockFn: func(m mocks) {
				m.song.EXPECT().
					DeleteByID(gomock.Any(), gomock.Any()).
					Return(repoerr.ErrNotFound).
					Times(1)
			},
		},
		{
			name:  "error internal server error",
			input: uuid.New().String(),
			code:  http.StatusInternalServerError,
			mockFn: func(m mocks) {
				m.song.EXPECT().
					DeleteByID(gomock.Any(), gomock.Any()).
					Return(errMock).
					Times(1)
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mocks := newMocks(t)
			tt.mockFn(mocks)

			server, err := newServer(UseCases{
				Song: mocks.song,
			})
			require.NoError(t, err)

			r := httptest.NewRequest("DELETE", "http://localhost/", nil)
			chiCtx := chi.NewRouteContext()
			chiCtx.URLParams.Add(idParam, tt.input)

			r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, chiCtx))

			w := httptest.NewRecorder()
			server.DeleteId(w, r)

			require.Equal(t, tt.code, w.Code)
		})
	}
}
