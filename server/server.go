package server

import (
	"context"
	"fmt"

	"github.com/Saser/strecku/auth"
	streckuv1 "github.com/Saser/strecku/saser/strecku/v1"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const (
	users                = "users"
	authenticatedUserKey = "authenticatedUser"
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

func (s *Server) authenticatedUser(ctx context.Context) (*streckuv1.User, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "Missing metadata.")
	}
	authorization := md["authorization"]
	if len(authorization) == 0 {
		return nil, status.Error(codes.Unauthenticated, `Missing "authorization" header.`)
	}
	b, err := auth.ParseBasic(authorization[0])
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, `Invalid "authorization" header.`)
	}
	key, ok := s.userKeys[b.Username]
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "Username and/or password incorrect.")
	}
	entry := s.users[s.userIndices[key]]
	if err := bcrypt.CompareHashAndPassword(entry.hash, []byte(b.Password)); err != nil {
		return nil, status.Error(codes.Unauthenticated, "Username and/or password incorrect.")
	}
	return entry.user, nil
}

func (s *Server) GetUser(ctx context.Context, req *streckuv1.GetUserRequest) (*streckuv1.User, error) {
	if req.Name == "" {
		return nil, status.Error(codes.InvalidArgument, "Name is required.")
	}
	au, err := s.authenticatedUser(ctx)
	if err != nil {
		return nil, err
	}
	if req.Name == "users/me" {
		return au, nil
	}
	if !au.Superuser {
		if req.Name == au.Name {
			return au, nil
		}
		return nil, status.Error(codes.PermissionDenied, "Permission denied.")
	}
	index, ok := s.userIndices[req.Name]
	if !ok {
		return nil, status.Error(codes.NotFound, "User not found.")
	}
	return s.users[index].user, nil
}

func (s *Server) CreateUser(ctx context.Context, req *streckuv1.CreateUserRequest) (*streckuv1.User, error) {
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
	au, err := s.authenticatedUser(ctx)
	if err != nil {
		return nil, err
	}
	if !au.Superuser {
		return nil, status.Error(codes.PermissionDenied, "Permission denied.")
	}
	if _, ok := s.userKeys[user.EmailAddress]; ok {
		return nil, status.Error(codes.AlreadyExists, "Email address must be unique.")
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
