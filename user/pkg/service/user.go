package service

import (
	"context"
	"firebase.google.com/go/auth"
	"fmt"
	v1 "github.com/jakubjano/todolist/apis/go-sdk/user/v1"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	middleware "jakubjano/todolist/user/internal/auth"
	"jakubjano/todolist/user/pkg/service/repository"
	"log"
	"net/http"
)

type AuthClientInterface interface {
	GetUserByEmail(ctx context.Context, email string) (*auth.UserRecord, error)
	DeleteUser(ctx context.Context, uid string) error
}

type UserService struct {
	v1.UnimplementedUserServiceServer // from proto, must be present
	authClient                        AuthClientInterface
	userRepo                          repository.FSUserInterface
}

func NewUserService(authClient AuthClientInterface, userRepo repository.FSUserInterface) *UserService {
	return &UserService{
		authClient: authClient,
		userRepo:   userRepo,
	}
}

func (s *UserService) UpdateUser(ctx context.Context, in *v1.User) (*v1.User, error) {
	userCtx := ctx.Value("user").(*middleware.UserContext)
	switch userCtx.Role {
	case "admin":
		log.Printf("Admin %s authorized\n", userCtx.Email)
	case "user":
		if userCtx.Email != in.Email {
			return &v1.User{}, status.Error(http.StatusUnauthorized, ErrUnauthorized.Error())
		}
	}
	fbUser, err := s.authClient.GetUserByEmail(ctx, in.Email)
	if err != nil {
		return &v1.User{}, err //status.Error(http.StatusBadRequest, err.Error())
	}
	user, err := s.userRepo.Update(ctx, fbUser.UID, repository.UserFromMsg(in))
	if err != nil {
		return &v1.User{}, status.Error(http.StatusInternalServerError, err.Error())
	}
	return user.ToApi(), nil
}

func (s *UserService) GetUser(ctx context.Context, in *v1.GetUserRequest) (*v1.User, error) {
	userCtx := ctx.Value("user").(*middleware.UserContext)
	switch userCtx.Role {
	case "admin":
		fmt.Println("Admin authorized")
	case "user":
		if userCtx.UserID != in.UserID {
			return &v1.User{}, status.Error(http.StatusUnauthorized, ErrUnauthorized.Error())
		}
	}
	fmt.Println(ctx.Value("user"))
	user, err := s.userRepo.Get(ctx, in.UserID)
	if err != nil {
		return &v1.User{}, status.Error(http.StatusBadRequest, err.Error())
	}
	return user.ToApi(), nil
}

func (s *UserService) DeleteUser(ctx context.Context, in *v1.DeleteUserRequest) (*emptypb.Empty, error) {
	userCtx := ctx.Value("user").(*middleware.UserContext)
	if userCtx.Role != "admin" {
		return &emptypb.Empty{}, status.Error(http.StatusUnauthorized, ErrUnauthorized.Error())
	}
	err := s.authClient.DeleteUser(ctx, in.UserID)
	if err != nil {
		log.Printf("error deleting user: %v\n", err)
		return &emptypb.Empty{}, status.Error(http.StatusBadRequest, err.Error())
	}
	log.Printf("Successfully deleted user: %s\n", in.UserID)
	err = s.userRepo.Delete(ctx, in.UserID)
	if err != nil {
		return &emptypb.Empty{}, status.Error(http.StatusInternalServerError, err.Error())
	}
	return &emptypb.Empty{}, nil
}
