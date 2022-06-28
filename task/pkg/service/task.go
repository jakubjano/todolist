package service

import (
	"context"
	"firebase.google.com/go/auth"
	"fmt"
	v1 "github.com/jakubjano/todolist/apis/go-sdk/task/v1"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	middleware "jakubjano/todolist/task/internal/auth"
	"jakubjano/todolist/task/pkg/service/repository"
	"net/http"
)

//type TaskClientInterface interface {
//}

type TaskService struct {
	v1.UnimplementedTaskServiceServer
	authClient *auth.Client
	taskRepo   repository.FSTaskInterface
}

func NewTaskService(authClient *auth.Client, taskRepo repository.FSTaskInterface) *TaskService {
	return &TaskService{
		authClient: authClient,
		taskRepo:   taskRepo,
	}
}

func (ts *TaskService) CreateTask(ctx context.Context, in *v1.Task) (*v1.Task, error) {
	//todo enable admin create tasks for others?
	//userCtx := ctx.Value("user").(*middleware.UserContext)
	task, err := ts.taskRepo.Create(ctx, repository.TaskFromMsg(in))
	if err != nil {
		fmt.Printf("error creating task: %v \n", err)
		return &v1.Task{}, status.Error(http.StatusInternalServerError, err.Error())

	}
	return repository.ToApi(task), nil
}

func (ts *TaskService) GetTask(ctx context.Context, in *v1.GetTaskRequest) (*v1.Task, error) {
	userCtx := ctx.Value("user").(*middleware.UserContext)
	task, err := ts.taskRepo.Get(ctx, in.TaskId)
	if err != nil {
		fmt.Printf("error getting task with id %s: %v \n", in.TaskId, err)
		return &v1.Task{}, err
	}
	switch userCtx.Role {
	case "admin":
		fmt.Println("admin authorized")
	case "user":
		if task.UserID != userCtx.UserID {
			return &v1.Task{}, status.Error(http.StatusUnauthorized, ErrUnauthorized.Error())
		}
	}
	return repository.ToApi(task), nil
}

func (ts *TaskService) UpdateTask(ctx context.Context, in *v1.UpdateTaskRequest) (*v1.Task, error) {
	userCtx := ctx.Value("user").(*middleware.UserContext)
	taskCheck, err := ts.taskRepo.Get(ctx, in.TaskId)
	if err != nil {
		fmt.Printf("error getting task with id %s: %v \n", in.TaskId, err)
		return &v1.Task{}, status.Error(http.StatusInternalServerError, err.Error())
	}
	switch userCtx.Role {
	case "admin":
		fmt.Println("admin authorized")
	case "user":
		if taskCheck.UserID != userCtx.UserID {
			return &v1.Task{}, status.Error(http.StatusUnauthorized, ErrUnauthorized.Error())
		}
	}
	fields := map[string]interface{}{
		"name":        in.NewName,
		"description": in.NewDescription,
		"time":        in.NewTime}
	task, err := ts.taskRepo.Update(ctx, fields, in.TaskId)
	if err != nil {
		fmt.Printf("error updating task with id %s: %v \n", in.TaskId, err)
		return &v1.Task{}, status.Error(http.StatusInternalServerError, err.Error())
	}
	return repository.ToApi(task), nil
}

//func (ts *TaskService) UpdateTask(ctx context.Context, in *v1.Task) (*v1.Task, error) {
//todo remake message for update request - tried but worked as expected
// takes whole Task structure instead of UpdateTaskRequest and replaces only new values

//	task, err := ts.taskRepo.Update(ctx, repository.TaskFromMsg(in), in.taskID)
//	if err != nil {
//		fmt.Printf("Error updating task: %v \n", err)
//	}
//	return repository.ToApi(task), nil
//}

func (ts *TaskService) DeleteTask(ctx context.Context, in *v1.DeleteTaskRequest) (*emptypb.Empty, error) {
	userCtx := ctx.Value("user").(*middleware.UserContext)
	taskCheck, err := ts.taskRepo.Get(ctx, in.TaskId)
	switch userCtx.Role {
	case "admin":
		fmt.Println("admin authorized")
	case "user":
		if taskCheck.UserID != userCtx.UserID {
			return &emptypb.Empty{}, status.Error(http.StatusUnauthorized, ErrUnauthorized.Error())
		}
	}
	err = ts.taskRepo.Delete(ctx, in.TaskId)
	if err != nil {
		fmt.Printf("error deleting task with id %s: %v \n", in.TaskId, err)
		return &emptypb.Empty{}, status.Error(http.StatusInternalServerError, err.Error())
	}
	return &emptypb.Empty{}, nil
}
