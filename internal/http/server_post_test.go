package http

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/khostya/effective-mobile/internal/http/api"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestServer_Post(t *testing.T) {
	t.Parallel()

	type (
		test struct {
			name   string
			input  any
			code   int
			mockFn func(m mocks)
		}
	)

	tests := []test{
		{
			name: "ok",
			input: api.PostJSONBody{
				Group: "13",
				Song:  "31",
			},
			code: http.StatusOK,
			mockFn: func(m mocks) {
				m.song.EXPECT().
					Create(gomock.Any(), gomock.Any()).
					Return(nil).
					Times(1)
			},
		},
		{
			name: "error internal server error",
			input: api.PostJSONBody{
				Group: "13",
				Song:  "31",
			},
			code: http.StatusInternalServerError,
			mockFn: func(m mocks) {
				m.song.EXPECT().
					Create(gomock.Any(), gomock.Any()).
					Return(errMock).
					Times(1)
			},
		},
		{
			name:  "error invalid body",
			input: struct{}{},
			code:  http.StatusBadRequest,
			mockFn: func(m mocks) {
			},
		},
		{
			name: "error invalid song",
			input: api.PostJSONBody{
				Group: "1",
				Song:  "",
			},
			code: http.StatusBadRequest,
			mockFn: func(m mocks) {
			},
		},
		{
			name: "error invalid group",
			input: api.PostJSONBody{
				Group: "1",
				Song:  "",
			},
			code: http.StatusBadRequest,
			mockFn: func(m mocks) {
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

			body, err := json.Marshal(tt.input)
			require.NoError(t, err)

			r := httptest.NewRequest("GET", "http://localhost/", bytes.NewReader(body))
			chiCtx := chi.NewRouteContext()

			r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, chiCtx))

			w := httptest.NewRecorder()
			server.PostCreate(w, r)

			require.Equal(t, tt.code, w.Code)
		})
	}
}
