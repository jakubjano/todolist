package auth

import (
	"context"
	"firebase.google.com/go/auth"
	"fmt"
	grpcAuth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"google.golang.org/grpc"
)

const (
	ContextUser = "user"
)

type TokenClient struct {
	authClient *auth.Client
}

func NewTokenClient(authClient *auth.Client) *TokenClient {
	return &TokenClient{
		authClient: authClient,
	}
}

type UserContext struct {
	UserID string
	Email  string
	Role   string
}

// UnaryServerInterceptor returns a new unary server interceptors that performs per-request auth.
func (t *TokenClient) CustomUnaryInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		newCtx, err := t.AuthFunc(ctx)
		if err != nil {
			// token can't be parsed without complete/authorized signature
			// so when error occurs and the request is intercepted there are no data
			// of the user that tried and failed to authorize
			fmt.Printf("error authorizing request: %v", err)
			return nil, err
		}
		return handler(newCtx, req)
	}
}

func (t *TokenClient) AuthFunc(ctx context.Context) (context.Context, error) {
	// AuthFromMD searches for Authorization header from request that is carried by context
	jwt, err := grpcAuth.AuthFromMD(ctx, "bearer")
	if err != nil {
		fmt.Printf("error obtaining auth header from metadata: %v", err)
		return nil, err
	}
	// VerifyIDToken searches for projectID in key automatically when client was initialized with service account
	// credentials
	token, err := t.authClient.VerifyIDToken(ctx, jwt)
	if err != nil {
		fmt.Printf("error verifying token: %v", err)
		return nil, err
	}
	data := token.Claims
	ctxUser := &UserContext{
		UserID: data["user_id"].(string),
		Email:  data["email"].(string),
		Role:   data["role"].(string),
	}
	newCtx := context.WithValue(ctx, ContextUser, ctxUser)
	return newCtx, nil
}
