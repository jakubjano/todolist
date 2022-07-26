package main

import (
	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"context"
	firebase "firebase.google.com/go"
	"fmt"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	v1 "github.com/jakubjano/todolist/apis/go-sdk/task/v1"
	"github.com/jakubjano/todolist/task/internal/auth"
	"github.com/jakubjano/todolist/task/pkg/service"
	"github.com/jakubjano/todolist/task/pkg/service/repository"
	"github.com/spf13/viper"
	cloudscheduler "google.golang.org/api/cloudscheduler/v1beta1"
	"google.golang.org/api/option"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
	"net"
	"net/http"
)

func main() {

	viper.SetDefault("grpc.port", ":8181")
	viper.SetDefault("gateway.port", ":8180")
	viper.SetDefault("firebase.secret", "projects/todolist-356712/secrets/firebase-key/versions/latest")
	viper.SetDefault("host", "smtp.mailtrap.io")
	viper.SetDefault("from", "jakubjanek8@gmail.com")
	viper.SetDefault("email.credentials", "projects/todolist-356712/secrets/email-credentials/versions/latest")
	viper.SetDefault("scheduler.path", "projects/todolist-356712/locations/europe-west3")

	ctx := context.Background()
	logger, err := service.NewLogger()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	secretClient, err := secretmanager.NewClient(ctx)
	if err != nil {
		panic(err)
	}
	defer secretClient.Close()
	secretManager := service.NewSecretManager(ctx, secretClient)
	firebaseSecret, err := secretManager.AccessSecret(
		viper.GetString("firebase.secret"))
	if err != nil {
		panic(err)
	}
	key := option.WithCredentialsJSON(firebaseSecret)

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
	schedulerService, err := cloudscheduler.NewService(ctx)
	if err != nil {
		panic(err)
	}
	scheduler := service.NewSchedulerService(schedulerService, viper.GetString("scheduler.path"))

	taskRepo := repository.NewFSTask(client.Collection(repository.CollectionUsers), client)
	credJsonData, err := secretManager.AccessSecret(
		viper.GetString("email.credentials"))
	if err != nil {
		panic(err)
	}
	emailCredentials, err := secretManager.MapSecretData(credJsonData)
	if err != nil {
		panic(err)
	}
	settings := &service.Settings{
		Host:     viper.GetString("host"),
		From:     viper.GetString("from"),
		UserName: emailCredentials.Username,
		Password: emailCredentials.Password,
	}
	emailSender := service.NewEmailSender(settings)
	reminder := service.NewReminder(taskRepo, logger, emailSender, client)
	taskService := service.NewTaskService(taskRepo, logger, reminder, scheduler)
	tokenClient := auth.NewTokenClient(authClient, logger)

	grpcPort := viper.GetString("grpc.port")
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
	v1.RegisterTaskServiceServer(s, taskService)
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
	err = v1.RegisterTaskServiceHandler(context.Background(), mux, conn)
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
