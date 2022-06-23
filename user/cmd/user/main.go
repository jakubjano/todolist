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

var (
	defaults = map[string]interface{}{
		"gatewayPort": ":8081",
		"httpAddr":    ":8080",
		"secretPath":  "secret/todolist-dd92e-firebase-adminsdk-9ase9-b03dcda63f.json",

		//todo firestore collection ref in config?
	}
)

func main() {
	//todo viper config
	// defaults + config file from env

	//todo
	// dotfiles $HOME/.config/ for viper and terraform ?

	logger, err := service.NewLogger()
	if err != nil {
		logger.Fatal(err.Error())
	}
	defer logger.Sync()

	for k, v := range defaults {
		viper.SetDefault(k, v)
	}
	// for future config files
	viper.AddConfigPath("$HOME/.appname")
	viper.AddConfigPath(".")
	err = viper.ReadInConfig()
	if err != nil {
		logger.Warn("error finding config file, using default values", zap.Error(err))
	}

	gwPort := viper.Get("gatewayPort").(string)
	ctx := context.Background()
	key := option.WithCredentialsFile(viper.Get("secretPath").(string))

	app, err := firebase.NewApp(ctx, nil, key)
	if err != nil {
		logger.Fatal(err.Error())
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		logger.Fatal(err.Error())
	}
	defer client.Close()

	authClient, err := app.Auth(ctx)
	if err != nil {
		logger.Fatal(err.Error())
	}

	userRepo := repository.NewFSUser(client.Collection("users"))
	userService := service.NewUserService(authClient, userRepo, logger)
	tokenClient := auth.NewTokenClient(authClient, logger)

	lis, err := net.Listen("tcp", gwPort)
	if err != nil {
		logger.Fatal(err.Error())
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
			logger.Fatal(err.Error())
		}
	}()

	conn, err := grpc.DialContext(
		context.Background(),
		gwPort,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		logger.Fatal(err.Error())
	}

	// Register gRPC server endpoint
	// Note: Make sure the gRPC server is running properly and accessible
	mux := runtime.NewServeMux()
	err = v1.RegisterUserServiceHandler(context.Background(), mux, conn)
	if err != nil {
		logger.Fatal(err.Error())
	}

	// Start HTTP server (and proxy calls to gRPC server endpoint)
	fmt.Printf("starting http server at '%s'\n", viper.Get("httpAddr").(string))

	err = http.ListenAndServe(viper.Get("httpAddr").(string), mux)
	if err != nil {
		logger.Fatal(err.Error())
	}

}
