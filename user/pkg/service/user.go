package service

import (
	"context"
	"firebase.google.com/go/auth"
	v1 "github.com/jakubjano/todolist/apis/go-sdk/user/v1"
	"go.uber.org/zap"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	middleware "jakubjano/todolist/user/internal/auth"
	"jakubjano/todolist/user/pkg/service/repository"
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
	logger                            *zap.Logger
}

func NewUserService(authClient AuthClientInterface, userRepo repository.FSUserInterface, logger *zap.Logger) *UserService {
	return &UserService{
		authClient: authClient,
		userRepo:   userRepo,
		logger:     logger,
	}
}

func (s *UserService) UpdateUser(ctx context.Context, in *v1.User) (*v1.User, error) {
	userCtx := ctx.Value("user").(*middleware.UserContext)
	log := s.logger.With(
		zap.String("caller_email", userCtx.Email),
		zap.String("caller_id", userCtx.UserId),
		zap.String("caller_role", userCtx.Role),
		zap.String("updated_user_email", in.Email),
	)
	switch userCtx.Role {
	case "admin":
		log.Info("Admin authorized")
	case "user":
		if userCtx.Email != in.Email {
			log.Error(ErrUnauthorized.Error())
			return &v1.User{}, status.Error(http.StatusUnauthorized, ErrUnauthorized.Error())
		}
	}
	fbUser, err := s.authClient.GetUserByEmail(ctx, in.Email)
	if err != nil {
		log.Error(err.Error())
		return &v1.User{}, err //status.Error(http.StatusBadRequest, err.Error())
	}
	user, err := s.userRepo.Update(ctx, fbUser.UID, repository.UserFromMsg(in))
	if err != nil {
		log.Error(err.Error())
		return &v1.User{}, status.Error(http.StatusInternalServerError, err.Error())
	}
	log.Info("update", zap.String("updated_user", user.Email), zap.String("updated_by", userCtx.Email))
	return user.ToApi(), nil
}

func (s *UserService) GetUser(ctx context.Context, in *v1.GetUserRequest) (*v1.User, error) {
	userCtx := ctx.Value("user").(*middleware.UserContext)
	log := s.logger.With(
		zap.String("caller_email", userCtx.Email),
		zap.String("caller_id", userCtx.UserId),
		zap.String("caller_role", userCtx.Role),
		zap.String("get_user_id", in.UserId),
	)
	switch userCtx.Role {
	case "admin":
		log.Info("Admin authorized")
	case "user":
		if userCtx.UserId != in.UserId {
			log.Error(ErrUnauthorized.Error())
			return &v1.User{}, status.Error(http.StatusUnauthorized, ErrUnauthorized.Error())
		}
	}
	user, err := s.userRepo.Get(ctx, in.UserId)
	if err != nil {
		log.Error(err.Error())
		return &v1.User{}, status.Error(http.StatusBadRequest, err.Error())
	}
	return user.ToApi(), nil
}

func (s *UserService) DeleteUser(ctx context.Context, in *v1.DeleteUserRequest) (*emptypb.Empty, error) {
	userCtx := ctx.Value("user").(*middleware.UserContext)
	log := s.logger.With(
		zap.String("caller_email", userCtx.Email),
		zap.String("caller_id", userCtx.UserId),
		zap.String("caller_role", userCtx.Role),
		zap.String("delete_user_id", in.UserId),
	)
	if userCtx.Role != "admin" {
		return &emptypb.Empty{}, status.Error(http.StatusUnauthorized, ErrUnauthorized.Error())
	}
	log.Info("Admin authorized")
	err := s.authClient.DeleteUser(ctx, in.UserId)
	if err != nil {
		log.Error(err.Error())
		return &emptypb.Empty{}, status.Error(http.StatusBadRequest, err.Error())
	}
	log.Info("deleted user from FB")
	err = s.userRepo.Delete(ctx, in.UserId)
	if err != nil {
		log.Error(err.Error())
		return &emptypb.Empty{}, status.Error(http.StatusInternalServerError, err.Error())
	}
	log.Info("deleted user from FS")
	return &emptypb.Empty{}, nil
}
