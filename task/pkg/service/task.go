package service

import (
	"context"
	"firebase.google.com/go/auth"
	v1 "github.com/jakubjano/todolist/apis/go-sdk/task/v1"
	"google.golang.org/protobuf/types/known/emptypb"
	"jakubjano/todolist/task/pkg/service/repository"
	"log"
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
		log.Fatalf("Error: %v", err)
	}
	return repository.ToApi(task), nil
}

func (ts *TaskService) GetTask(ctx context.Context, in *v1.GetTaskRequest) (*v1.Task, error) {
	task, err := ts.taskRepo.Get(ctx, in.TaskID)
	if err != nil {
		log.Fatalf("Error getting task with id %s: %v", in.TaskID, err)
		return &v1.Task{}, err
	}
	return repository.ToApi(task), nil
}

func (ts *TaskService) UpdateTask(ctx context.Context, in *v1.UpdateTaskRequest) (*v1.Task, error) {
	fields := map[string]interface{}{
		// passing keys like this does not seem right
		"Name":        in.NewName,
		"Description": in.NewDescription,
		"Time":        in.NewTime}
	task, err := ts.taskRepo.Update(ctx, fields, in.TaskID)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	return repository.ToApi(task), nil
}

//func (ts *TaskService) UpdateTask(ctx context.Context, in *v1.Task) (*v1.Task, error) {
//todo remake message for update request - tried but worked as expected
// takes whole Task structure instead of UpdateTaskRequest and replaces only new values

//	task, err := ts.taskRepo.Update(ctx, repository.TaskFromMsg(in), in.TaskID)
//	if err != nil {
//		log.Fatalf("Error updating task: %v \n", err)
//	}
//	return repository.ToApi(task), nil
//}

func (ts *TaskService) DeleteTask(ctx context.Context, in *v1.DeleteTaskRequest) (*emptypb.Empty, error) {
	err := ts.taskRepo.Delete(ctx, in.TaskID)
	if err != nil {
		log.Fatalf("error deleting task with id %s: %v \n", in.TaskID, err)
		return &emptypb.Empty{}, err
	}
	return &emptypb.Empty{}, nil
}
