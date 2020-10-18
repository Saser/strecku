package service

import (
	"context"
	"errors"

	pb "github.com/Saser/strecku/api/v1"
	"github.com/Saser/strecku/resources/users"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

var internalError = status.Error(codes.Internal, "internal error")

type Service struct {
	pb.UnimplementedStreckUServer

	userRepo *users.Repository
}

func New(userRepo *users.Repository) *Service {
	return &Service{
		userRepo: userRepo,
	}
}

func (s *Service) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.User, error) {
	name := req.Name
	if err := users.ValidateName(name); err != nil {
		switch err {
		case users.ErrNameEmpty, users.ErrNameInvalidFormat:
			return nil, status.Errorf(codes.InvalidArgument, "invalid name: %v", err)
		default:
			return nil, internalError
		}
	}
	user, err := s.userRepo.LookupUser(ctx, name)
	if err != nil {
		if notFound := new(users.NotFoundError); errors.As(err, &notFound) {
			return nil, status.Error(codes.NotFound, notFound.Error())
		}
		return nil, internalError
	}
	return user, nil
}

func (s *Service) ListUsers(ctx context.Context, req *pb.ListUsersRequest) (*pb.ListUsersResponse, error) {
	if req.PageSize < 0 {
		return nil, status.Errorf(codes.InvalidArgument, "negative page size: %d", req.PageSize)
	}
	if req.PageSize > 0 || req.PageToken != "" {
		return nil, status.Error(codes.Unimplemented, "pagination is not implemented")
	}
	allUsers, err := s.userRepo.ListUsers(ctx)
	if err != nil {
		return nil, internalError
	}
	return &pb.ListUsersResponse{
		Users:         allUsers,
		NextPageToken: "",
	}, nil
}

func (s *Service) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.User, error) {
	user := req.User
	user.Name = users.GenerateName()
	if err := users.Validate(user); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid user: %v", err)
	}
	if err := users.ValidatePassword(req.Password); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid password: %v", err)
	}
	if err := s.userRepo.CreateUser(ctx, user, req.Password); err != nil {
		if exists := new(users.ExistsError); errors.As(err, &exists) {
			return nil, status.Error(codes.AlreadyExists, exists.Error())
		}
		return nil, internalError
	}
	return user, nil
}

func (s *Service) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.User, error) {
	src := req.User
	dst, err := s.GetUser(ctx, &pb.GetUserRequest{Name: src.Name})
	if err != nil {
		return nil, err
	}
	mask := req.UpdateMask
	if mask == nil {
		dst = src
	} else {
		if !mask.IsValid(dst) {
			return nil, status.Error(codes.InvalidArgument, "invalid update mask")
		}
		for _, path := range mask.Paths {
			switch path {
			case "email_address":
				dst.EmailAddress = src.EmailAddress
			case "display_name":
				dst.DisplayName = src.DisplayName
			default:
				return nil, status.Errorf(codes.Internal, "update not implemented for path %q", path)
			}
		}
	}
	if err := s.userRepo.UpdateUser(ctx, dst); err != nil {
		if exists := new(users.ExistsError); errors.As(err, &exists) {
			return nil, status.Error(codes.AlreadyExists, exists.Error())
		}
		return nil, status.Errorf(codes.InvalidArgument, "invalid user: %v", err)
	}
	return dst, nil
}

func (s *Service) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*emptypb.Empty, error) {
	if err := users.ValidateName(req.Name); err != nil {
		switch err {
		case users.ErrNameEmpty, users.ErrNameInvalidFormat:
			return nil, status.Errorf(codes.InvalidArgument, "invalid name: %v", err)
		default:
			return nil, internalError
		}
	}
	if err := s.userRepo.DeleteUser(ctx, req.Name); err != nil {
		if notFound := new(users.NotFoundError); errors.As(err, &notFound) {
			return nil, status.Error(codes.NotFound, notFound.Error())
		}
		return nil, internalError
	}
	return new(emptypb.Empty), nil
}
