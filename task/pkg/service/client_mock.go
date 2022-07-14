package service

import (
	"github.com/stretchr/testify/mock"
	"net/smtp"
)

type ClientMock struct {
	mock.Mock
}

func NewClientMock() *ClientMock {
	return &ClientMock{}
}

func (m *ClientMock) SendMail(addr string, a smtp.Auth, from string, to []string, msg []byte) error {
	args := m.Called(addr, a, from, to, msg)
	return args.Error(0)
}
