package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"sync"

	streckuv1 "github.com/Saser/strecku/backend/gen/api/v1"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserAPIImpl struct {
	mu    sync.Mutex
	users map[string]*streckuv1.User
}

func (u *UserAPIImpl) GetUser(_ context.Context, req *streckuv1.GetUserRequest) (*streckuv1.GetUserResponse, error) {
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

func (u *UserAPIImpl) ListUsers(_ context.Context, req *streckuv1.ListUsersRequest) (*streckuv1.ListUsersResponse, error) {
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

func (u *UserAPIImpl) CreateUser(_ context.Context, req *streckuv1.CreateUserRequest) (*streckuv1.CreateUserResponse, error) {
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

type StoreAPIImpl struct {
	mu     sync.Mutex
	stores map[string]*streckuv1.Store
}

func (s *StoreAPIImpl) GetStore(_ context.Context, req *streckuv1.GetStoreRequest) (*streckuv1.GetStoreResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	store, ok := s.stores[req.Name]
	if !ok {
		return nil, status.Error(codes.NotFound, "Store resource not found")
	}
	return &streckuv1.GetStoreResponse{
		Store: store,
	}, nil
}

func (s *StoreAPIImpl) ListStores(context.Context, *streckuv1.ListStoresRequest) (*streckuv1.ListStoresResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	stores := make([]*streckuv1.Store, 0, len(s.stores))
	for _, store := range s.stores {
		stores = append(stores, store)
	}
	return &streckuv1.ListStoresResponse{
		Stores: stores,
	}, nil
}

func (s *StoreAPIImpl) CreateStore(_ context.Context, req *streckuv1.CreateStoreRequest) (*streckuv1.CreateStoreResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	newStore := req.Store
	newStore.Name = fmt.Sprintf("stores/%s", uuid.New().String())
	s.stores[newStore.Name] = newStore
	return &streckuv1.CreateStoreResponse{
		Store: newStore,
	}, nil
}

func main() {
	userAPIImpl := &UserAPIImpl{
		users: make(map[string]*streckuv1.User),
	}
	storeAPIImpl := &StoreAPIImpl{
		stores: make(map[string]*streckuv1.Store),
	}
	server := grpc.NewServer()
	streckuv1.RegisterUserAPIServer(server, userAPIImpl)
	streckuv1.RegisterStoreAPIServer(server, storeAPIImpl)
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("error: %+v", err)
	}
	if err := server.Serve(listener); err != nil {
		log.Fatalf("error: %+v", err)
	}
	log.Println("goodbye!")
}
