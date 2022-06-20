package service

import (
	"context"
	"firebase.google.com/go/auth"
	"github.com/stretchr/testify/mock"
)

type FBClientMock struct {
	mock.Mock
}

func NewFBClientMock() *FBClientMock {
	return &FBClientMock{}
}

func (m *FBClientMock) GetUserByEmail(ctx context.Context, email string) (*auth.UserRecord, error) {
	args := m.Called(ctx, email)
	return args.Get(0).(*auth.UserRecord), args.Error(1)
}

func (m *FBClientMock) DeleteUser(ctx context.Context, userID string) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}
