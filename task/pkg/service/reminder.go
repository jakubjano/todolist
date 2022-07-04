package service

import (
	"context"
	"go.uber.org/zap"
	"jakubjano/todolist/task/pkg/service/repository"
	"net/smtp"
)

type Reminder struct {
	taskRepo  repository.FSTaskInterface
	logger    *zap.Logger
	emailAuth smtp.Auth
}

func NewReminder(taskRepo repository.FSTaskInterface, logger *zap.Logger, emailAuth smtp.Auth) *Reminder {
	return &Reminder{
		taskRepo:  taskRepo,
		logger:    logger,
		emailAuth: emailAuth,
	}
}

// RemindUserViaEmail checks if there are any reminders to send out to users
// then iterates through each task and sends it via smtp to the corresponding email address with a prebuilt message
// after the reminders are sent, RemindUserViaEmail flags the tasks with boolean and updates the database
func (r *Reminder) RemindUserViaEmail(ctx context.Context, host, port, from string) (map[string][]repository.Task, error) {
	sentReminders := make(map[string][]repository.Task)
	reminders, err := r.taskRepo.SearchForExpiringTasks(ctx)
	if err != nil {
		r.logger.Error(err.Error())
		return sentReminders, err
	}
	if len(reminders) < 1 {
		r.logger.Info("No expiring tasks")
		return sentReminders, nil
	}
	for email, tasks := range reminders {
		for i, task := range tasks {
			log := r.logger.With(
				zap.String("email", email),
				zap.String("task", task.Name),
			)
			message := []byte("Your task is expiring soon: " + task.Name)
			err = smtp.SendMail(host+port, r.emailAuth, from, []string{email}, message)
			if err != nil {
				log.Error(err.Error())
				return sentReminders, err
			}
			tasks[i].ReminderSent = true
			updatedTask, err := r.taskRepo.Update(ctx, tasks[i], tasks[i].UserID, tasks[i].TaskID)
			if err != nil {
				log.Error(err.Error())
			}
			sentReminders[email] = append(sentReminders[email], updatedTask)
			log.Info("reminder sent")
		}

	}
	return sentReminders, nil
}
