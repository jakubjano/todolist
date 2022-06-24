package repository

import (
	"cloud.google.com/go/firestore"
	"context"
	"log"
)

//TODO test on emulated FS dtb
// dockerize emulator via docker-compose.yml

type FSUserInterface interface {
	Get(ctx context.Context, UserId string) (User, error)
	Update(ctx context.Context, UserId string, user User) (User, error)
	Delete(ctx context.Context, UserId string) error
}

type FSUser struct {
	fs *firestore.CollectionRef
}

func NewFSUser(fs *firestore.CollectionRef) *FSUser {
	return &FSUser{
		fs: fs,
	}
}

func (a *FSUser) Get(ctx context.Context, UserId string) (User, error) {
	doc, err := a.fs.Doc(UserId).Get(ctx)
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

func (a *FSUser) Update(ctx context.Context, UserId string, user User) (User, error) {
	_, err := a.fs.Doc(UserId).Set(ctx, user)
	if err != nil {
		return User{}, err
	}
	return user, nil
}

func (a *FSUser) Delete(ctx context.Context, UserId string) error {
	_, err := a.fs.Doc(UserId).Delete(ctx)
	if err != nil {
		log.Printf("Error deleting user with id %s", UserId)
		return err
	}
	return nil
}
