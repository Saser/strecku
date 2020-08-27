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

func (s *Server) AuthenticateUser(_ context.Context, req *streckuv1.AuthenticateUserRequest) (*streckuv1.User, error) {
	if req.EmailAddress == "" {
		return nil, status.Error(codes.InvalidArgument, "Email address is required.")
	}
	if req.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "Password is required.")
	}
	name, ok := s.userKeys[req.EmailAddress]
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "Authentication failed.")
	}
	entry := s.users[s.userIndices[name]]
	if err := bcrypt.CompareHashAndPassword(entry.hash, []byte(req.Password)); err != nil {
		return nil, status.Error(codes.Unauthenticated, "Authentication failed.")
	}
	return entry.user, nil
}

func newUserName() string {
	return fmt.Sprintf("%s/%s", users, uuid.New())
}
