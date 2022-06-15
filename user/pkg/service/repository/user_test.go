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
)

type RepoUserTestSuite struct {
	suite.Suite
	client   *firestore.Client
	userRepo FSUserInterface
}

// runs once at the beginning
func (s *RepoUserTestSuite) SetupSuite() {
	//initialize FS client
	//initialize user collection
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, os.Getenv("PROJECT_ID"))
	if err != nil {
		panic(err)
	}
	userRepo := NewFSUser(client.Collection("User"))
	s.userRepo = userRepo
	s.client = client
}

// runs before every test
func (s *RepoUserTestSuite) SetupTest() {
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
			Email:     "",
			FirstName: "",
			LastName:  "",
			Phone:     "",
			Address:   "",
		},

		{
			UserID:    "four",
			Email:     "example4@tst.com",
			FirstName: "fn4",
			LastName:  "ln4",
			Phone:     "p4",
			Address:   "a4",
		},

		{
			UserID:    "5",
			Email:     "",
			FirstName: "fn5",
			LastName:  "ln5",
			Phone:     "p5",
			Address:   "a5",
		},
	}

	for _, user := range users {
		docRef := s.client.Collection("User").Doc(user.UserID)
		batchCreate.Set(docRef, user)
	}
	_, err := batchCreate.Commit(ctx)
	s.NoError(err)
}

func (s *RepoUserTestSuite) TearDownTest() {
	// clear all data from DB after every test
	ctx := context.Background()
	docs, err := s.client.Collection("User").Documents(ctx).GetAll()
	s.NoError(err)
	batch := s.client.Batch()
	for _, doc := range docs {
		batch.Delete(doc.Ref)
	}
	_, err = batch.Commit(ctx)
	s.NoError(err)

}

func (s *RepoUserTestSuite) TearDownSuite() {
	err := s.client.Close()
	s.NoError(err)

}

func (s *RepoUserTestSuite) TestGetUser() {
	ctx := context.Background()
	candidates := []struct {
		UserID         string
		ExpectedResult User
		ExpectedError  error
		ExpectedCode   codes.Code
	}{
		// candidate 1: valid input
		{
			UserID: "1",
			ExpectedResult: User{
				UserID:    "1",
				Email:     "example1@tst.com",
				FirstName: "fn1",
				LastName:  "ln1",
				Phone:     "p1",
				Address:   "a1",
			},
			ExpectedError: nil,
		},
		// candidate 2: valid input
		{
			UserID: "5",
			ExpectedResult: User{
				UserID:    "5",
				Email:     "",
				FirstName: "fn5",
				LastName:  "ln5",
				Phone:     "p5",
				Address:   "a5",
			},
		},
		// candidate 3: valid input
		{
			UserID:         "3",
			ExpectedResult: User{UserID: "3"},
			ExpectedError:  nil,
		},
		// candidate 4: doc not found
		{
			UserID:         "999",
			ExpectedResult: User{},
			ExpectedError: status.Error(codes.NotFound,
				"\"projects/dummy-project-id/databases/(default)/documents/User/999\" not found"),
			ExpectedCode: codes.NotFound,
		},
		// candidate 4: invalid input
		{
			UserID:         "",
			ExpectedResult: User{},
			ExpectedError: status.Error(codes.InvalidArgument,
				"Document name \"projects/dummy-project-id/databases/(default)/documents/User/\" has invalid trailing \"/\"."),
			ExpectedCode: codes.InvalidArgument,
		},
	}

	for i, candidate := range candidates {
		user, err := s.userRepo.Get(ctx, candidate.UserID)
		s.Equalf(candidate.ExpectedResult, user, "candidate %d", i+1)
		s.Equalf(candidate.ExpectedError, err, "candidate %d", i+1)
		s.Equalf(candidate.ExpectedCode, status.Code(err), "candidate %d", i+1)
	}

}

func (s *RepoUserTestSuite) TestUpdateUser() {
	ctx := context.Background()
	candidates := []struct {
		UserID         string
		ExpectedResult User
		ExpectedError  error
	}{
		{
			UserID: "1",
			ExpectedResult: User{
				UserID:    "1",
				Email:     "example1@tst.com",
				FirstName: "fn1",
				LastName:  "ln1",
				Phone:     "p1",
				Address:   "a1",
			},
			ExpectedError: nil,
		},

		{UserID: "999",
			ExpectedResult: User{},
			ExpectedError:  nil,
		},
	}

	for _, candidate := range candidates {
		user, err := s.userRepo.Update(ctx, candidate.UserID, candidate.ExpectedResult)
		s.Equal(candidate.ExpectedResult, user)
		s.Equal(candidate.ExpectedError, err)
	}
}

func (s *RepoUserTestSuite) TestDeleteUser() {
	ctx := context.Background()
	candidates := []struct {
		UserID        string
		ExpectedError error
	}{
		{
			UserID:        "1",
			ExpectedError: nil,
		},
		{
			UserID:        "999",
			ExpectedError: nil,
		},
	}
	for _, candidate := range candidates {
		err := s.userRepo.Delete(ctx, candidate.UserID)
		s.Equal(candidate.ExpectedError, err)
	}
}

func TestUserRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(RepoUserTestSuite))
}
