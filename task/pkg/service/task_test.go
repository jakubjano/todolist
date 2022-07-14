package service

import (
	"context"
	v1 "github.com/jakubjano/todolist/apis/go-sdk/task/v1"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	middleware "jakubjano/todolist/task/internal/auth"
	"jakubjano/todolist/task/pkg/service/repository"
	"testing"
)

type ServiceTaskTestSuite struct {
	suite.Suite
	ts       *TaskService
	mockRepo *repository.FSTaskMock
}

func (s *ServiceTaskTestSuite) SetupSuite() {
	logger, _ := zap.NewProduction()
	taskRepo := repository.NewMockRepo()
	ts := NewTaskService(taskRepo, logger)
	s.mockRepo = taskRepo
	s.ts = ts
}

func (s *ServiceTaskTestSuite) TestCreateTask() {
	ctx := context.Background()
	candidates := []struct {
		ctx            context.Context
		in             *v1.Task
		expectedResult *v1.Task
		mockReturn     repository.Task
		expectedError  error
		expectedCode   codes.Code
	}{
		// valid input
		{
			ctx: context.WithValue(ctx, middleware.ContextUser, &middleware.UserContext{
				UserID: "1",
				Email:  "example1@tst.com",
				Role:   "user",
			}),
			in: &v1.Task{
				TaskId:      "tid1",
				CreatedAt:   1,
				Name:        "task1",
				Description: "task1 desc",
				Time:        5,
				UserId:      "1",
				UserEmail:   "example1@tst.com",
			},
			mockReturn: repository.Task{
				CreatedAt:    1,
				Name:         "task1",
				Description:  "task1 desc",
				UserID:       "1",
				UserEmail:    "example1@tst.com",
				Time:         5,
				TaskID:       "tid1",
				ReminderSent: false,
			},
			expectedResult: &v1.Task{
				TaskId:      "tid1",
				CreatedAt:   1,
				Name:        "task1",
				Description: "task1 desc",
				Time:        5,
				UserId:      "1",
				UserEmail:   "example1@tst.com",
			},
			expectedError: nil,
			expectedCode:  codes.OK,
		},
	}
	for i, candidate := range candidates {
		s.mockRepo.On("Create", candidate.ctx, repository.TaskFromMsg(candidate.in)).
			Return(candidate.mockReturn, candidate.expectedError)
		task, err := s.ts.CreateTask(candidate.ctx, candidate.in)
		task.TaskId = candidate.expectedResult.TaskId
		task.CreatedAt = candidate.expectedResult.CreatedAt
		s.mockRepo.AssertCalled(s.T(), "Create", candidate.ctx, repository.TaskFromMsg(candidate.in))
		s.Equalf(candidate.expectedResult, task, "candidate %d", i+1)
		s.Equalf(candidate.expectedCode, status.Code(err), "candidate %:", i+1)
	}
}

