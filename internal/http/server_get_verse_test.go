package http

import (
	"context"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestServer_GetByVerse(t *testing.T) {
	t.Parallel()

	type (
		input struct {
			page string
			size string
			id   string
		}
		test struct {
			name   string
			input  input
			code   int
			mockFn func(m mocks)
		}
	)

	tests := []test{
		{
			name: "ok",
			input: input{
				id:   gofakeit.UUID(),
				page: "1",
				size: "3",
			},
			code: http.StatusOK,
			mockFn: func(m mocks) {
				m.song.EXPECT().
					GetByVerse(gomock.Any(), gomock.Any()).
					Return(nil, nil).
					Times(1)
			},
		},
		{
			name: "error id invalid",
			input: input{
				id:   "dfs",
				page: "13",
				size: "3",
			},
			code: http.StatusBadRequest,
			mockFn: func(m mocks) {
			},
		},
		{
			name: "error invalid page",
			input: input{
				page: gofakeit.UUID(),
				size: gofakeit.UUID(),
				id:   gofakeit.UUID(),
			},
			code: http.StatusBadRequest,
			mockFn: func(m mocks) {
			},
		},
		{
			name: "error internal server error",
			input: input{
				page: "4",
				size: "1",
				id:   gofakeit.UUID(),
			},
			code: http.StatusInternalServerError,
			mockFn: func(m mocks) {
				m.song.EXPECT().
					GetByVerse(gomock.Any(), gomock.Any()).
					Return(nil, errMock).
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

			r := httptest.NewRequest("GET", "http://localhost/", nil)
			chiCtx := chi.NewRouteContext()
			chiCtx.URLParams.Add(idParam, tt.input.id)

			q := r.URL.Query()
			q.Set(pageParam, tt.input.page)
			q.Set(sizeParam, tt.input.size)
			r.URL.RawQuery = q.Encode()

			r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, chiCtx))

			w := httptest.NewRecorder()
			server.GetVerseId(w, r)

			require.Equal(t, tt.code, w.Code)
		})
	}
}
