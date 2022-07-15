package service

import (
	"github.com/stretchr/testify/mock"
)

type ClientMock struct {
	mock.Mock
}

func NewClientMock() *ClientMock {
	return &ClientMock{}
}

func (m *ClientMock) Send(to []string, message []byte) error {
	args := m.Called(to, message)
	return args.Error(0)
}
