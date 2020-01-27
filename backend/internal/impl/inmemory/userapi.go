package inmemory

import (
	"context"
	"fmt"

	streckuv1 "github.com/Saser/strecku/backend/gen/api/v1"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (i *Impl) GetUser(_ context.Context, req *streckuv1.GetUserRequest) (*streckuv1.GetUserResponse, error) {
	i.mu.Lock()
	defer i.mu.Unlock()
	user, ok := i.users[req.Name]
	if !ok {
		return nil, status.Error(codes.NotFound, "User resource not found")
	}
	return &streckuv1.GetUserResponse{
		User: user,
	}, nil
}

func (i *Impl) ListUsers(context.Context, *streckuv1.ListUsersRequest) (*streckuv1.ListUsersResponse, error) {
	i.mu.Lock()
	defer i.mu.Unlock()
	users := make([]*streckuv1.User, 0, len(i.users))
	for _, user := range i.users {
		users = append(users, user)
	}
	return &streckuv1.ListUsersResponse{
		Users: users,
	}, nil
}

func (i *Impl) CreateUser(_ context.Context, req *streckuv1.CreateUserRequest) (*streckuv1.CreateUserResponse, error) {
	i.mu.Lock()
	defer i.mu.Unlock()
	newUser := req.User
	if newUser.EmailAddress == "" {
		return nil, status.Error(codes.InvalidArgument, "empty email address")
	}
	if newUser.DisplayName == "" {
		return nil, status.Error(codes.InvalidArgument, "empty display name")
	}
	for _, user := range i.users {
		if user.EmailAddress == newUser.EmailAddress {
			return nil, status.Error(codes.AlreadyExists, "email address already exists")
		}
	}
	id := uuid.New()
	newUser.Name = fmt.Sprintf("users/%s", id.String())
	i.users[newUser.Name] = newUser
	return &streckuv1.CreateUserResponse{
		User: newUser,
	}, nil
}
