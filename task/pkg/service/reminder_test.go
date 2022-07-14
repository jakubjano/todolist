//go:build integration
// +build integration

package service

import (
	"cloud.google.com/go/firestore"
	"context"
	"github.com/stretchr/testify/suite"
	"jakubjano/todolist/task/pkg/service/repository"
	"net/smtp"
	"os"
	"testing"
)

type ReminderTestSuite struct {
	suite.Suite
	client     *firestore.Client
	taskRepo   repository.FSTaskInterface
	clientMock *ClientMock
	testAuth   smtp.Auth
}

func (s *ReminderTestSuite) SetupSuite() {
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, os.Getenv("PROJECT_ID"))
	s.NoError(err)
	taskRepo := repository.NewFSTask(client.Collection(repository.CollectionUsers), client)
	clientMock := NewClientMock()
	testAuth := smtp.PlainAuth("", "username", "password", "host")
	s.client = client
	s.taskRepo = taskRepo
	s.clientMock = clientMock
	s.testAuth = testAuth
}

func (s *ReminderTestSuite) TestRemindUserViaEmail() {
	ctx := context.Background()
	tasks, err := s.taskRepo.SearchForExpiringTasks(ctx)
	s.NoError(err)
	batch := s.client.Batch()
	for email, tasks := range tasks {
		for i, task := range tasks {
			message := []byte("Your task is expiring soon: " + task.Name)
			s.clientMock.On("SendMail",
				"test_host",
				s.testAuth,
				"test_from",
				[]string{email},
				message).Return(nil)
			//update task_list duplicate collection
			batch.Set(s.client.Collection(repository.TaskList).Doc(task.TaskID), map[string]interface{}{
				"reminderSent": true,
			}, firestore.MergeAll)
			_, err = batch.Commit(ctx)
			s.NoError(err)

			taskCheck, err := s.taskRepo.Get(ctx, task.UserID, task.TaskID)
			s.NoError(err)
			s.Equalf(true, taskCheck.ReminderSent, "task number %d of user %s", i+1, task.UserEmail)
		}
	}
}

func TestReminderTestSuite(t *testing.T) {
	suite.Run(t, new(ReminderTestSuite))
}
