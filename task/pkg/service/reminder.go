package service

import (
	"context"
	"go.uber.org/zap"
	"jakubjano/todolist/task/pkg/service/repository"
	"net/smtp"
	"strings"
)

// get all tasks from the database
// query each user's tasks and get ones approaching deadline
// get user's ID and ID of the task that's about to expire
// log this information
// email those users
// automate via cron

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

func (r *Reminder) SendEmail(ctx context.Context, host, port, from string) error {
	//todo don't send reminders to the same user every 30 seconds
	// more reliable solution than interval cutting needs more data than email and task name
	// -> pass whole user and task into a nested map in repository.SearchForExpiredTasks and return more data(i.e ids) ?
	// -> then return ids of sent reminders and check them in the next run of cron
	reminders, err := r.taskRepo.SearchForExpiringTasks(ctx)
	if err != nil {
		return err
	}
	// temporary log
	if len(reminders) < 1 {
		r.logger.Info("No expiring tasks")
	}
	for user, tasks := range reminders {
		tasksJoined := strings.Join(tasks, "\n")
		log := r.logger.With(
			zap.String("user", user),
			zap.String("task_names", tasksJoined),
		)
		message := []byte("Some of your tasks are expiring soon: " + tasksJoined)
		err = smtp.SendMail(host+port, r.emailAuth, from, []string{user}, message)
		if err != nil {
			return err
		}
		log.Info("reminder sent")
	}
	return nil
}
