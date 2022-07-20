package service

import (
	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"context"
	secretmanagerpb "google.golang.org/genproto/googleapis/cloud/secretmanager/v1"
)

type SecretManager struct {
	ctx    context.Context
	client *secretmanager.Client
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