func (s *ServiceTaskTestSuite) TestGetTask() {
	ctx := context.Background()
	candidates := []struct {
		ctx            context.Context
		in             *v1.GetTaskRequest
		expectedResult *v1.Task
		testForError   bool
		expectedError  error
		expectedCode   codes.Code
	}{
		// valid input
		{
			ctx: context.WithValue(ctx, middleware.ContextUser, &middleware.UserContext{
				UserID: "1",
				Email:  "example1@tst.com",
				Role:   "user",
			}),
			in: &v1.GetTaskRequest{TaskId: "tid1"},
			expectedResult: &v1.Task{
				TaskId:      "tid1",
				CreatedAt:   1,
				Name:        "task1",
				Description: "task1 desc",
				Time:        2,
				UserId:      "1",
				UserEmail:   "example1@tst.com",
			},
			testForError:  false,
			expectedError: nil,
			expectedCode:  codes.OK,
		},
		// non-existing task
		{
			ctx: context.WithValue(ctx, middleware.ContextUser, &middleware.UserContext{
				UserID: "1",
				Email:  "example1@tst.com",
				Role:   "user",
			}),
			in:             &v1.GetTaskRequest{TaskId: "tid999"},
			expectedResult: &v1.Task{},
			testForError:   true,
			expectedError:  status.Error(codes.NotFound, ""),
			expectedCode:   codes.NotFound,
		},
		// different user id from context and task
		// won't find the task under different user
		{
			ctx: context.WithValue(ctx, middleware.ContextUser, &middleware.UserContext{
				UserID: "1",
				Email:  "example1@tst.com",
				Role:   "user",
			}),
			in:             &v1.GetTaskRequest{TaskId: "tid2"},
			expectedResult: &v1.Task{},
			testForError:   true,
			expectedError:  status.Error(codes.NotFound, ""),
			expectedCode:   codes.NotFound,
		},
		// invalid input
		{
			ctx: context.WithValue(ctx, middleware.ContextUser, &middleware.UserContext{
				UserID: "1",
				Email:  "example1@tst.com",
				Role:   "user",
			}),
			in:             &v1.GetTaskRequest{TaskId: ""},
			expectedResult: &v1.Task{},
			testForError:   true,
			expectedError:  status.Error(codes.InvalidArgument, ""),
			expectedCode:   codes.InvalidArgument,
		},
	}
	for i, candidate := range candidates {
		userCtx := candidate.ctx.Value(middleware.ContextUser).(*middleware.UserContext)
		s.mockRepo.On("Get", candidate.ctx, userCtx.UserID, candidate.in.TaskId).Return(repository.Task{
			CreatedAt:    candidate.expectedResult.CreatedAt,
			Name:         candidate.expectedResult.Name,
			Description:  candidate.expectedResult.Description,
			UserID:       candidate.expectedResult.UserId,
			UserEmail:    candidate.expectedResult.UserEmail,
			Time:         candidate.expectedResult.Time,
			TaskID:       candidate.expectedResult.TaskId,
			ReminderSent: false,
		}, candidate.expectedError)
		task, err := s.ts.GetTask(candidate.ctx, candidate.in)
		// check if method was called correctly
		s.mockRepo.AssertCalled(s.T(), "Get", candidate.ctx, userCtx.UserID, candidate.in.TaskId)
		if candidate.testForError {
			s.Contains(err.Error(), candidate.expectedCode.String())
		}
		s.Equalf(candidate.expectedResult, task, "candidate %d", i+1)
	}
}

func (s *ServiceTaskTestSuite) TestUpdateTask() {
	ctx := context.Background()
	candidates := []struct {
		ctx            context.Context
		in             *v1.Task
		expectedResult *v1.Task
		expectedError  error
		expectedCode   codes.Code
	}{
		// valid input
		{
			ctx: context.WithValue(ctx, middleware.ContextUser, &middleware.UserContext{
				UserID: "1",
				Email:  "example1@tst.com",
				Role:   "user",
			}),
			in: &v1.Task{
				TaskId:      "tid1",
				CreatedAt:   1,
				Name:        "updated name",
				Description: "updated desc",
				Time:        21,
				UserId:      "1",
				UserEmail:   "example1@tst.com",
			},
			expectedResult: &v1.Task{
				TaskId:      "tid1",
				CreatedAt:   1,
				Name:        "updated name",
				Description: "updated desc",
				Time:        21,
				UserId:      "1",
				UserEmail:   "example1@tst.com",
			},
			expectedError: nil,
			expectedCode:  codes.OK,
		},
		// wrong task id
		// creates new task when provided taskID not found
		{
			ctx: context.WithValue(ctx, middleware.ContextUser, &middleware.UserContext{
				UserID: "1",
				Email:  "example1@tst.com",
				Role:   "user",
			}),
			in: &v1.Task{
				TaskId:      "non-existent task id",
				CreatedAt:   1,
				Name:        "updated name",
				Description: "updated desc",
				Time:        21,
				UserId:      "1",
				UserEmail:   "example1@tst.com",
			},
			expectedResult: &v1.Task{
				TaskId:      "non-existent task id",
				CreatedAt:   1,
				Name:        "updated name",
				Description: "updated desc",
				Time:        21,
				UserId:      "1",
				UserEmail:   "example1@tst.com",
			},
			expectedError: nil,
			expectedCode:  codes.OK,
		},
	}
	for i, candidate := range candidates {
		s.mockRepo.On("Update", candidate.ctx,
			repository.TaskFromMsg(candidate.in), candidate.in.UserId, candidate.in.TaskId).Return(repository.Task{
			CreatedAt:    candidate.expectedResult.CreatedAt,
			Name:         candidate.expectedResult.Name,
			Description:  candidate.expectedResult.Description,
			UserID:       candidate.expectedResult.UserId,
			UserEmail:    candidate.expectedResult.UserEmail,
			Time:         candidate.expectedResult.Time,
			TaskID:       candidate.expectedResult.TaskId,
			ReminderSent: false,
		}, candidate.expectedError)
		task, err := s.ts.UpdateTask(candidate.ctx, candidate.in)
		s.mockRepo.AssertCalled(s.T(), "Update", candidate.ctx, repository.TaskFromMsg(candidate.in),
			candidate.in.UserId, candidate.in.TaskId)
		s.Equalf(candidate.expectedResult, task, "candidate %d", i+1)
		s.Equalf(candidate.expectedCode, status.Code(err), "candidate %:", i+1)
	}
}

