package usecase

import (
	"context"
	"errors"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	"github.com/khostya/effective-mobile/internal/domain"
	"github.com/khostya/effective-mobile/internal/dto"
	mock_repository "github.com/khostya/effective-mobile/internal/repo/mocks"
	mock_transactor "github.com/khostya/effective-mobile/internal/repo/transactor/mocks"
	mock_api "github.com/khostya/effective-mobile/internal/usecase/api/mocks"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"testing"
)

var (
	errMock = errors.New("mock error")
)

type mocks struct {
	mockGroupRepository *mock_repository.MockgroupRepo
	mockSongRepository  *mock_repository.MocksongRepo
	mockTransactor      *mock_transactor.MockTransactor
	mockAPI             *mock_api.MocksongApi
}

func newMocks(t *testing.T) mocks {
	ctrl := gomock.NewController(t)

	return mocks{
		mockTransactor:      mock_transactor.NewMockTransactor(ctrl),
		mockSongRepository:  mock_repository.NewMocksongRepo(ctrl),
		mockGroupRepository: mock_repository.NewMockgroupRepo(ctrl),
		mockAPI:             mock_api.NewMocksongApi(ctrl),
	}
}

func TestSong_DeleteByID(t *testing.T) {
	t.Parallel()

	type test struct {
		name   string
		input  uuid.UUID
		err    error
		mockFn func(m mocks)
	}
	var ctx = context.Background()
	tests := []test{
		{
			name:  "ok",
			input: uuid.New(),
			err:   nil,
			mockFn: func(m mocks) {
				m.mockSongRepository.EXPECT().
					Delete(gomock.Any(), gomock.Any()).
					Return(nil).
					Times(1)
			},
		},
		{
			name:  "error",
			input: uuid.New(),
			err:   errMock,
			mockFn: func(m mocks) {
				m.mockSongRepository.EXPECT().
					Delete(gomock.Any(), gomock.Any()).
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
			orderService := NewSongUseCase(SongDeps{
				SongRepo:  mocks.mockSongRepository,
				GroupRepo: mocks.mockGroupRepository,
				Tm:        mocks.mockTransactor,
			})

			err := orderService.DeleteByID(ctx, tt.input)

			require.ErrorIs(t, err, tt.err)
		})
	}
}

func TestSong_Update(t *testing.T) {
	t.Parallel()

	type test struct {
		name   string
		input  dto.UpdateSongParam
		err    error
		mockFn func(m mocks)
	}
	var ctx = context.Background()
	tests := []test{
		{
			name:  "ok",
			input: dto.UpdateSongParam{},
			err:   nil,
			mockFn: func(m mocks) {
				m.mockSongRepository.EXPECT().
					Update(gomock.Any(), gomock.Any()).
					Return(nil).
					Times(1)
			},
		},
		{
			name:  "error",
			input: dto.UpdateSongParam{},
			err:   errMock,
			mockFn: func(m mocks) {
				m.mockSongRepository.EXPECT().
					Update(gomock.Any(), gomock.Any()).
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
			songUseCase := NewSongUseCase(SongDeps{
				SongRepo:  mocks.mockSongRepository,
				GroupRepo: mocks.mockGroupRepository,
				Tm:        mocks.mockTransactor,
			})

			err := songUseCase.Update(ctx, tt.input)

			require.ErrorIs(t, err, tt.err)
		})
	}
}

func TestSong_Get(t *testing.T) {
	t.Parallel()

	type test struct {
		name   string
		input  dto.GetSongsParam
		err    error
		mockFn func(m mocks)
	}
	var ctx = context.Background()
	tests := []test{
		{
			name: "ok",
			input: dto.GetSongsParam{
				Song:  gofakeit.UUID(),
				Group: gofakeit.Animal(),
			},
			err: nil,
			mockFn: func(m mocks) {
				m.mockSongRepository.EXPECT().
					Get(gomock.Any(), gomock.Any()).
					Return(nil, nil).
					Times(1)
			},
		},
		{
			name: "error",
			input: dto.GetSongsParam{
				Song:  gofakeit.Email(),
				Group: gofakeit.Email(),
			},
			err: errMock,
			mockFn: func(m mocks) {
				m.mockSongRepository.EXPECT().
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
			songUseCase := NewSongUseCase(SongDeps{
				SongRepo:  mocks.mockSongRepository,
				GroupRepo: mocks.mockGroupRepository,
				Tm:        mocks.mockTransactor,
			})

			_, err := songUseCase.Get(ctx, tt.input)

			require.ErrorIs(t, err, tt.err)
		})
	}
}

