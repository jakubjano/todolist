package repository

import v1 "github.com/jakubjano/todolist/apis/go-sdk/user/v1"

const (
	CollectionUsers = "users"
)

type User struct {
	UserID    string `firestore:"UserId"`
	Email     string `firestore:"Email"`
	FirstName string `firestore:"FirstName"`
	LastName  string `firestore:"LastName"`
	Phone     string `firestore:"Phone"`
	Address   string `firestore:"Address"`
}

func (u User) ToApi() *v1.User {
	return &v1.User{
		LastName:  u.LastName,
		FirstName: u.FirstName,
		Phone:     u.Phone,
		Address:   u.Address,
		Email:     u.Email,
		UserId:    u.UserID,
	}

}
func UserFromMsg(msg *v1.User) User {
	return User{
		UserID:    msg.UserId,
		Email:     msg.Email,
		FirstName: msg.FirstName,
		LastName:  msg.LastName,
		Phone:     msg.Phone,
		Address:   msg.Address,
	}
}
