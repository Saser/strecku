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
	users  = "users"
	stores = "stores"

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

	stores       []*streckuv1.Store
	storeIndices map[string]int // name -> index into stores
}

func New() *Server {
	return &Server{
		userIndices: make(map[string]int),
		userKeys:    make(map[string]string),

		storeIndices: make(map[string]int),
	}
}

func newUserName() string {
	return fmt.Sprintf("%s/%s", users, uuid.New())
}

func newStoreName() string {
	return fmt.Sprintf("%s/%s", stores, uuid.New())
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

func (s *Server) ListUsers(ctx context.Context, req *streckuv1.ListUsersRequest) (*streckuv1.ListUsersResponse, error) {
	if req.PageSize < 0 {
		return nil, status.Error(codes.InvalidArgument, "Page size must be non-negative.")
	}
	if req.PageSize > 0 {
		return nil, status.Error(codes.Unimplemented, "Pagination is not implemented.")
	}
	if req.PageToken != "" {
		return nil, status.Error(codes.Unimplemented, "Pagination is not implemented.")
	}
	au, err := s.authenticatedUser(ctx)
	if err != nil {
		return nil, err
	}
	if !au.Superuser {
		return &streckuv1.ListUsersResponse{
			Users:         []*streckuv1.User{au},
			NextPageToken: "",
		}, nil
	}
	users := make([]*streckuv1.User, len(s.users))
	for i, entry := range s.users {
		users[i] = entry.user
	}
	return &streckuv1.ListUsersResponse{Users: users}, nil
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

func (s *Server) CreateStore(ctx context.Context, req *streckuv1.CreateStoreRequest) (*streckuv1.Store, error) {
	store := req.Store
	if store.DisplayName == "" {
		return nil, status.Error(codes.InvalidArgument, "Display name is required.")
	}
	store.Name = newStoreName()
	s.stores = append(s.stores, store)
	s.storeIndices[store.Name] = len(s.stores) - 1
	return store, nil
}
