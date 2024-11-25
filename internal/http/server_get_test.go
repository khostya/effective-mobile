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
	"time"
)

func TestServer_Get(t *testing.T) {
	t.Parallel()

	type (
		input struct {
			page           string
			size           string
			group          string
			song           string
			releaseDateLte string
			releaseDateGte string
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
				group:          gofakeit.UUID(),
				song:           gofakeit.UUID(),
				page:           "1",
				size:           "3",
				releaseDateGte: time.Now().Format(time.DateOnly),
				releaseDateLte: time.Now().Format(time.DateOnly),
			},
			code: http.StatusOK,
			mockFn: func(m mocks) {
				m.song.EXPECT().
					Get(gomock.Any(), gomock.Any()).
					Return(nil, nil).
					Times(1)
			},
		},
		{
			name: "error page release_date_lte",
			input: input{
				group:          gofakeit.UUID(),
				song:           gofakeit.UUID(),
				page:           "1",
				size:           "3",
				releaseDateGte: time.Now().Format(time.DateOnly),
				releaseDateLte: "Gweg",
			},
			code: http.StatusBadRequest,
			mockFn: func(m mocks) {
			},
		},
		{
			name: "error page release_date_gte",
			input: input{
				group:          gofakeit.UUID(),
				song:           gofakeit.UUID(),
				page:           "1",
				size:           "3",
				releaseDateGte: "GFg",
				releaseDateLte: time.Now().Format(time.DateOnly),
			},
			code: http.StatusBadRequest,
			mockFn: func(m mocks) {
			},
		},
		{
			name: "error page invalid",
			input: input{
				group: gofakeit.UUID(),
				song:  gofakeit.UUID(),
				page:  gofakeit.UUID(),
				size:  "3",
			},
			code: http.StatusBadRequest,
			mockFn: func(m mocks) {
			},
		},
		{
			name: "error internal server error",
			input: input{
				group: gofakeit.UUID(),
				song:  gofakeit.UUID(),
				page:  "2",
				size:  "3",
			},
			code: http.StatusInternalServerError,
			mockFn: func(m mocks) {
				m.song.EXPECT().
					Get(gomock.Any(), gomock.Any()).
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

			q := r.URL.Query()
			q.Set(pageParam, tt.input.page)
			q.Set(sizeParam, tt.input.size)
			q.Set(songParam, tt.input.song)
			q.Set(groupParam, tt.input.group)
			q.Set(releaseDateGteParam, tt.input.releaseDateGte)
			q.Set(releaseDateLteParam, tt.input.releaseDateLte)

			r.URL.RawQuery = q.Encode()

			r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, chiCtx))

			w := httptest.NewRecorder()
			server.Get(w, r)

			require.Equal(t, tt.code, w.Code)
		})
	}
}
