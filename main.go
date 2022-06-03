package main

import (
	"context"
	firebase "firebase.google.com/go"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	v1 "github.com/jakubjano/todolist/apis/go-sdk/user/v1"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/api/option"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	service2 "jakubjano/todolist/user/pkg/service"
	"jakubjano/todolist/user/pkg/service/repository"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
)

func main() {

	// Use the application default credentials
	ctx := context.Background()
	key := option.WithCredentialsFile("todolist-dd92e-firebase-adminsdk-9ase9-a45f00f8a4.json")

	lis, err := net.Listen("tcp", ":8081")
	if err != nil {
		panic(err)
	}
	app, err := firebase.NewApp(ctx, nil, key)
	if err != nil {
		panic(err)
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		panic(err)
	}
	defer client.Close()

	userRepo := repository.NewFSUser(client)
	service := service2.NewUserService(userRepo)

	grpcServer := grpc.NewServer()
	v1.RegisterUserServiceServer(grpcServer, service)
	reflection.Register(grpcServer)
	mux := runtime.NewServeMux()
	err = v1.RegisterUserServiceHandlerServer(ctx, mux, service)
	if err != nil {
		panic(err)
	}
	srv := &http.Server{
		Addr:    "8081",
		Handler: grpcHandlerFunc(grpcServer),
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		fmt.Println("Shutting down server")
		grpcServer.GracefulStop()
		_ = srv.Close()
	}()

	err = srv.Serve(lis)
	if err != nil {
		panic(err)
	}

}

func grpcHandlerFunc(grpcServer *grpc.Server) http.Handler {
	return h2c.NewHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.ProtoMajor == 2 && strings.Contains(r.Header.Get("Content-Type"), "application/grpc") {
			grpcServer.ServeHTTP(w, r)
		}
	}), &http2.Server{})
}