func (s *ServiceTaskTestSuite) TestDeleteTask() {
	ctx := context.Background()
	candidates := []struct {
		ctx           context.Context
		in            *v1.DeleteTaskRequest
		expectedError error
		expectedCode  codes.Code
	}{
		// valid input
		{
			ctx: context.WithValue(ctx, middleware.ContextUser, &middleware.UserContext{
				UserID: "user1",
				Email:  "example1@tst.com",
				Role:   "user",
			}),
			in:            &v1.DeleteTaskRequest{TaskId: "tid1"},
			expectedError: nil,
			expectedCode:  codes.OK,
		},
		// non existent task
		// delete doesn't do anything
		{
			ctx: context.WithValue(ctx, middleware.ContextUser, &middleware.UserContext{
				UserID: "1",
				Email:  "example1@tst.com",
				Role:   "user",
			}),
			in:            &v1.DeleteTaskRequest{TaskId: "tid999"},
			expectedError: nil,
			expectedCode:  codes.OK,
		},
	}
	for i, candidate := range candidates {
		userCtx := candidate.ctx.Value(middleware.ContextUser).(*middleware.UserContext)
		s.mockRepo.On("Delete", candidate.ctx, userCtx.UserID, candidate.in.TaskId).
			Return(candidate.expectedError)
		_, err := s.ts.DeleteTask(candidate.ctx, candidate.in)
		s.mockRepo.AssertCalled(s.T(), "Delete", candidate.ctx, userCtx.UserID, candidate.in.TaskId)
		s.Equalf(candidate.expectedCode, status.Code(err), "candidate %:", i+1)
	}
}

