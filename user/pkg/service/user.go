package service

import (
	"context"
	"fmt"
	v1 "github.com/jakubjano/todolist/apis/go-sdk/user/v1"
	"jakubjano/todolist/user/pkg/service/repository"
)

type UserService struct {
	v1.UnimplementedUserServiceServer // from proto, must be present
	userRepo                          repository.FSUserInterface
}

//func createClient(ctx context.Context) *firestore.Client {
//	// Sets your Google Cloud Platform project ID.
//	projectID := "todolist-dd92e"
//	client, err := firestore.NewClient(ctx, projectID)
//	if err != nil {
//		panic(err)
//	}
//	// Close client when done with
//	// defer client.Close()
//	return client
//}

//TODO this is wrong, don't want to create firebase instance with every func call ()

func NewUserService(userRepo repository.FSUserInterface) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

func (s *UserService) UpdateUser(ctx context.Context, in *v1.User) (*v1.User, error) {

	// Get old user
	// Update user
	// return user

	return in, nil
}

func (s *UserService) GetUser(ctx context.Context, in *v1.GetUserRequest) (*v1.User, error) {

	v := s.userRepo
	fmt.Println(v)

	return nil, nil
}
