package main

import (
	"context"
	firebase "firebase.google.com/go"
	"fmt"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	v1 "github.com/jakubjano/todolist/apis/go-sdk/user/v1"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"google.golang.org/api/option"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
	"jakubjano/todolist/user/internal/auth"
	"jakubjano/todolist/user/pkg/service"
	"jakubjano/todolist/user/pkg/service/repository"
	"net"
	"net/http"
)

func main() {
	//todo
	// dotfiles $HOME/.config/ for viper and terraform ?

	logger, err := service.NewLogger()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	viper.SetConfigName("user_config")
	viper.SetConfigType("json")
	viper.AddConfigPath("./secret")
	viper.AddConfigPath(".")
	err = viper.ReadInConfig()
	if err != nil {
		logger.Warn("error finding config file, using default values", zap.Error(err))
	}

	grpcPort := viper.GetString("grpc.port")
	ctx := context.Background()
	key := option.WithCredentialsFile(viper.GetString("secret.path"))

	app, err := firebase.NewApp(ctx, nil, key)
	if err != nil {
		panic(err)
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		panic(err)
	}
	defer client.Close()

	authClient, err := app.Auth(ctx)
	if err != nil {
		panic(err)
	}

	userRepo := repository.NewFSUser(client.Collection(repository.CollectionUsers))
	userService := service.NewUserService(authClient, userRepo, logger)
	tokenClient := auth.NewTokenClient(authClient, logger)

	lis, err := net.Listen("tcp", grpcPort)
	if err != nil {
		panic(err)
	}

	s := grpc.NewServer(
		grpc_middleware.WithUnaryServerChain(
			grpc_recovery.UnaryServerInterceptor(),
			tokenClient.CustomUnaryInterceptor(),
		),
	)
	v1.RegisterUserServiceServer(s, userService)
	reflection.Register(s)

	go func() {
		err = s.Serve(lis)
		if err != nil {
			panic(err)
		}
	}()

	conn, err := grpc.DialContext(
		context.Background(),
		grpcPort,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		panic(err)
	}

	// Register gRPC server endpoint
	// Note: Make sure the gRPC server is running properly and accessible
	mux := runtime.NewServeMux()
	err = v1.RegisterUserServiceHandler(context.Background(), mux, conn)
	if err != nil {
		panic(err)
	}

	// Start HTTP server (and proxy calls to gRPC server endpoint)
	fmt.Printf("starting http server at '%s'\n", viper.GetString("gateway.port"))

	err = http.ListenAndServe(viper.GetString("gateway.port"), mux)
	if err != nil {
		panic(err)
	}

}
