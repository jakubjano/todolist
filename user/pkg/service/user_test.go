package service

import (
	"context"
	"firebase.google.com/go/auth"
	v1 "github.com/jakubjano/todolist/apis/go-sdk/user/v1"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc/status"
	middleware "jakubjano/todolist/user/internal/auth"
	"jakubjano/todolist/user/pkg/service/repository"
	"net/http"
	"testing"
)

type ServiceUserTestSuite struct {
	suite.Suite
	us         *UserService
	mockRepo   *repository.FSUserMock
	mockClient *FBClientMock
}

func (s *ServiceUserTestSuite) SetupSuite() {
	mockRepo := repository.NewMockRepo()
	mockClient := NewFBClientMock()
	us := NewUserService(mockClient, mockRepo)
	s.mockRepo = mockRepo
	s.mockClient = mockClient
	s.us = us
}

func (s *ServiceUserTestSuite) TestGetUser() {
	ctx := context.Background()

	candidates := []struct {
		ctx            context.Context
		in             *v1.GetUserRequest
		ExpectedResult *v1.User
		ExpectedError  error
	}{
		// user role authorized, valid input
		{
			ctx: context.WithValue(ctx, "user", &middleware.UserContext{
				UserID: "id1",
				Email:  "test@example.com",
				Role:   "user",
			}),
			in: &v1.GetUserRequest{UserID: "id1"},
			ExpectedResult: &v1.User{
				LastName:  "anon",
				FirstName: "user",
				Phone:     "09500600",
				Address:   "ad1",
				Email:     "test@example.com",
				UserID:    "id1",
			},
			ExpectedError: nil,
		},

		// user role unauthorized, valid input
		{
			ctx: context.WithValue(ctx, "user", &middleware.UserContext{
				UserID: "id2",
				Email:  "test2@example.com",
				Role:   "user",
			}),
			in:             &v1.GetUserRequest{UserID: "idnot2"},
			ExpectedResult: &v1.User{},
			ExpectedError:  status.Error(http.StatusUnauthorized, "unauthorized entry"),
		},

		//admin role authorized invalid input
		{
			ctx: context.WithValue(ctx, "user", &middleware.UserContext{
				UserID: "idadmin",
				Email:  "jakub@test.com",
				Role:   "admin",
			}),
			in:             &v1.GetUserRequest{UserID: ""},
			ExpectedResult: &v1.User{},
			ExpectedError:  nil,
		},

		//admin role authorized valid input
		{
			ctx: context.WithValue(ctx, "user", &middleware.UserContext{
				UserID: "idadmin",
				Email:  "jakub@test.com",
				Role:   "admin",
			}),
			in: &v1.GetUserRequest{UserID: "id1"},
			ExpectedResult: &v1.User{
				LastName:  "user",
				FirstName: "other",
				Phone:     "09500600",
				Address:   "ad1",
				Email:     "test@example.com",
				UserID:    "id1",
			},
			ExpectedError: nil,
		},
	}

	for i, candidate := range candidates {
		s.mockRepo.On("Get", candidate.ctx, candidate.in.UserID).Return(repository.User{
			UserID:    candidate.ExpectedResult.UserID,
			Email:     candidate.ExpectedResult.Email,
			FirstName: candidate.ExpectedResult.FirstName,
			LastName:  candidate.ExpectedResult.LastName,
			Phone:     candidate.ExpectedResult.Phone,
			Address:   candidate.ExpectedResult.Address,
		}, candidate.ExpectedError)
		user, err := s.us.GetUser(candidate.ctx, candidate.in)
		s.Equalf(candidate.ExpectedResult, user, "candidate %d", i+1)
		s.Equalf(candidate.ExpectedError, err, "candidate %:", i+1)

	}
}

