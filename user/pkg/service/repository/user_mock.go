package repository

import (
	"context"
	"github.com/stretchr/testify/mock"
)

type FSUserMock struct {
	mock.Mock
}

func (m *FSUserMock) Get(ctx context.Context, userID string) (User, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(User), args.Error(1)
}

func (m *FSUserMock) Update(ctx context.Context, userID string, user User) (User, error) {
	args := m.Called(ctx, userID, user)
	return args.Get(0).(User), args.Error(1)
}

func (m *FSUserMock) Delete(ctx context.Context, userID string) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func NewMockRepo() *FSUserMock {
	return &FSUserMock{}
}
