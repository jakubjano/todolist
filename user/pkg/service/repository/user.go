package repository

import (
	"cloud.google.com/go/firestore"
	"context"
)

type FSUserInterface interface {
	Get() (User, error)
}

type FSUser struct {
	client *firestore.Client
}

func NewFSUser(client *firestore.Client) *FSUser {
	return &FSUser{
		client: client,
	}
}

func (a *FSUser) Get(userID string) (*User, error) {
	ctx := context.Background()
	doc, err := a.client.Collection("Users").Doc(userID).Get(ctx)
	if err != nil {
		panic(err)
	}
	user := User{}
	err = doc.DataTo(&user)
	if err != nil {
		panic(err)
	}

	return &user, nil
}
