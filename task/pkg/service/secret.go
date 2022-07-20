package service

import (
	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"context"
	"encoding/json"
	secretmanagerpb "google.golang.org/genproto/googleapis/cloud/secretmanager/v1"
)

type SecretManager struct {
	ctx    context.Context
	client *secretmanager.Client
}

type EmailCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func NewSecretManager(ctx context.Context, client *secretmanager.Client) *SecretManager {
	return &SecretManager{
		ctx:    ctx,
		client: client,
	}
}

func (s *SecretManager) AccessSecret(name string) ([]byte, error) {
	req := &secretmanagerpb.AccessSecretVersionRequest{
		Name: name,
	}
	result, err := s.client.AccessSecretVersion(s.ctx, req)
	if err != nil {
		return nil, err
	}
	return result.Payload.Data, nil
}

func (s *SecretManager) MapSecretData(data []byte) (*EmailCredentials, error) {
	emailCredentials := &EmailCredentials{}
	err := json.Unmarshal(data, emailCredentials)
	if err != nil {
		return nil, err
	}
	return emailCredentials, nil
}
