package server

import (
	"context"
	"errors"

	"github.com/Saser/strecku/auth"
	"github.com/Saser/strecku/resources/stores"
	"github.com/Saser/strecku/resources/users"
	streckuv1 "github.com/Saser/strecku/saser/strecku/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type Server struct {
	streckuv1.UnimplementedStreckUServer

	userRepo  *users.Repository
	storeRepo *stores.Repository
}

func New(userRepo *users.Repository, storeRepo *stores.Repository) *Server {
	return &Server{
		userRepo:  userRepo,
		storeRepo: storeRepo,
	}
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
	user, err := s.userRepo.LookupUserByEmail(ctx, b.Username)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "Invalid email address and/or password.")
	}
	if err := s.userRepo.Authenticate(ctx, user.Name, b.Password); err != nil {
		return nil, status.Error(codes.Unauthenticated, "Invalid email address and/or password.")
	}
	return user, nil
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
	user, err := s.userRepo.LookupUser(ctx, req.Name)
	if err != nil {
		return nil, status.Error(codes.NotFound, "User not found.")
	}
	return user, nil
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
	allUsers, err := s.userRepo.ListUsers(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, "Internal error.")
	}
	return &streckuv1.ListUsersResponse{
		Users: allUsers,
	}, nil
}

func (s *Server) CreateUser(ctx context.Context, req *streckuv1.CreateUserRequest) (*streckuv1.User, error) {
	user := req.User
	user.Name = users.GenerateName()
	if err := users.Validate(user); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
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
	if err := s.userRepo.CreateUser(ctx, req.User, req.Password); err != nil {
		if exists := new(users.UserExistsError); errors.As(err, &exists) {
			return nil, status.Error(codes.AlreadyExists, "User email address already exists.")
		}
		return nil, status.Error(codes.Internal, "Internal error.")
	}
	return user, nil
}

func (s *Server) GetStore(ctx context.Context, req *streckuv1.GetStoreRequest) (*streckuv1.Store, error) {
	if req.Name == "" {
		return nil, status.Error(codes.InvalidArgument, "Name is required.")
	}
	au, err := s.authenticatedUser(ctx)
	if err != nil {
		return nil, err
	}
	store, err := s.storeRepo.LookupStore(ctx, req.Name)
	if err != nil {
		if notFound := new(stores.StoreNotFoundError); errors.As(err, &notFound) {
			if !au.Superuser {
				return nil, status.Error(codes.PermissionDenied, "Permission denied.")
			}
			return nil, status.Error(codes.NotFound, "Store not found.")
		}
		return nil, status.Error(codes.Internal, "Internal error.")
	}
	return store, nil
}

func (s *Server) ListStores(ctx context.Context, req *streckuv1.ListStoresRequest) (*streckuv1.ListStoresResponse, error) {
	if req.PageSize < 0 {
		return nil, status.Error(codes.InvalidArgument, "Page size must be non-negative.")
	}
	if req.PageSize > 0 || req.PageToken != "" {
		return nil, status.Error(codes.Unimplemented, "Pagination is not implemented.")
	}
	allStores, err := s.storeRepo.ListStores(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, "Internal error.")
	}
	return &streckuv1.ListStoresResponse{
		Stores:        allStores,
		NextPageToken: "",
	}, nil
}

func (s *Server) CreateStore(ctx context.Context, req *streckuv1.CreateStoreRequest) (*streckuv1.Store, error) {
	store := req.Store
	store.Name = stores.GenerateName()
	if err := stores.Validate(store); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	if err := s.storeRepo.CreateStore(ctx, store); err != nil {
		return nil, status.Error(codes.Internal, "Internal error.")
	}
	return store, nil
}
