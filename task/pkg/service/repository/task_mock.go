package repository

import (
	"context"
	"github.com/stretchr/testify/mock"
)

type FSTaskMock struct {
	mock.Mock
}

func NewMockRepo() *FSTaskMock {
	return &FSTaskMock{}
}

func (m *FSTaskMock) Get(ctx context.Context, userID, taskID string) (Task, error) {
	args := m.Called(ctx, userID, taskID)
	return args.Get(0).(Task), args.Error(1)
}

func (m *FSTaskMock) Create(ctx context.Context, in Task) (Task, error) {
	args := m.Called(ctx, in)
	return args.Get(0).(Task), args.Error(1)
}

func (m *FSTaskMock) Update(ctx context.Context, newTask Task, userID, taskID string) (Task, error) {
	args := m.Called(ctx, newTask, userID, taskID)
	return args.Get(0).(Task), args.Error(1)
}

func (m *FSTaskMock) Delete(ctx context.Context, userID, taskID string) error {
	args := m.Called(ctx, userID, taskID)
	return args.Error(0)
}

func (m *FSTaskMock) GetLastN(ctx context.Context, userID string, n int32) (tasks []Task, err error) {
	args := m.Called(ctx, userID, n)
	return args.Get(0).([]Task), args.Error(1)
}

func (m *FSTaskMock) GetExpired(ctx context.Context, userID string) (expiredTasks []Task, err error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]Task), args.Error(1)
}

func (m *FSTaskMock) SearchForExpiringTasks(ctx context.Context) (map[string][]Task, error) {
	args := m.Called(ctx)
	return args.Get(0).(map[string][]Task), args.Error(1)
}
