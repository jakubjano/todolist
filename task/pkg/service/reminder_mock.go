package service

import (
	"context"
	"github.com/stretchr/testify/mock"
)

type ReminderMock struct {
	mock.Mock
}

func NewReminderMock() *ReminderMock {
	return &ReminderMock{}
}

func (m *ReminderMock) RemindUserViaEmail(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}
