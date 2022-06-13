package main

import (
	"context"
	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"flag"
	"google.golang.org/api/option"
)

func main() {
	key := option.WithCredentialsFile("secret/todolist-dd92e-firebase-adminsdk-9ase9-b03dcda63f.json")
	ctx := context.Background()
	app, err := firebase.NewApp(ctx, nil, key)
	if err != nil {
		panic(err)
	}

	authClient, err := app.Auth(ctx)
	if err != nil {
		panic(err)
	}
	//TODO get input for setUserRole from terminal
	email := flag.String("email", "", "User's email who's role is going to be assigned")
	role := flag.String("role", "", "Role for authorization")
	flag.Parse()
	err = setUserRole(ctx, *email, *role, authClient)
	if err != nil {
		panic(err)
	}
}

func setUserRole(ctx context.Context, email string, role string, authClient *auth.Client) error {
	user, err := authClient.GetUserByEmail(ctx, email)
	if err != nil {
		return err
	}
	claims := map[string]interface{}{"role": role}
	err = authClient.SetCustomUserClaims(ctx, user.UID, claims)
	if err != nil {
		return err
	}
	return nil
}
