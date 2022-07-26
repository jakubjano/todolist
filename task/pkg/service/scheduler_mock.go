package service

import (
	"context"
	"github.com/stretchr/testify/mock"
	cloudscheduler "google.golang.org/api/cloudscheduler/v1beta1"
)

type SchedulerMock struct {
	mock.Mock
}

func NewSchedulerMock() *SchedulerMock {
	return &SchedulerMock{}
}

func (m *SchedulerMock) CreateScheduledJob(ctx context.Context, name, schedule, target, method, description string) (*cloudscheduler.Job, error) {
	args := m.Called(ctx, name, schedule, target, method, description)
	return args.Get(0).(*cloudscheduler.Job), args.Error(1)
}

func (m *SchedulerMock) PatchScheduledJob(ctx context.Context, name, schedule, description string) (*cloudscheduler.Job, error) {
	args := m.Called(ctx, name, schedule, description)
	return args.Get(0).(*cloudscheduler.Job), args.Error(1)
}

func (m *SchedulerMock) ListScheduledJobs(ctx context.Context) ([]*cloudscheduler.Job, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*cloudscheduler.Job), args.Error(1)
}

func (m *SchedulerMock) PauseScheduledJob(ctx context.Context, name string) (*cloudscheduler.Job, error) {
	args := m.Called(ctx, name)
	return args.Get(0).(*cloudscheduler.Job), args.Error(1)
}

func (m *SchedulerMock) ResumeScheduledJob(ctx context.Context, name string) (*cloudscheduler.Job, error) {
	args := m.Called(ctx, name)
	return args.Get(0).(*cloudscheduler.Job), args.Error(1)
}

func (m *SchedulerMock) DeleteScheduledJob(ctx context.Context, name string) error {
	args := m.Called(ctx, name)
	return args.Error(0)
}

func (m *SchedulerMock) RunScheduledJob(ctx context.Context, name string) (*cloudscheduler.Job, error) {
	args := m.Called(ctx, name)
	return args.Get(0).(*cloudscheduler.Job), args.Error(1)
}
