package repository

import (
	"cloud.google.com/go/firestore"
	"context"
	"log"
)

type FSUserInterface interface {
	Get(ctx context.Context, UserID string) (User, error)
	Update(ctx context.Context, userID string, user User) (User, error)
	Delete(ctx context.Context, UserID string) error
}

type FSUser struct {
	fs *firestore.CollectionRef
}

func NewFSUser(fs *firestore.CollectionRef) *FSUser {
	return &FSUser{
		fs: fs,
	}
}

func (a *FSUser) Get(ctx context.Context, UserID string) (User, error) {
	doc, err := a.fs.Doc(UserID).Get(ctx)
	if err != nil {
		return User{}, err
	}
	user := User{}
	err = doc.DataTo(&user)
	if err != nil {
		return User{}, err
	}
	return user, nil
}

func (a *FSUser) Update(ctx context.Context, userID string, user User) (User, error) {
	_, err := a.fs.Doc(userID).Set(ctx, user)
	if err != nil {
		return User{}, err
	}
	return user, nil
}

func (a *FSUser) Delete(ctx context.Context, UserID string) error {
	_, err := a.fs.Doc(UserID).Delete(ctx)
	if err != nil {
		log.Printf("Error deleting user with id %s", UserID)
		return err
	}
	return nil
}
