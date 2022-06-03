package repository

type User struct {
	UserID    string `firestore:"UserID"`
	Email     string `firestore:"Email"`
	FirstName string `firestore:"FirstName"`
	LastName  string `firestore:"LastName"`
	Phone     string `firestore:"Phone"`
	Address   string `firestore:"Address"`
}
