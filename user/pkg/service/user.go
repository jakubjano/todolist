package service

import (
	"context"
	"firebase.google.com/go/auth"
	"fmt"
	v1 "github.com/jakubjano/todolist/apis/go-sdk/user/v1"
	"jakubjano/todolist/user/pkg/service/repository"
	"log"
)

type UserService struct {
	v1.UnimplementedUserServiceServer // from proto, must be present
	authClient                        *auth.Client
	userRepo                          repository.FSUserInterface
}

func NewUserService(authClient *auth.Client, userRepo repository.FSUserInterface) *UserService {
	return &UserService{
		authClient: authClient,
		userRepo:   userRepo,
	}
}

func (s *UserService) UpdateUser(ctx context.Context, in *v1.User) (*v1.User, error) {
	// check the user in firebase auth
	// if user does not exist in firebase auth return error
	_, err := s.authClient.GetUserByEmail(ctx, in.Email)
	if err != nil {
		return &v1.User{}, err
	}
	fmt.Println("User found")
	user, err := s.userRepo.Update(ctx, repository.UserFromMsg(in))
	if err != nil {
		fmt.Println("error updating user on the database layer")
		return &v1.User{}, err
	}
	return user.ToApi(), nil
}

func (s *UserService) GetUser(ctx context.Context, in *v1.GetUserRequest) (*v1.User, error) {
	user, err := s.userRepo.Get(ctx, in.UserID)
	if err != nil {
		log.Printf("error getting user with id:%s", in.UserID)
		return &v1.User{}, err
	}
	return user.ToApi(), nil
}

// no deletion of users on endpoints ?

//func (s *UserService) DeleteUser(ctx context.Context, in *v1.GetUserRequest) error {
//	err := s.authClient.DeleteUser(ctx, in.UserID)
//	if err != nil {
//		log.Printf("error deleting user: %v\n", err)
//		return err
//	}
//	log.Printf("Successfully deleted user: %s\n", in.UserID)
//	err = s.userRepo.Delete(ctx, in.UserID)
//	if err != nil {
//		return err
//	}
//	return nil
//}

// No creation of new FB users on endpoints?

//func (s *UserService) CreateUser(ctx context.Context, in *v1.User) (*v1.User, error) {
//
//	params := (&auth.UserToCreate{}).
//		Email(in.Email).
//		PhoneNumber(in.Phone).
//		DisplayName(in.FirstName + " " + in.LastName)
//	u, err := s.authClient.CreateUser(ctx, params)
//	if err != nil {
//		log.Fatalf("error creating user: %v\n", err)
//	}
//	log.Printf("Successfully created user: %v\n", u)
//
//	//TODO create user in FS db
//
//	return &v1.User{}, nil
//}
