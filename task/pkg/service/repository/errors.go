package repository

import (
	"errors"
)

var (
	ErrInvalidTime   = errors.New("invalid time provided, task can not be set in the past")
	ErrNameMaxLength = errors.New("maximum length of the task name is 50 characters")
	ErrDescMaxLength = errors.New("maximum length of the task description is 300 characters")
)
