package inmemory

import (
	"context"
	"fmt"
	"sync"

	streckuv1 "github.com/Saser/strecku/backend/gen/api/v1"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserAPI struct {
	mu    sync.Mutex
	users map[string]*streckuv1.User
}

func NewUserAPI() *UserAPI {
	return &UserAPI{
		users: make(map[string]*streckuv1.User),
	}
}

func (u *UserAPI) GetUser(_ context.Context, req *streckuv1.GetUserRequest) (*streckuv1.GetUserResponse, error) {
	u.mu.Lock()
	defer u.mu.Unlock()
	user, ok := u.users[req.Name]
	if !ok {
		return nil, status.Error(codes.NotFound, "User resource not found")
	}
	return &streckuv1.GetUserResponse{
		User: user,
	}, nil
}

func (u *UserAPI) ListUsers(context.Context, *streckuv1.ListUsersRequest) (*streckuv1.ListUsersResponse, error) {
	u.mu.Lock()
	defer u.mu.Unlock()
	users := make([]*streckuv1.User, 0, len(u.users))
	for _, user := range u.users {
		users = append(users, user)
	}
	return &streckuv1.ListUsersResponse{
		Users: users,
	}, nil
}

func (u *UserAPI) CreateUser(_ context.Context, req *streckuv1.CreateUserRequest) (*streckuv1.CreateUserResponse, error) {
	u.mu.Lock()
	defer u.mu.Unlock()
	newUser := req.User
	if newUser.EmailAddress == "" {
		return nil, status.Error(codes.InvalidArgument, "empty email address")
	}
	if newUser.DisplayName == "" {
		return nil, status.Error(codes.InvalidArgument, "empty display name")
	}
	for _, user := range u.users {
		if user.EmailAddress == newUser.EmailAddress {
			return nil, status.Error(codes.AlreadyExists, "email address already exists")
		}
	}
	id := uuid.New()
	newUser.Name = fmt.Sprintf("users/%s", id.String())
	u.users[newUser.Name] = newUser
	return &streckuv1.CreateUserResponse{
		User: newUser,
	}, nil
}
