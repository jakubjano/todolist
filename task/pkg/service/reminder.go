package service

import (
	"cloud.google.com/go/firestore"
	"context"
	"github.com/jakubjano/todolist/task/pkg/service/repository"
	"go.uber.org/zap"
	"net/smtp"
)

const (
	batchSize = 500
)

type Reminder struct {
	taskRepo    repository.FSTaskInterface
	logger      *zap.Logger
	emailSender EmailSender
	fs          *firestore.Client
}

type EmailSender interface {
	Send(to []string, message []byte) error
}

type Settings struct {
	Host     string
	Port     string
	From     string
	UserName string
	Password string
}

type emailSender struct {
	emailSettings *Settings
}

func (e *emailSender) Send(to []string, message []byte) error {
	addr := e.emailSettings.Host + ":" + e.emailSettings.Port
	auth := smtp.PlainAuth("", e.emailSettings.UserName, e.emailSettings.Password, e.emailSettings.Host)
	return smtp.SendMail(addr, auth, e.emailSettings.From, to, message)
}

func NewEmailSender(emailSetting *Settings) EmailSender {
	return &emailSender{emailSetting}
}

func NewReminder(taskRepo repository.FSTaskInterface, logger *zap.Logger, emailSender EmailSender,
	fs *firestore.Client) *Reminder {
	return &Reminder{
		taskRepo:    taskRepo,
		logger:      logger,
		emailSender: emailSender,
		fs:          fs,
	}
}

// RemindUserViaEmail checks if there are any reminders to send out to users
// then iterates through each task and sends it via smtp to the corresponding email address with a prebuilt message
// after the reminders are sent, RemindUserViaEmail flags the tasks with boolean and updates the database
func (r *Reminder) RemindUserViaEmail(ctx context.Context) error {
	reminders, err := r.taskRepo.SearchForExpiringTasks(ctx)
	if err != nil {
		r.logger.Error(err.Error())
		return err
	}
	if len(reminders) < 1 {
		r.logger.Info("No expiring tasks")
		return nil
	}
	batch := r.fs.Batch()
	taskCount := 0
	mapLength := len(reminders)
	emailCount := 0
	for email, tasks := range reminders {
		emailCount += 1
		for i, task := range tasks {
			log := r.logger.With(
				zap.String("email", email),
				zap.String("task", task.Name),
			)
			message := []byte("Your task is expiring soon: " + task.Name)
			err = r.emailSender.Send([]string{email}, message)
			if err != nil {
				log.Error(err.Error())
				return err
			}
			log.Info("reminder sent")
			//update task_list duplicate collection
			batch.Set(r.fs.Collection(repository.TaskList).Doc(task.TaskID), map[string]interface{}{
				"reminderSent": true,
			}, firestore.MergeAll)
			// increment task count after sending reminder and adding Set operation into batch
			// commit when taskCount reaches the max limit of operations in batch write or when the iteration
			//		through the entire map is done (last iteration over tasks in the last email key)
			taskCount += 1
			if taskCount == batchSize || (emailCount == mapLength && i == len(tasks)-1) {
				_, err = batch.Commit(ctx)
				if err != nil {
					log.Error(err.Error())
					return err
				}
				taskCount = 0
			}
		}
	}
	return nil
}
