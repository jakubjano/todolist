package main

import (
	"context"
	firebase "firebase.google.com/go"
	"fmt"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	v1 "github.com/jakubjano/todolist/apis/go-sdk/task/v1"
	"github.com/spf13/viper"
	"google.golang.org/api/option"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
	"jakubjano/todolist/task/pkg/service"
	"jakubjano/todolist/task/pkg/service/repository"
	"log"
	"net"
	"net/http"
)

func main() {
	//todo
	// logger task service

	viper.SetDefault("gateway.port", ":8181")
	viper.SetDefault("http.address", ":8180")
	viper.SetDefault("secret.path", "secret/todolist-dd92e-firebase-adminsdk-9ase9-b03dcda63f.json")

	//todo for future config files - can't panic here because config doesn't exist yet
	viper.AddConfigPath("$HOME/.appname")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		log.Printf("error config not found: %v \n", err)
	}

	gwPort := viper.GetString("gateway.port")
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

	taskRepo := repository.NewFSTask(client.Collection(repository.COLLECTION_TASKS))
	taskService := service.NewTaskService(authClient, taskRepo)

	lis, err := net.Listen("tcp", gwPort)
	if err != nil {
		panic(err)
	}

	s := grpc.NewServer(
		grpc_middleware.WithUnaryServerChain(
			grpc_recovery.UnaryServerInterceptor(),
			//tokenClient.CustomUnaryInterceptor(),
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
		gwPort,
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
	fmt.Printf("starting http server at '%s'\n", viper.GetString("http.address"))

	err = http.ListenAndServe(viper.GetString("http.address"), mux)
	if err != nil {
		panic(err)
	}

}
