package repository

import (
	v1 "github.com/jakubjano/todolist/apis/go-sdk/task/v1"
)

const (
	CollectionTasks = "tasks"
	CollectionUsers = "users"
	TaskList        = "task_list"
)

type Task struct {
	CreatedAt    int64  `firestore:"createdAt"`
	Name         string `firestore:"name"`
	Description  string `firestore:"description"`
	UserID       string `firestore:"userID"`
	UserEmail    string `firestore:"email"`
	Time         int64  `firestore:"time"`
	TaskID       string `firestore:"taskID"`
	ReminderSent bool   `firestore:"reminderSent"`
}

// User type redefined in the task microservice to maintain its independence on the user microservice
type User struct {
	UserID    string `firestore:"userID"`
	Email     string `firestore:"email"`
	FirstName string `firestore:"firstName"`
	LastName  string `firestore:"lastName"`
	Phone     string `firestore:"phone"`
	Address   string `firestore:"address"`
}

func TaskFromMsg(msg *v1.Task) Task {
	return Task{
		CreatedAt:   msg.CreatedAt,
		Name:        msg.Name,
		Description: msg.Description,
		UserID:      msg.UserId,
		Time:        msg.Time,
		TaskID:      msg.TaskId,
		UserEmail:   msg.UserEmail,
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
		UserEmail:   task.UserEmail,
	}
}

func SliceToApi(tasks []Task) *v1.TaskList {
	apiTasks := make([]*v1.Task, len(tasks))
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
	return &v1.TaskList{Tasks: apiTasks}
}
