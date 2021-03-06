//go:build integration
// +build integration

package repository

import (
	"cloud.google.com/go/firestore"
	"context"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"os"
	"testing"
	"time"
)

type RepoTaskTestSuite struct {
	suite.Suite
	client   *firestore.Client
	taskRepo FSTaskInterface
}

// runs once at the beginning
func (s *RepoTaskTestSuite) SetupSuite() {
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, os.Getenv("PROJECT_ID"))
	s.NoError(err)
	taskRepo := NewFSTask(client.Collection(CollectionUsers), client)
	s.client = client
	s.taskRepo = taskRepo

}

// runs before every test
func (s *RepoTaskTestSuite) SetupTest() {
	// Add test data to DB
	ctx := context.Background()
	batchCreate := s.client.Batch()
	users := []User{
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

	tasks := []Task{
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
		userRef := s.client.Collection(CollectionUsers).Doc(user.UserID)
		batchCreate.Set(userRef, user)
	}
	_, err := batchCreate.Commit(ctx)
	s.NoError(err)

	taskBatch := s.client.Batch()
	for _, task := range tasks {
		taskRef := s.client.Collection(CollectionUsers).Doc(task.UserID).Collection(CollectionTasks).Doc(task.TaskID)
		taskBatch.Set(taskRef, task)
	}
	_, err = taskBatch.Commit(ctx)
	s.NoError(err)

	taskListBatch := s.client.Batch()
	for _, task := range tasks {
		taskListRef := s.client.Collection(TaskList).Doc(task.TaskID)
		taskListBatch.Set(taskListRef, task)
	}
	_, err = taskListBatch.Commit(ctx)
	s.NoError(err)
}

func (s *RepoTaskTestSuite) TearDownTest() {
	// clear all data from DB after every test
	// deleting users will delete nested tasks as well
	ctx := context.Background()
	docs, err := s.client.Collection(CollectionUsers).Documents(ctx).GetAll()
	s.NoError(err)
	batch := s.client.Batch()
	for _, doc := range docs {
		batch.Delete(doc.Ref)
	}
	_, err = batch.Commit(ctx)
	s.NoError(err)

	// delete task list too
	taskListBatch := s.client.Batch()
	taskListDocs, err := s.client.Collection(TaskList).Documents(ctx).GetAll()
	s.NoError(err)
	for _, doc := range taskListDocs {
		taskListBatch.Delete(doc.Ref)
	}
	_, err = taskListBatch.Commit(ctx)
	s.NoError(err)
}

func (s *RepoTaskTestSuite) TearDownSuite() {
	err := s.client.Close()
	s.NoError(err)
}

func (s *RepoTaskTestSuite) TestGetTask() {
	ctx := context.Background()
	candidates := []struct {
		taskID         string
		userID         string
		expectedResult Task
		expectedCode   codes.Code
	}{
		// valid input, exists
		{
			taskID: "tid1",
			userID: "1",
			expectedResult: Task{
				CreatedAt:    1,
				Name:         "task1",
				Description:  "desc1",
				UserID:       "1",
				UserEmail:    "example1@tst.com",
				Time:         7,
				TaskID:       "tid1",
				ReminderSent: false,
			},
			expectedCode: codes.OK,
		},
		// valid input, task does not exist
		{
			taskID:         "tid999",
			userID:         "1",
			expectedResult: Task{},
			expectedCode:   codes.NotFound,
		},
		// valid input, user does not exist
		{
			taskID:         "tid3",
			userID:         "9999",
			expectedResult: Task{},
			expectedCode:   codes.NotFound,
		},
		// invalid input
		{
			taskID:         "",
			expectedResult: Task{},
			expectedCode:   codes.InvalidArgument,
		},
	}
	for i, candidate := range candidates {
		task, err := s.taskRepo.Get(ctx, candidate.userID, candidate.taskID)
		s.NotNilf(task.TaskID, "candidate %d", i+1)
		task.TaskID = candidate.expectedResult.TaskID
		s.Equalf(candidate.expectedResult, task, "candidate %d", i+1)
		s.Equalf(candidate.expectedCode, status.Code(err), "candidate %d", i+1)
	}
}

func (s *RepoTaskTestSuite) TestCreateTask() {
	ctx := context.Background()
	candidates := []struct {
		input          Task
		expectedResult Task
		expectedCode   codes.Code
	}{
		// valid input
		{
			input: Task{
				CreatedAt:    time.Now().Unix(),
				Name:         "task1",
				Description:  "task1desc",
				UserID:       "5",
				UserEmail:    "example5@tst.com",
				Time:         10,
				TaskID:       "6",
				ReminderSent: false,
			},
			expectedResult: Task{
				CreatedAt:    time.Now().Unix(),
				Name:         "task1",
				Description:  "task1desc",
				UserID:       "5",
				UserEmail:    "example5@tst.com",
				Time:         10,
				TaskID:       "6",
				ReminderSent: false,
			},
			expectedCode: codes.OK,
		},
		// invalid input - no data from context
		//{
		//	input: Task{
		//		CreatedAt:    0,
		//		Name:         "",
		//		Description:  "",
		//		UserID:       "",
		//		UserEmail:    "",
		//		Time:         0,
		//		TaskID:       "",
		//		ReminderSent: false,
		//	},
		//	expectedResult: Task{},
		//	expectedCode:   codes.InvalidArgument,
		//},
		// invalid email -> email validation takes part in the User MS
		//-> userID, userEmail is passed through context
		{
			input: Task{
				CreatedAt:    time.Now().Unix(),
				Name:         "wrong_email",
				Description:  "task with invalid email",
				UserID:       "4",
				UserEmail:    "@@@gamil.cz",
				Time:         time.Now().Add(time.Minute * 60).Unix(),
				TaskID:       "55",
				ReminderSent: false,
			},
			expectedResult: Task{
				CreatedAt:    time.Now().Unix(),
				Name:         "wrong_email",
				Description:  "task with invalid email",
				UserID:       "4",
				UserEmail:    "@@@gamil.cz",
				Time:         time.Now().Add(time.Minute * 60).Unix(),
				TaskID:       "55",
				ReminderSent: false,
			},
			expectedCode: codes.OK,
		},
	}
	for i, candidate := range candidates {
		task, err := s.taskRepo.Create(ctx, candidate.input)
		s.NotNilf(task.TaskID, "candidate %d", i+1)
		task.TaskID = candidate.input.TaskID
		s.Equalf(candidate.expectedResult, task, "candidate %d", i+1)
		s.Equalf(candidate.expectedCode, status.Code(err), "candidate %d", i+1)
	}
}

func (s *RepoTaskTestSuite) TestUpdateTask() {
	ctx := context.Background()
	candidates := []struct {
		userID         string
		taskID         string
		input          Task
		expectedResult Task
		expectedCode   codes.Code
	}{
		// valid input
		{
			userID: "1",
			taskID: "tid1",
			input: Task{
				CreatedAt:    5,
				Name:         "newName",
				Description:  "newDesc",
				UserID:       "1",
				UserEmail:    "example1@tst.com",
				Time:         11,
				TaskID:       "tid1",
				ReminderSent: false,
			},
			expectedResult: Task{
				CreatedAt:    5,
				Name:         "newName",
				Description:  "newDesc",
				UserID:       "1",
				UserEmail:    "example1@tst.com",
				Time:         11,
				TaskID:       "tid1",
				ReminderSent: false,
			},
			expectedCode: codes.OK,
		},
		// wrong taskID -> new task
		{
			userID: "1",
			taskID: "tid777",
			input: Task{
				CreatedAt:    5,
				Name:         "newName",
				Description:  "newDesc",
				UserID:       "1",
				UserEmail:    "example1@tst.com",
				Time:         11,
				TaskID:       "tid777",
				ReminderSent: false,
			},
			expectedResult: Task{
				CreatedAt:    5,
				Name:         "newName",
				Description:  "newDesc",
				UserID:       "1",
				UserEmail:    "example1@tst.com",
				Time:         11,
				TaskID:       "tid777",
				ReminderSent: false,
			},
			expectedCode: codes.OK,
		},
	}
	for i, candidate := range candidates {
		task, err := s.taskRepo.Update(ctx, candidate.input, candidate.userID, candidate.taskID)
		task.CreatedAt = candidate.input.CreatedAt
		s.Equalf(candidate.expectedResult, task, "candidate %d", i+1)
		s.Equalf(candidate.expectedCode, status.Code(err), "candidate %d", i+1)
	}
}

func (s *RepoTaskTestSuite) TestDeleteTask() {
	ctx := context.Background()
	candidates := []struct {
		userID       string
		taskID       string
		expectedCode codes.Code
	}{
		//valid input
		{
			userID:       "1",
			taskID:       "tid1",
			expectedCode: codes.OK,
		},
		//non-existing task
		{
			userID:       "1",
			taskID:       "tid9999",
			expectedCode: codes.OK,
		},
		//non-existing user
		{
			userID:       "9999",
			taskID:       "tid1",
			expectedCode: codes.OK,
		},
	}

	for i, candidate := range candidates {
		err := s.taskRepo.Delete(ctx, candidate.userID, candidate.taskID)
		s.NoError(err)
		_, err = s.taskRepo.Get(ctx, candidate.userID, candidate.taskID)
		// check if deleted correctly
		s.Equalf(codes.NotFound, status.Code(err), "candidate %d", i+1)
	}
}

func (s *RepoTaskTestSuite) TestGetLastNTasks() {
	ctx := context.Background()
	candidates := []struct {
		userID         string
		n              int32
		expectedResult []Task
		expectedCode   codes.Code
	}{
		//valid input
		{
			userID: "1",
			n:      2,
			expectedResult: []Task{
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
					CreatedAt:    2,
					Name:         "task2",
					Description:  "desc2",
					UserID:       "1",
					UserEmail:    "example1@tst.com",
					Time:         time.Now().Add(time.Minute * 1).Unix(),
					TaskID:       "tid2",
					ReminderSent: false,
				},
			},
			expectedCode: codes.OK,
		},
		// more n than tasks
		{
			userID: "2",
			n:      3,
			expectedResult: []Task{
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
			},
			expectedCode: codes.OK,
		},
	}
	for i, candidate := range candidates {
		tasks, err := s.taskRepo.GetLastN(ctx, candidate.userID, candidate.n)
		s.Equalf(candidate.expectedResult, tasks, "candidate %d", i+1)

		s.Equalf(candidate.expectedCode, status.Code(err), "candidate %d", i+1)
	}
}

