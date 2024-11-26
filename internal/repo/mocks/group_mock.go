// Code generated by MockGen. DO NOT EDIT.
// Source: ./mocks/group.go
//
// Generated by this command:
//
//	mockgen -source ./mocks/group.go -destination=./mocks/group_mock.go -package=mock_repository
//

// Package mock_repository is a generated GoMock package.
package mock_repository

import (
	context "context"
	reflect "reflect"

	domain "github.com/khostya/effective-mobile/internal/domain"
	gomock "go.uber.org/mock/gomock"
)

// MockgroupRepo is a mock of groupRepo interface.
type MockgroupRepo struct {
	ctrl     *gomock.Controller
	recorder *MockgroupRepoMockRecorder
	isgomock struct{}
}

// MockgroupRepoMockRecorder is the mock recorder for MockgroupRepo.
type MockgroupRepoMockRecorder struct {
	mock *MockgroupRepo
}

// NewMockgroupRepo creates a new mock instance.
func NewMockgroupRepo(ctrl *gomock.Controller) *MockgroupRepo {
	mock := &MockgroupRepo{ctrl: ctrl}
	mock.recorder = &MockgroupRepoMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockgroupRepo) EXPECT() *MockgroupRepoMockRecorder {
	return m.recorder
}

// CreateOnConflictDoNothing mocks base method.
func (m *MockgroupRepo) CreateOnConflictDoNothing(ctx context.Context, group domain.Group) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateOnConflictDoNothing", ctx, group)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateOnConflictDoNothing indicates an expected call of CreateOnConflictDoNothing.
func (mr *MockgroupRepoMockRecorder) CreateOnConflictDoNothing(ctx, group any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateOnConflictDoNothing", reflect.TypeOf((*MockgroupRepo)(nil).CreateOnConflictDoNothing), ctx, group)
}

// GetByID mocks base method.
func (m *MockgroupRepo) GetByID(ctx context.Context, title string) (*domain.Group, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByID", ctx, title)
	ret0, _ := ret[0].(*domain.Group)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByID indicates an expected call of GetByID.
func (mr *MockgroupRepoMockRecorder) GetByID(ctx, title any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByID", reflect.TypeOf((*MockgroupRepo)(nil).GetByID), ctx, title)
}
