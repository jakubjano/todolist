package service

import (
	"context"
	"fmt"
	v1 "github.com/jakubjano/todolist/apis/go-sdk/task/v1"
	"go.uber.org/zap"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	middleware "jakubjano/todolist/task/internal/auth"
	"jakubjano/todolist/task/pkg/service/repository"
	"net/http"
)

type TaskService struct {
	v1.UnimplementedTaskServiceServer
	taskRepo repository.FSTaskInterface
	logger   *zap.Logger
}

func NewTaskService(taskRepo repository.FSTaskInterface, logger *zap.Logger) *TaskService {
	return &TaskService{
		taskRepo: taskRepo,
		logger:   logger,
	}
}

func (ts *TaskService) CreateTask(ctx context.Context, in *v1.Task) (*v1.Task, error) {
	userCtx := ctx.Value(middleware.ContextUser).(*middleware.UserContext)
	log := ts.logger.With(
		zap.String("caller_email", userCtx.Email),
		zap.String("caller_id", userCtx.UserID),
	)
	in.UserId = userCtx.UserID
	in.UserEmail = userCtx.Email
	task, err := ts.taskRepo.Create(ctx, repository.TaskFromMsg(in))
	if err != nil {
		log.Error(err.Error(), zap.String("task_id", task.TaskID))
		return &v1.Task{}, status.Error(http.StatusInternalServerError, err.Error())
	}
	log.Info("Created task ", zap.String("task_id", task.TaskID))
	return repository.ToApi(task), nil
}

func (ts *TaskService) GetTask(ctx context.Context, in *v1.GetTaskRequest) (*v1.Task, error) {
	userCtx := ctx.Value(middleware.ContextUser).(*middleware.UserContext)
	log := ts.logger.With(
		zap.String("caller_email", userCtx.Email),
		zap.String("caller_id", userCtx.UserID),
		zap.String("task_id", in.TaskId),
	)
	task, err := ts.taskRepo.Get(ctx, userCtx.UserID, in.TaskId)
	if err != nil {
		log.Error(err.Error())
		return &v1.Task{}, status.Error(http.StatusBadRequest, err.Error())
	}
	return repository.ToApi(task), nil
}

func (ts *TaskService) UpdateTask(ctx context.Context, in *v1.Task) (*v1.Task, error) {
	userCtx := ctx.Value(middleware.ContextUser).(*middleware.UserContext)
	log := ts.logger.With(
		zap.String("caller_email", userCtx.Email),
		zap.String("caller_id", userCtx.UserID),
		zap.String("task_id", in.TaskId),
	)
	task, err := ts.taskRepo.Update(ctx, repository.TaskFromMsg(in), userCtx.UserID, in.TaskId)
	log.Info("Updated task ")
	if err != nil {
		log.Error(err.Error())
		return &v1.Task{}, status.Error(http.StatusInternalServerError, err.Error())
	}
	return repository.ToApi(task), nil
}

func (ts *TaskService) DeleteTask(ctx context.Context, in *v1.DeleteTaskRequest) (*emptypb.Empty, error) {
	userCtx := ctx.Value(middleware.ContextUser).(*middleware.UserContext)
	log := ts.logger.With(
		zap.String("caller_email", userCtx.Email),
		zap.String("caller_id", userCtx.UserID),
		zap.String("task_id", in.TaskId),
	)
	err := ts.taskRepo.Delete(ctx, userCtx.UserID, in.TaskId)
	log.Info("Deleted task ")
	if err != nil {
		log.Error(err.Error())
		return &emptypb.Empty{}, status.Error(http.StatusInternalServerError, err.Error())
	}
	return &emptypb.Empty{}, nil
}

func (ts *TaskService) GetLastN(ctx context.Context, in *v1.GetLastNRequest) (*v1.TaskList, error) {
	userCtx := ctx.Value(middleware.ContextUser).(*middleware.UserContext)
	log := ts.logger.With(
		zap.String("caller_email", userCtx.Email),
		zap.String("caller_id", userCtx.UserID),
	)
	tasks, err := ts.taskRepo.GetLastN(ctx, userCtx.UserID, in.N)
	if err != nil {
		log.Error(err.Error())
		return &v1.TaskList{Tasks: nil}, err
	}
	return repository.SliceToApi(tasks), nil
}

func (ts *TaskService) GetExpired(ctx context.Context, in *v1.GetExpiredRequest) (*v1.TaskList, error) {
	userCtx := ctx.Value(middleware.ContextUser).(*middleware.UserContext)
	log := ts.logger.With(
		zap.String("caller_email", userCtx.Email),
		zap.String("caller_id", userCtx.UserID),
	)
	tasks, err := ts.taskRepo.GetExpired(ctx, userCtx.UserID)
	if err != nil {
		log.Error(err.Error())
		return &v1.TaskList{Tasks: nil}, err
	}
	fmt.Println(repository.SliceToApi(tasks))
	return repository.SliceToApi(tasks), nil
}