func (s *RepoTaskTestSuite) TestGetExpiredTasks() {
	ctx := context.Background()
	candidates := []struct {
		userID         string
		expectedResult []Task
		expectedCode   codes.Code
	}{
		//valid input
		{
			userID: "1",
			expectedResult: []Task{
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
					CreatedAt:    3,
					Name:         "task3",
					Description:  "desc3",
					UserID:       "1",
					UserEmail:    "example1@tst.com",
					Time:         10,
					TaskID:       "tid3",
					ReminderSent: false,
				},
			},
			expectedCode: codes.OK,
		},
		// non existent user
		{
			userID:         "999",
			expectedResult: nil,
			expectedCode:   codes.OK,
		},
	}
	for i, candidate := range candidates {
		tasks, err := s.taskRepo.GetExpired(ctx, candidate.userID)
		s.Equalf(candidate.expectedResult, tasks, "candidate %d", i+1)
		s.Equalf(candidate.expectedCode, status.Code(err), "candidate %d", i+1)
	}
}

func (s *RepoTaskTestSuite) TestSearchForExpiringTasks() {
	ctx := context.Background()
	expectedResult := make(map[string][]Task)
	expectedResult["example1@tst.com"] = []Task{
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
	}
	expectedResult["example6@tst.com"] = []Task{
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
	}
	expectedResult["example7@tst.com"] = []Task{
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
	tasks, err := s.taskRepo.SearchForExpiringTasks(ctx)
	s.NoError(err)
	s.Equalf(expectedResult, tasks, "ok")
}

func TestTaskRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(RepoTaskTestSuite))
}
