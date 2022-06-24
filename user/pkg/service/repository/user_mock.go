package repository

import (
	"context"
	"github.com/stretchr/testify/mock"
)

type FSUserMock struct {
	mock.Mock
}

func (m *FSUserMock) Get(ctx context.Context, UserId string) (User, error) {
	args := m.Called(ctx, UserId)
	return args.Get(0).(User), args.Error(1)
}

func (m *FSUserMock) Update(ctx context.Context, UserId string, user User) (User, error) {
	args := m.Called(ctx, UserId, user)
	return args.Get(0).(User), args.Error(1)
}

func (m *FSUserMock) Delete(ctx context.Context, UserId string) error {
	args := m.Called(ctx, UserId)
	return args.Error(0)
}

func NewMockRepo() *FSUserMock {
	return &FSUserMock{}
}
