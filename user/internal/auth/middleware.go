package auth

import (
	"context"
	"firebase.google.com/go/auth"
	"fmt"
	grpcAuth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"google.golang.org/grpc"
)

// notes
// interceptor
// grpc unary interceptors examples

// generate jwt *done
// input for authfunc :token (from request header- see postman article)
// verify signature - by secret google firebase verify token
// ENDPOINT is called

type TokenClient struct {
	authClient *auth.Client
}

func NewTokenClient(authClient *auth.Client) *TokenClient {
	return &TokenClient{
		authClient: authClient}
}

//type AF func(ctx context.Context) (context.Context, error)

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
			return nil, err
		}
		return handler(newCtx, req)
	}
}

func (t *TokenClient) AuthFunc(ctx context.Context) (context.Context, error) {
	//get token form post request's header
	// verify token
	// parse email

	// AuthFromMD searches for Authorization header from request that is carried by context
	jwt, err := grpcAuth.AuthFromMD(ctx, "bearer")
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	// VerifyIDToken searches for projectID in key automatically when client was initialized with service account
	// credentials
	token, err := t.authClient.VerifyIDToken(ctx, jwt)
	if err != nil {
		fmt.Printf("error verifying ID token: %v\n\n", err)
		return nil, err
	}
	//fmt.Printf("Verified ID token: %v\n", token)

	data := token.Claims
	//fmt.Printf("data from token.Claims: %v", data)

	ctxUser := &UserContext{
		UserID: data["user_id"].(string),
		Email:  data["email"].(string),
		Role:   data["role"].(string),
	}

	newCtx := context.WithValue(ctx, "user", ctxUser)
	//fmt.Println(newCtx.Value("user"))

	//TODO
	// parse user email,id from token
	// new type userContext: userid, email, role
	// token.firebase...
	// create new context with parameters from token
	// context.WithValue()

	return newCtx, nil
}

func (t *TokenClient) VerifyUserRole(ctx context.Context) error {
	// to be called in auth function for certain methods on UserService
	// e.g. delete firebase user needs to be authorized by an admin
	userID := ctx.Value("user").(UserContext).UserID
	fmt.Printf("user ID from context :%s\n", userID)
	user, err := t.authClient.GetUser(ctx, userID)
	fmt.Printf("user from authclient:%v\n", user)
	if err != nil {
		fmt.Printf("Error getting user")
		return err
	}

	//// map custom claims to userContext type ?
	//if admin, ok := user.CustomClaims["admin"]; ok {
	//	if admin.(bool) {
	//		return nil
	//	}
	//	err = errors.New("unauthorized request")
	//}
	return err
}
