package repository

import "time"

type Task struct {
	CreatedAt   time.Time `firestorm:"CreatedAt"`
	Name        string    `firestorm:"Name"`
	Description string    `firestorm:"Description"`
	UserID      string    `firestorm:"UserID"`
	Time        time.Time `firestorm:"Time"`
}