func (s *ServiceUserTestSuite) TestUpdateUser() {
	ctx := context.Background()
	candidates := []struct {
		ctx            context.Context
		in             *v1.User
		ExpectedResult *v1.User
		ExpectedError  error
	}{
		// user role authorized valid input
		{
			ctx: context.WithValue(ctx, "user", &middleware.UserContext{
				UserID: "id1",
				Email:  "user@test.com",
				Role:   "user",
			}),
			in: &v1.User{
				LastName:  "test",
				FirstName: "user",
				Phone:     "123",
				Address:   "a1",
				Email:     "user@test.com",
				UserID:    "id1",
			},
			ExpectedResult: &v1.User{
				LastName:  "test",
				FirstName: "user",
				Phone:     "123",
				Address:   "a1",
				Email:     "user@test.com",
				UserID:    "id1",
			},
			ExpectedError: nil,
		},

		// admin role authorized valid input
		{
			ctx: context.WithValue(ctx, "user", &middleware.UserContext{
				UserID: "id1",
				Email:  "admin@test.com",
				Role:   "admin",
			}),
			in: &v1.User{
				LastName:  "test",
				FirstName: "user",
				Phone:     "123",
				Address:   "a1",
				Email:     "user@test.com",
				UserID:    "id1",
			},
			ExpectedResult: &v1.User{
				LastName:  "test",
				FirstName: "user",
				Phone:     "123",
				Address:   "a1",
				Email:     "user@test.com",
				UserID:    "id1",
			},
			ExpectedError: nil,
		},

		// admin role authorized invalid input
		{
			ctx: context.WithValue(ctx, "user", &middleware.UserContext{
				UserID: "id1",
				Email:  "admin@test.com",
				Role:   "admin",
			}),
			in: &v1.User{
				LastName:  "test",
				FirstName: "user",
				Phone:     "123",
				Address:   "a1",
				Email:     "@@bad_email",
				UserID:    "",
			},
			ExpectedResult: &v1.User{},
			ExpectedError: status.Error(http.StatusBadRequest,
				"malformed email string: @@bad_email"),
		},

		// user role unauthorized valid input
		{
			ctx: context.WithValue(ctx, "user", &middleware.UserContext{
				UserID: "id1",
				Email:  "user@test.com",
				Role:   "user",
			}),
			in: &v1.User{
				LastName:  "test",
				FirstName: "user",
				Phone:     "123",
				Address:   "a1",
				Email:     "notuser@test.com",
				UserID:    "id33",
			},
			ExpectedResult: &v1.User{},
			ExpectedError:  status.Error(http.StatusUnauthorized, "unauthorized entry"),
		},
	}

	for i, candidate := range candidates {
		s.mockClient.On("GetUserByEmail", candidate.ctx, candidate.in.Email).
			Return(&auth.UserRecord{
				UserInfo: &auth.UserInfo{
					UID: candidate.in.UserID,
				},
			}, candidate.ExpectedError)

		s.mockRepo.On("Update", candidate.ctx, candidate.in.UserID, repository.UserFromMsg(candidate.in)).
			Return(repository.User{
				UserID:    candidate.ExpectedResult.UserID,
				Email:     candidate.ExpectedResult.Email,
				FirstName: candidate.ExpectedResult.FirstName,
				LastName:  candidate.ExpectedResult.LastName,
				Phone:     candidate.ExpectedResult.Phone,
				Address:   candidate.ExpectedResult.Address,
			}, candidate.ExpectedError)
		user, err := s.us.UpdateUser(candidate.ctx, candidate.in)
		s.Equalf(candidate.ExpectedResult, user, "candidate %d", i+1)
		s.Equalf(candidate.ExpectedError, err, "candidate %d", i+1)
	}

}

func (s *ServiceUserTestSuite) TestDeleteUser() {
	ctx := context.Background()
	candidates := []struct {
		ctx           context.Context
		in            *v1.DeleteUserRequest
		ExpectedError error
	}{

		// user trying to delete valid input
		{
			ctx: context.WithValue(ctx, "user", &middleware.UserContext{
				UserID: "id1",
				Email:  "test@example.com",
				Role:   "user",
			}),
			in:            &v1.DeleteUserRequest{UserID: "id1"},
			ExpectedError: status.Error(http.StatusUnauthorized, ErrUnauthorized.Error()),
		},

		// admin trying to delete valid input
		{
			ctx: context.WithValue(ctx, "user", &middleware.UserContext{
				UserID: "idadmin",
				Email:  "test@admin.com",
				Role:   "admin",
			}),
			in:            &v1.DeleteUserRequest{UserID: "id1"},
			ExpectedError: nil,
		},

		// admin trying to delete invalid input
		{
			ctx: context.WithValue(ctx, "user", &middleware.UserContext{
				UserID: "idadmin",
				Email:  "test@admin.com",
				Role:   "admin",
			}),
			in:            &v1.DeleteUserRequest{UserID: ""},
			ExpectedError: nil,
		},
	}
	for i, candidate := range candidates {
		s.mockClient.On("DeleteUser", candidate.ctx, candidate.in.UserID).Return(candidate.ExpectedError)
		s.mockRepo.On("Delete", candidate.ctx, candidate.in.UserID).Return(candidate.ExpectedError)
		_, err := s.us.DeleteUser(candidate.ctx, candidate.in)
		s.Equalf(candidate.ExpectedError, err, "candidate %d", i+1)
	}

}

func TestUserServiceTestSuite(t *testing.T) {
	suite.Run(t, new(ServiceUserTestSuite))
}