func TestSong_Create(t *testing.T) {
	t.Parallel()

	type test struct {
		name   string
		input  dto.CreateSongParam
		err    error
		mockFn func(m mocks)
	}

	var ctx = context.Background()
	tests := []test{
		{
			name: "ok",
			input: dto.CreateSongParam{
				Song:  gofakeit.UUID(),
				Group: gofakeit.Animal(),
			},
			err: nil,
			mockFn: func(m mocks) {
				m.mockAPI.EXPECT().
					GetInfo(gomock.Any(), gomock.Any()).
					Return(&dto.SongDetail{}, nil).
					Times(1)
				m.mockGroupRepository.EXPECT().
					CreateOnConflictDoNothing(gomock.Any(), gomock.Any()).
					Return(nil).
					Times(1)
				m.mockSongRepository.EXPECT().
					Create(gomock.Any(), gomock.Any()).
					Return(nil).
					Times(1)
				m.mockTransactor.EXPECT().Unwrap(nil).Times(1).Return(nil)
				m.mockTransactor.EXPECT().RunRepeatableRead(gomock.Any(), gomock.Any()).Times(1).
					DoAndReturn(func(ctx context.Context, transaction func(ctx context.Context) error) error {
						return transaction(ctx)
					})
			},
		},
		{
			name: "error get info",
			input: dto.CreateSongParam{
				Song:  gofakeit.UUID(),
				Group: gofakeit.Animal(),
			},
			err: errMock,
			mockFn: func(m mocks) {
				m.mockAPI.EXPECT().
					GetInfo(gomock.Any(), gomock.Any()).
					Return(&dto.SongDetail{}, errMock).
					Times(1)
			},
		},
		{
			name: "error create group on conflict do nothing",
			input: dto.CreateSongParam{
				Song:  gofakeit.UUID(),
				Group: gofakeit.Animal(),
			},
			err: errMock,
			mockFn: func(m mocks) {
				m.mockAPI.EXPECT().
					GetInfo(gomock.Any(), gomock.Any()).
					Return(&dto.SongDetail{}, nil).
					Times(1)
				m.mockGroupRepository.EXPECT().
					CreateOnConflictDoNothing(gomock.Any(), gomock.Any()).
					Return(errMock).
					Times(1)
				m.mockTransactor.EXPECT().Unwrap(errMock).Times(1).Return(errMock)
				m.mockTransactor.EXPECT().RunRepeatableRead(gomock.Any(), gomock.Any()).Times(1).
					DoAndReturn(func(ctx context.Context, transaction func(ctx context.Context) error) error {
						return transaction(ctx)
					})
			},
		},
		{
			name: "error create song",
			input: dto.CreateSongParam{
				Song:  gofakeit.UUID(),
				Group: gofakeit.Animal(),
			},
			err: errMock,
			mockFn: func(m mocks) {
				m.mockAPI.EXPECT().
					GetInfo(gomock.Any(), gomock.Any()).
					Return(&dto.SongDetail{}, nil).
					Times(1)
				m.mockGroupRepository.EXPECT().
					CreateOnConflictDoNothing(gomock.Any(), gomock.Any()).
					Return(nil).
					Times(1)
				m.mockSongRepository.EXPECT().
					Create(gomock.Any(), gomock.Any()).
					Return(errMock).
					Times(1)
				m.mockTransactor.EXPECT().Unwrap(errMock).Times(1).Return(errMock)
				m.mockTransactor.EXPECT().RunRepeatableRead(gomock.Any(), gomock.Any()).Times(1).
					DoAndReturn(func(ctx context.Context, transaction func(ctx context.Context) error) error {
						return transaction(ctx)
					})
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mocks := newMocks(t)
			tt.mockFn(mocks)

			songUseCase := NewSongUseCase(SongDeps{
				SongRepo:  mocks.mockSongRepository,
				GroupRepo: mocks.mockGroupRepository,
				Tm:        mocks.mockTransactor,
				InfoSong:  mocks.mockAPI,
			})

			err := songUseCase.Create(ctx, tt.input)

			require.ErrorIs(t, err, tt.err)
		})
	}
}

func TestSong_GetByVerse(t *testing.T) {
	t.Parallel()

	type test struct {
		name   string
		input  dto.GetByVerseParam
		err    error
		mockFn func(m mocks)
	}

	var ctx = context.Background()
	tests := []test{
		{
			name: "ok",
			input: dto.GetByVerseParam{
				ID: uuid.New(),
				Page: dto.Page{
					Size: 10,
					Page: 1,
				},
			},
			err: nil,
			mockFn: func(m mocks) {
				m.mockSongRepository.EXPECT().
					GetByID(gomock.Any(), gomock.Any()).
					Return(domain.Song{Verses: []string{"ля-ля-ля"}}, nil).
					Times(1)
			},
		},
		{
			name: "error get song by id",
			input: dto.GetByVerseParam{
				ID: uuid.New(),
				Page: dto.Page{
					Size: 10,
					Page: 1,
				},
			},
			err: errMock,
			mockFn: func(m mocks) {
				m.mockSongRepository.EXPECT().
					GetByID(gomock.Any(), gomock.Any()).
					Return(domain.Song{}, errMock).
					Times(1)
			},
		},
		{
			name: "error out of range",
			input: dto.GetByVerseParam{
				ID: uuid.New(),
				Page: dto.Page{
					Size: 10,
					Page: 10,
				},
			},
			err: ErrOutOfRange,
			mockFn: func(m mocks) {
				m.mockSongRepository.EXPECT().
					GetByID(gomock.Any(), gomock.Any()).
					Return(domain.Song{}, nil).
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

			songUseCase := NewSongUseCase(SongDeps{
				SongRepo:  mocks.mockSongRepository,
				GroupRepo: mocks.mockGroupRepository,
				Tm:        mocks.mockTransactor,
				InfoSong:  mocks.mockAPI,
			})

			_, err := songUseCase.GetByVerse(ctx, tt.input)

			require.ErrorIs(t, err, tt.err)
		})
	}
}
