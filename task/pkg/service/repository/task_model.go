package repository

import (
	v1 "github.com/jakubjano/todolist/apis/go-sdk/task/v1"
)

const (
	COLLECTION_TASKS = "tasks"
)

type Task struct {
	CreatedAt   int64  `firestore:"createdAt"`
	Name        string `firestore:"name"`
	Description string `firestore:"description"`
	UserId      string `firestore:"userId"`
	Time        int64  `firestore:"time"`
	TaskId      string `firestore:"taskId"`
}

func TaskFromMsg(msg *v1.Task) Task {
	return Task{
		CreatedAt:   msg.CreatedAt,
		Name:        msg.Name,
		Description: msg.Description,
		UserId:      msg.UserId,
		Time:        msg.Time,
		TaskId:      msg.TaskId,
	}
}

func ToApi(task Task) *v1.Task {
	return &v1.Task{
		TaskId:      task.TaskId,
		CreatedAt:   task.CreatedAt,
		Name:        task.Name,
		Description: task.Description,
		Time:        task.Time,
		UserId:      task.UserId,
	}
}
