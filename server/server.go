package server

import (
	"context"
	"fmt"

	streckuv1 "github.com/Saser/strecku/saser/strecku/v1"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	users = "users"
)

type userEntry struct {
	user *streckuv1.User
	hash []byte
}

type Server struct {
	streckuv1.UnimplementedStreckUServer

	users       []*userEntry
	userIndices map[string]int    // name -> index into users
	userKeys    map[string]string // email address -> name
}

func New() *Server {
	return &Server{
		userIndices: make(map[string]int),
		userKeys:    make(map[string]string),
	}
}

func newUserName() string {
	return fmt.Sprintf("%s/%s", users, uuid.New())
}

func (s *Server) CreateUser(_ context.Context, req *streckuv1.CreateUserRequest) (*streckuv1.User, error) {
	user := req.User
	if user.EmailAddress == "" {
		return nil, status.Error(codes.InvalidArgument, "Email address is required.")
	}
	if user.DisplayName == "" {
		return nil, status.Error(codes.InvalidArgument, "Display name is required.")
	}
	if req.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "Password is required.")
	}
	user.Name = newUserName()
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, status.Error(codes.Internal, "Internal error.")
	}
	s.users = append(s.users, &userEntry{
		user: user,
		hash: hash,
	})
	s.userIndices[user.Name] = len(s.users) - 1
	s.userKeys[user.EmailAddress] = user.Name
	return user, nil
}
