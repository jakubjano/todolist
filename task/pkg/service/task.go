package service

import (
	"context"
	"firebase.google.com/go/auth"
	"fmt"
	v1 "github.com/jakubjano/todolist/apis/go-sdk/task/v1"
	"google.golang.org/protobuf/types/known/emptypb"
	"jakubjano/todolist/task/pkg/service/repository"
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
	//todo handle time
	task, err := ts.taskRepo.Create(ctx, repository.TaskFromMsg(in))
	if err != nil {
		fmt.Printf("error crating task: %v", err)
		return &v1.Task{}, err
	}
	return repository.ToApi(task), nil
}

func (ts *TaskService) GetTask(ctx context.Context, in *v1.GetTaskRequest) (*v1.Task, error) {
	task, err := ts.taskRepo.Get(ctx, in.TaskId)
	if err != nil {
		fmt.Printf("error getting task with id %s: %v", in.TaskId, err)
		return &v1.Task{}, err
	}
	return repository.ToApi(task), nil
}

func (ts *TaskService) UpdateTask(ctx context.Context, in *v1.UpdateTaskRequest) (*v1.Task, error) {
	fields := map[string]interface{}{
		"name":        in.NewName,
		"description": in.NewDescription,
		"time":        in.NewTime}
	task, err := ts.taskRepo.Update(ctx, fields, in.TaskId)
	if err != nil {
		fmt.Printf("error updating task with id %s: %v", in.TaskId, err)
		return &v1.Task{}, err
	}
	return repository.ToApi(task), nil
}

//func (ts *TaskService) UpdateTask(ctx context.Context, in *v1.Task) (*v1.Task, error) {
//todo remake message for update request - tried but worked as expected
// takes whole Task structure instead of UpdateTaskRequest and replaces only new values

//	task, err := ts.taskRepo.Update(ctx, repository.TaskFromMsg(in), in.TaskID)
//	if err != nil {
//		fmt.Printf("Error updating task: %v \n", err)
//	}
//	return repository.ToApi(task), nil
//}

func (ts *TaskService) DeleteTask(ctx context.Context, in *v1.DeleteTaskRequest) (*emptypb.Empty, error) {
	err := ts.taskRepo.Delete(ctx, in.TaskId)
	if err != nil {
		fmt.Printf("error deleting task with id %s: %v \n", in.TaskId, err)
		return &emptypb.Empty{}, err
	}
	return &emptypb.Empty{}, nil
}
