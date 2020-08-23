package server

import (
	"context"
	"fmt"

	streckuv1 "github.com/Saser/strecku/saser/strecku/v1"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	users = "users"
)

type Server struct {
	streckuv1.UnimplementedStreckUServer

	users     []*streckuv1.User
	userKeys  map[string]int
	passwords map[string]string
}

func New() *Server {
	return &Server{
		userKeys:  make(map[string]int),
		passwords: make(map[string]string),
	}
}

func (s *Server) AuthenticateUser(_ context.Context, req *streckuv1.AuthenticateUserRequest) (*streckuv1.User, error) {
	for _, user := range s.users {
		if user.EmailAddress == req.EmailAddress && s.passwords[user.Name] == req.Password {
			return user, nil
		}
	}
	return nil, status.Error(codes.Unauthenticated, "Authentication failed.")
}

func newUserName() string {
	return fmt.Sprintf("%s/%s", users, uuid.New())
}
