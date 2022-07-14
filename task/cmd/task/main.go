package main

import (
	"context"
	firebase "firebase.google.com/go"
	"fmt"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	v1 "github.com/jakubjano/todolist/apis/go-sdk/task/v1"
	"github.com/robfig/cron/v3"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"google.golang.org/api/option"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
	"jakubjano/todolist/task/internal/auth"
	"jakubjano/todolist/task/pkg/service"
	"jakubjano/todolist/task/pkg/service/repository"
	"net"
	"net/http"
	"net/smtp"
)

func main() {

	logger, err := service.NewLogger()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	viper.SetConfigName("task_config")
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

	taskRepo := repository.NewFSTask(client.Collection(repository.CollectionUsers), client)
	taskService := service.NewTaskService(taskRepo, logger)
	tokenClient := auth.NewTokenClient(authClient)

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

	// cron reminders
	emailAuth := smtp.PlainAuth("", viper.GetString("username"), viper.GetString("password"),
		viper.GetString("host"))
	reminder := service.NewReminder(taskRepo, logger, emailAuth, client)
	c := cron.New()
	c.AddFunc("@every 30s", func() {
		err := reminder.RemindUserViaEmail(ctx, viper.GetString("host"),
			viper.GetString("smtp_port"),
			viper.GetString("from"))
		if err != nil {
			panic(err)
		}
	},
	)
	c.Start()

	// Start HTTP server (and proxy calls to gRPC server endpoint)
	fmt.Printf("starting http server at '%s'\n", viper.GetString("gateway.port"))
	err = http.ListenAndServe(viper.GetString("gateway.port"), mux)
	if err != nil {
		panic(err)
	}
}
