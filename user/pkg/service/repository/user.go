package repository

import (
	"cloud.google.com/go/firestore"
	"context"
	"log"
)

type FSUserInterface interface {
	Get(ctx context.Context, UserID string) (User, error)
	Update(ctx context.Context, user User) (User, error)
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
	user := User{}
	doc, err := a.fs.Doc(UserID).Get(ctx)
	if err != nil {
		return User{}, err
	}
	err = doc.DataTo(&user)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

func (a *FSUser) Update(ctx context.Context, user User) (User, error) {

	_, err := a.fs.Doc(user.UserID).Set(ctx, user)
	if err != nil {
		log.Printf("Error updating user with id(%s) on the database layer", user.UserID)
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

//func (a *FSUser) Create(ctx context.Context, user User) (User, error) {
//
//	return User{}, nil
//}