func (s *ServiceTaskTestSuite) TestGetLastNTasks() {
	ctx := context.Background()
	candidates := []struct {
		ctx            context.Context
		in             *v1.GetLastNRequest
		expectedResult *v1.TaskList
		mockReturn     []repository.Task
		expectedError  error
		expectedCode   codes.Code
	}{
		// valid input
		{
			ctx: context.WithValue(ctx, middleware.ContextUser, &middleware.UserContext{
				UserID: "1",
				Email:  "example1@tst.com",
				Role:   "user",
			}),
			in: &v1.GetLastNRequest{N: 2},
			mockReturn: []repository.Task{
				{
					CreatedAt:    2,
					Name:         "task2",
					Description:  "task2 desc",
					UserID:       "1",
					UserEmail:    "example1@tst.com",
					Time:         6,
					TaskID:       "tid2",
					ReminderSent: false,
				},
				{
					CreatedAt:    1,
					Name:         "task1",
					Description:  "task1 desc",
					UserID:       "1",
					UserEmail:    "example1@tst.com",
					Time:         5,
					TaskID:       "tid1",
					ReminderSent: false,
				},
			},
			expectedResult: &v1.TaskList{
				Tasks: []*v1.Task{
					{
						TaskId:      "tid2",
						CreatedAt:   2,
						Name:        "task2",
						Description: "task2 desc",
						Time:        6,
						UserId:      "1",
						UserEmail:   "example1@tst.com",
					},
					{
						TaskId:      "tid1",
						CreatedAt:   1,
						Name:        "task1",
						Description: "task1 desc",
						Time:        5,
						UserId:      "1",
						UserEmail:   "example1@tst.com",
					},
				},
			},
			expectedError: nil,
			expectedCode:  codes.OK,
		},
	}
	for i, candidate := range candidates {
		userCtx := candidate.ctx.Value(middleware.ContextUser).(*middleware.UserContext)
		s.mockRepo.On("GetLastN", candidate.ctx, userCtx.UserID, candidate.in.N).Return(candidate.mockReturn,
			candidate.expectedError)
		tasks, err := s.ts.GetLastN(candidate.ctx, candidate.in)
		s.mockRepo.AssertCalled(s.T(), "GetLastN", candidate.ctx, userCtx.UserID, candidate.in.N)
		s.Equalf(candidate.expectedResult, tasks, "candidate %d", i+1)
		s.Equalf(candidate.expectedCode, status.Code(err), "candidate %:", i+1)
	}
}

func (s *ServiceTaskTestSuite) TestGetExpiredTasks() {
	ctx := context.Background()
	candidates := []struct {
		ctx            context.Context
		in             *v1.GetExpiredRequest
		mockReturn     []repository.Task
		expectedResult *v1.TaskList
		expectedError  error
		expectedCode   codes.Code
	}{
		//
		{
			ctx: context.WithValue(ctx, middleware.ContextUser, &middleware.UserContext{
				UserID: "1",
				Email:  "example1@tst.com",
				Role:   "user",
			}),
			in: &v1.GetExpiredRequest{},
			mockReturn: []repository.Task{
				{
					CreatedAt:    1,
					Name:         "task1",
					Description:  "task1 desc",
					UserID:       "1",
					UserEmail:    "example1@tst.com",
					Time:         2,
					TaskID:       "tid1",
					ReminderSent: false,
				},
				{
					CreatedAt:    1,
					Name:         "task2",
					Description:  "task2 desc",
					UserID:       "1",
					UserEmail:    "example1@tst.com",
					Time:         3,
					TaskID:       "tid2",
					ReminderSent: false,
				},
			},
			expectedResult: &v1.TaskList{
				Tasks: []*v1.Task{
					{
						TaskId:      "tid1",
						CreatedAt:   1,
						Name:        "task1",
						Description: "task1 desc",
						Time:        2,
						UserId:      "1",
						UserEmail:   "example1@tst.com",
					},
					{
						TaskId:      "tid2",
						CreatedAt:   1,
						Name:        "task2",
						Description: "task2 desc",
						Time:        3,
						UserId:      "1",
						UserEmail:   "example1@tst.com",
					},
				},
			},
			expectedError: nil,
			expectedCode:  codes.OK,
		},
	}
	for i, candidate := range candidates {
		userCtx := candidate.ctx.Value(middleware.ContextUser).(*middleware.UserContext)
		s.mockRepo.On("GetExpired", candidate.ctx, userCtx.UserID).Return(candidate.mockReturn,
			candidate.expectedError)
		tasks, err := s.ts.GetExpired(candidate.ctx, candidate.in)
		s.mockRepo.AssertCalled(s.T(), "GetExpired", candidate.ctx, userCtx.UserID)
		s.Equalf(candidate.expectedResult, tasks, "candidate %d", i+1)
		s.Equalf(candidate.expectedCode, status.Code(err), "candidate %:", i+1)
	}
}

func TestServiceTaskTestSuite(t *testing.T) {
	suite.Run(t, new(ServiceTaskTestSuite))
}
