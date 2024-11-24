package http

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/go-chi/chi/v5"
	"github.com/khostya/effective-mobile/internal/http/api"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestServer_Put(t *testing.T) {
	t.Parallel()

	type (
		input struct {
			body any
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
				body: api.PutIdJSONBody{},
				id:   gofakeit.UUID(),
			},
			code: http.StatusOK,
			mockFn: func(m mocks) {
				m.song.EXPECT().
					Update(gomock.Any(), gomock.Any()).
					Return(nil).
					Times(1)
			},
		},
		{
			name: "error internal server error",
			input: input{
				body: api.PutIdJSONBody{},
				id:   gofakeit.UUID(),
			},
			code: http.StatusInternalServerError,
			mockFn: func(m mocks) {
				m.song.EXPECT().
					Update(gomock.Any(), gomock.Any()).
					Return(errMock).
					Times(1)
			},
		},
		{
			name: "error invalid id",
			input: input{
				body: api.PutIdJSONBody{},
				id:   "Grgegre",
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
			chiCtx.URLParams.Add(idParam, tt.input.id)
			r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, chiCtx))

			w := httptest.NewRecorder()
			server.Put(w, r)

			require.Equal(t, tt.code, w.Code)
		})
	}
}
