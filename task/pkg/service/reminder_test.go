//go:build integration
// +build integration

package service

import (
	"cloud.google.com/go/firestore"
	"context"
	"github.com/jakubjano/todolist/task/pkg/service/repository"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"os"
	"testing"
	"time"
)

type ReminderTestSuite struct {
	suite.Suite
	client      *firestore.Client
	taskRepo    repository.FSTaskInterface
	clientMock  *ClientMock
	reminder    Reminder
	emailSender EmailSender
}

func (s *ReminderTestSuite) SetupSuite() {
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, os.Getenv("PROJECT_ID"))
	s.NoError(err)
	taskRepo := repository.NewFSTask(client.Collection(repository.CollectionUsers), client)
	logger, err := NewLogger()
	s.NoError(err)
	clientMock := NewClientMock()
	reminder := NewReminder(taskRepo, logger, clientMock, client)
	s.reminder = reminder
	s.client = client
	s.taskRepo = taskRepo
	s.clientMock = clientMock

}

func (s *ReminderTestSuite) SetupTest() {
	// Add test data to DB
	ctx := context.Background()
	batchCreate := s.client.Batch()
	users := []repository.User{
		{
			UserID:    "1",
			Email:     "example1@tst.com",
			FirstName: "fn1",
			LastName:  "ln1",
			Phone:     "p1",
			Address:   "a1",
		},

		{
			UserID:    "2",
			Email:     "example2@tst.com",
			FirstName: "fn2",
			LastName:  "ln2",
			Phone:     "p2",
			Address:   "a2",
		},

		{
			UserID:    "3",
			Email:     "example3@tst.com",
			FirstName: "fn3",
			LastName:  "ln3",
			Phone:     "p3",
			Address:   "a3",
		},

		{
			UserID:    "4",
			Email:     "example6@tst.com",
			FirstName: "fn4",
			LastName:  "ln4",
			Phone:     "p4",
			Address:   "a4",
		},

		{
			UserID:    "5",
			Email:     "example5@tst.com",
			FirstName: "fn5",
			LastName:  "ln5",
			Phone:     "p5",
			Address:   "a5",
		},

		// users for reminder batch write testing
		{
			UserID:    "6",
			Email:     "example6@tst.com",
			FirstName: "fn6",
			LastName:  "ln6",
			Phone:     "p6",
			Address:   "a6",
		},
		{
			UserID:    "7",
			Email:     "example7@tst.com",
			FirstName: "fn7",
			LastName:  "ln7",
			Phone:     "p7",
			Address:   "a7",
		},
	}

	tasks := []repository.Task{
		{
			CreatedAt:    1,
			Name:         "task1",
			Description:  "desc1",
			UserID:       "1",
			UserEmail:    "example1@tst.com",
			Time:         7,
			TaskID:       "tid1",
			ReminderSent: false,
		},
		{
			CreatedAt:    2,
			Name:         "task2",
			Description:  "desc2",
			UserID:       "1",
			UserEmail:    "example1@tst.com",
			Time:         time.Now().Add(time.Minute * 1).Unix(),
			TaskID:       "tid2",
			ReminderSent: false,
		},
		{
			CreatedAt:    3,
			Name:         "task3",
			Description:  "desc3",
			UserID:       "1",
			UserEmail:    "example1@tst.com",
			Time:         10,
			TaskID:       "tid3",
			ReminderSent: false,
		},
		{
			CreatedAt:    5,
			Name:         "task4",
			Description:  "desc4",
			UserID:       "2",
			UserEmail:    "example2@tst.com",
			Time:         10,
			TaskID:       "tid4",
			ReminderSent: false,
		},
		{
			CreatedAt:    5,
			Name:         "task5",
			Description:  "desc5",
			UserID:       "3",
			UserEmail:    "example3@tst.com",
			Time:         10,
			TaskID:       "tid5",
			ReminderSent: false,
		},
		{
			CreatedAt:    5,
			Name:         "task6",
			Description:  "desc6",
			UserID:       "4",
			UserEmail:    "example6@tst.com",
			Time:         10,
			TaskID:       "tid6",
			ReminderSent: false,
		},
		// tasks for reminder batch write testing
		// User 6
		{
			CreatedAt:    1,
			Name:         "task7",
			Description:  "desc7",
			UserID:       "6",
			UserEmail:    "example6@tst.com",
			Time:         time.Now().Add(time.Minute * 4).Unix(),
			TaskID:       "tid7",
			ReminderSent: false,
		},
		{
			CreatedAt:    1,
			Name:         "task8",
			Description:  "desc8",
			UserID:       "6",
			UserEmail:    "example6@tst.com",
			Time:         time.Now().Add(time.Minute * 4).Unix(),
			TaskID:       "tid8",
			ReminderSent: false,
		},
		{
			CreatedAt:    1,
			Name:         "task9",
			Description:  "desc9",
			UserID:       "6",
			UserEmail:    "example6@tst.com",
			Time:         time.Now().Add(time.Minute * 4).Unix(),
			TaskID:       "tid9",
			ReminderSent: false,
		},
		{
			CreatedAt:    1,
			Name:         "task10",
			Description:  "desc10",
			UserID:       "6",
			UserEmail:    "example6@tst.com",
			Time:         time.Now().Add(time.Minute * 4).Unix(),
			TaskID:       "tid10",
			ReminderSent: false,
		},
		{
			CreatedAt:    1,
			Name:         "task11",
			Description:  "desc11",
			UserID:       "6",
			UserEmail:    "example6@tst.com",
			Time:         time.Now().Add(time.Minute * 4).Unix(),
			TaskID:       "tid11",
			ReminderSent: false,
		},
		{
			CreatedAt:    1,
			Name:         "task12",
			Description:  "desc12",
			UserID:       "6",
			UserEmail:    "example6@tst.com",
			Time:         time.Now().Add(time.Minute * 4).Unix(),
			TaskID:       "tid12",
			ReminderSent: false,
		},
		// User 7
		{
			CreatedAt:    1,
			Name:         "task13",
			Description:  "desc13",
			UserID:       "7",
			UserEmail:    "example7@tst.com",
			Time:         time.Now().Add(time.Minute * 3).Unix(),
			TaskID:       "tid13",
			ReminderSent: false,
		},
		{
			CreatedAt:    1,
			Name:         "task14",
			Description:  "desc14",
			UserID:       "7",
			UserEmail:    "example7@tst.com",
			Time:         time.Now().Add(time.Minute * 3).Unix(),
			TaskID:       "tid14",
			ReminderSent: false,
		},
		{
			CreatedAt:    1,
			Name:         "task15",
			Description:  "desc15",
			UserID:       "7",
			UserEmail:    "example7@tst.com",
			Time:         time.Now().Add(time.Minute * 3).Unix(),
			TaskID:       "tid15",
			ReminderSent: false,
		},
		{
			CreatedAt:    1,
			Name:         "task16",
			Description:  "desc16",
			UserID:       "7",
			UserEmail:    "example7@tst.com",
			Time:         time.Now().Add(time.Minute * 3).Unix(),
			TaskID:       "tid16",
			ReminderSent: false,
		},
		{
			CreatedAt:    1,
			Name:         "task17",
			Description:  "desc17",
			UserID:       "7",
			UserEmail:    "example7@tst.com",
			Time:         time.Now().Add(time.Minute * 3).Unix(),
			TaskID:       "tid17",
			ReminderSent: false,
		},
		{
			CreatedAt:    1,
			Name:         "task18",
			Description:  "desc18",
			UserID:       "7",
			UserEmail:    "example7@tst.com",
			Time:         time.Now().Add(time.Minute * 3).Unix(),
			TaskID:       "tid18",
			ReminderSent: false,
		},
	}

	for _, user := range users {
		userRef := s.client.Collection(repository.CollectionUsers).Doc(user.UserID)
		batchCreate.Set(userRef, user)
	}
	_, err := batchCreate.Commit(ctx)
	s.NoError(err)

	taskBatch := s.client.Batch()
	for _, task := range tasks {
		taskRef := s.client.Collection(repository.CollectionUsers).Doc(task.UserID).
			Collection(repository.CollectionTasks).Doc(task.TaskID)
		taskBatch.Set(taskRef, task)
	}
	_, err = taskBatch.Commit(ctx)
	s.NoError(err)

	taskListBatch := s.client.Batch()
	for _, task := range tasks {
		taskListRef := s.client.Collection(repository.TaskList).Doc(task.TaskID)
		taskListBatch.Set(taskListRef, task)
	}
	_, err = taskListBatch.Commit(ctx)
	s.NoError(err)
}

