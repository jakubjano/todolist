package service

import "errors"

var (
	ErrUnauthorized    = errors.New("unauthorized entry")
	ErrNoExpiringTasks = errors.New("no expiring tasks")
)
