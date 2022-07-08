package service

import (
	"cloud.google.com/go/firestore"
	"context"
	"go.uber.org/zap"
	"jakubjano/todolist/task/pkg/service/repository"
	"net/smtp"
)

type Reminder struct {
	taskRepo  repository.FSTaskInterface
	logger    *zap.Logger
	emailAuth smtp.Auth
	fs        *firestore.Client
}

func NewReminder(taskRepo repository.FSTaskInterface, logger *zap.Logger, emailAuth smtp.Auth, fs *firestore.Client) *Reminder {
	return &Reminder{
		taskRepo:  taskRepo,
		logger:    logger,
		emailAuth: emailAuth,
		fs:        fs,
	}
}

// RemindUserViaEmail checks if there are any reminders to send out to users
// then iterates through each task and sends it via smtp to the corresponding email address with a prebuilt message
// after the reminders are sent, RemindUserViaEmail flags the tasks with boolean and updates the database
func (r *Reminder) RemindUserViaEmail(ctx context.Context, host, port, from string) error {
	reminders, err := r.taskRepo.SearchForExpiringTasks(ctx)
	if err != nil {
		r.logger.Error(err.Error())
		return err
	}
	if len(reminders) < 1 {
		return nil
	}
	batch := r.fs.Batch()
	counter := 0
	for email, tasks := range reminders {
		for _, task := range tasks {
			log := r.logger.With(
				zap.String("email", email),
				zap.String("task", task.Name),
			)
			message := []byte("Your task is expiring soon: " + task.Name)
			err = smtp.SendMail(host+port, r.emailAuth, from, []string{email}, message)
			if err != nil {
				log.Error(err.Error())
				return err
			}
			log.Info("reminder sent")

			//update task_list duplicate collection
			batch.Set(r.fs.Collection(repository.TaskList).Doc(task.TaskID), map[string]interface{}{
				"reminderSent": true,
			}, firestore.MergeAll)
			counter++
			if counter == 500 {
				_, err = batch.Commit(ctx)
				if err != nil {
					r.logger.Error(err.Error())
					return err
				}
				continue
			}
		}
	}
	return nil
}