func (s *ReminderTestSuite) TearDownTest() {
	// clear all data from DB after every test
	// deleting users will delete nested tasks as well
	ctx := context.Background()
	docs, err := s.client.Collection(repository.CollectionUsers).Documents(ctx).GetAll()
	s.NoError(err)
	batch := s.client.Batch()
	for _, doc := range docs {
		batch.Delete(doc.Ref)
	}
	_, err = batch.Commit(ctx)
	s.NoError(err)

	// delete task list too
	taskListBatch := s.client.Batch()
	taskListDocs, err := s.client.Collection(repository.TaskList).Documents(ctx).GetAll()
	s.NoError(err)
	for _, doc := range taskListDocs {
		taskListBatch.Delete(doc.Ref)
	}
	_, err = taskListBatch.Commit(ctx)
	s.NoError(err)
}

func (s *ReminderTestSuite) TearDownSuite() {
	err := s.client.Close()
	s.NoError(err)
}

func (s *ReminderTestSuite) TestRemindUserViaEmail() {
	ctx := context.Background()
	s.clientMock.On("Send",
		mock.Anything,
		mock.Anything,
	).Return(nil)
	err := s.reminder.RemindUserViaEmail(ctx)
	s.clientMock.AssertNumberOfCalls(s.T(),
		"Send",
		13,
	)
	s.NoError(err)
}

func TestReminderTestSuite(t *testing.T) {
	suite.Run(t, new(ReminderTestSuite))
}
