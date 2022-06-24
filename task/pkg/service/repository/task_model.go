package repository

import (
	v1 "github.com/jakubjano/todolist/apis/go-sdk/task/v1"
)

type Task struct {
	CreatedAt   int64  `firestore:"createdAt"`
	Name        string `firestore:"name"`
	Description string `firestore:"description"`
	UserID      string `firestore:"userID"`
	Time        int64  `firestore:"time"`
	TaskID      string `firestore:"taskID"`
}

func TaskFromMsg(msg *v1.Task) Task {
	return Task{
		CreatedAt:   msg.CreatedAt,
		Name:        msg.Name,
		Description: msg.Description,
		UserID:      msg.UserID,
		Time:        msg.Time,
		TaskID:      msg.TaskID,
	}
}

func ToApi(task Task) *v1.Task {
	return &v1.Task{
		TaskID:      task.TaskID,
		CreatedAt:   task.CreatedAt,
		Name:        task.Name,
		Description: task.Description,
		Time:        task.Time,
		UserID:      task.UserID,
	}
}
