package repository

import (
	v1 "github.com/jakubjano/todolist/apis/go-sdk/task/v1"
)

const (
	CollectionTasks = "tasks"
	CollectionUsers = "users"
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
		UserID:      msg.UserId,
		Time:        msg.Time,
		TaskID:      msg.TaskId,
	}
}

func ToApi(task Task) *v1.Task {
	return &v1.Task{
		TaskId:      task.TaskID,
		CreatedAt:   task.CreatedAt,
		Name:        task.Name,
		Description: task.Description,
		Time:        task.Time,
		UserId:      task.UserID,
	}
}

func SliceToApi(tasks []Task) *v1.TaskSlice {
	apiTasks := v1.TaskSlice{}.Tasks
	for _, task := range tasks {
		apiTasks = append(apiTasks, &v1.Task{
			TaskId:      task.TaskID,
			CreatedAt:   task.CreatedAt,
			Name:        task.Name,
			Description: task.Description,
			Time:        task.Time,
			UserId:      task.UserID,
		})
	}
	return &v1.TaskSlice{Tasks: apiTasks}
}
