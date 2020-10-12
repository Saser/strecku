package service

import (
	"context"
	"errors"

	pb "github.com/Saser/strecku/api/v1"
	"github.com/Saser/strecku/resources/users"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
