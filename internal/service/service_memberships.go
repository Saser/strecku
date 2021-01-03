package service

import (
	"context"
	"errors"

	pb "github.com/Saser/strecku/api/v1"
	"github.com/Saser/strecku/resourcename"
	"github.com/Saser/strecku/resources/stores"
	"github.com/Saser/strecku/resources/stores/memberships"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *Service) GetMembership(ctx context.Context, req *pb.GetMembershipRequest) (*pb.Membership, error) {
	name := req.Name
	if err := memberships.ValidateName(name); err != nil {
		switch {
		case errors.Is(err, resourcename.ErrInvalidName):
			return nil, status.Error(codes.InvalidArgument, err.Error())
		default:
			return nil, internalError
		}
	}
	membership, err := s.membershipRepo.LookupMembership(ctx, name)
	if err != nil {
		if notFound := new(memberships.NotFoundError); errors.As(err, &notFound) {
			return nil, status.Error(codes.NotFound, notFound.Error())
		}
		return nil, internalError
	}
	return membership, nil
}

func (s *Service) ListMemberships(ctx context.Context, req *pb.ListMembershipsRequest) (*pb.ListMembershipsResponse, error) {
	if err := stores.ValidateName(req.Parent); err != nil {
		switch {
		case errors.Is(err, resourcename.ErrInvalidName):
			return nil, status.Errorf(codes.InvalidArgument, "invalid parent: %v", err)
		default:
			return nil, internalError
		}
	}
	if req.PageSize < 0 {
		return nil, status.Errorf(codes.InvalidArgument, "negative page size: %d", req.PageSize)
	}
	if req.PageSize > 0 || req.PageToken != "" {
		return nil, status.Error(codes.Unimplemented, "pagination is not implemented")
	}
	predicate := func(membership *pb.Membership) bool {
		parent, err := memberships.Parent(membership.Name)
		if err != nil {
			return false
		}
		return parent == req.Parent
	}
	filtered, err := s.membershipRepo.FilterMemberships(ctx, predicate)
	if err != nil {
		return nil, internalError
	}
	return &pb.ListMembershipsResponse{
		Memberships:   filtered,
		NextPageToken: "",
	}, nil
}

func (s *Service) CreateMembership(ctx context.Context, req *pb.CreateMembershipRequest) (*pb.Membership, error) {
	if err := stores.ValidateName(req.Parent); err != nil {
		switch {
		case errors.Is(err, resourcename.ErrInvalidName):
			return nil, status.Errorf(codes.InvalidArgument, "invalid parent: %v", err)
		default:
			return nil, internalError
		}
	}
	membership := req.Membership
	membership.Name = memberships.GenerateName(req.Parent)
	if err := memberships.Validate(membership); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid membership: %v", err)
	}
	if err := s.membershipRepo.CreateMembership(ctx, membership); err != nil {
		if exists := new(memberships.ExistsError); errors.As(err, &exists) {
			return nil, status.Error(codes.AlreadyExists, exists.Error())
		}
		return nil, internalError
	}
	return membership, nil
}

func (s *Service) UpdateMembership(ctx context.Context, req *pb.UpdateMembershipRequest) (*pb.Membership, error) {
	src := req.Membership
	dst, err := s.GetMembership(ctx, &pb.GetMembershipRequest{Name: src.Name})
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
			case "user":
				return nil, status.Errorf(codes.InvalidArgument, `field "user" cannot be updated`)
			case "administrator":
				dst.Administrator = src.Administrator
			case "discount":
				dst.Discount = src.Discount
			default:
				return nil, status.Errorf(codes.Internal, "update not implemented for path %q", path)
			}
		}
	}
	if err := memberships.Validate(dst); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid membership: %v", err)
	}
	if err := s.membershipRepo.UpdateMembership(ctx, dst); err != nil {
		switch err {
		case memberships.ErrUpdateUser:
			return nil, status.Errorf(codes.InvalidArgument, "invalid update: %v", err)
		default:
			return nil, internalError
		}
	}
	return dst, nil
}

func (s *Service) DeleteMembership(ctx context.Context, req *pb.DeleteMembershipRequest) (*emptypb.Empty, error) {
	if err := memberships.ValidateName(req.Name); err != nil {
		switch {
		case errors.Is(err, resourcename.ErrInvalidName):
			return nil, status.Error(codes.InvalidArgument, err.Error())
		default:
			return nil, internalError
		}
	}
	if err := s.membershipRepo.DeleteMembership(ctx, req.Name); err != nil {
		if notFound := new(memberships.NotFoundError); errors.As(err, &notFound) {
			return nil, status.Error(codes.NotFound, notFound.Error())
		}
		return nil, internalError
	}
	return new(emptypb.Empty), nil
}
