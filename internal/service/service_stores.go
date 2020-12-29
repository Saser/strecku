package service

import (
	"context"
	"errors"

	pb "github.com/Saser/strecku/api/v1"
	"github.com/Saser/strecku/internal/repositories"
	"github.com/Saser/strecku/resources/stores"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *Service) GetStore(ctx context.Context, req *pb.GetStoreRequest) (*pb.Store, error) {
	name := req.Name
	if err := stores.ValidateName(name); err != nil {
		switch err {
		case stores.ErrNameInvalidFormat:
			return nil, status.Errorf(codes.InvalidArgument, "invalid name: %v", err)
		default:
			return nil, internalError
		}
	}
	store, err := s.storeRepo.Lookup(ctx, name)
	if err != nil {
		if notFound := new(repositories.NotFound); errors.As(err, &notFound) {
			return nil, status.Error(codes.NotFound, notFound.Error())
		}
		return nil, internalError
	}
	return store, nil
}

func (s *Service) ListStores(ctx context.Context, req *pb.ListStoresRequest) (*pb.ListStoresResponse, error) {
	if req.PageSize < 0 {
		return nil, status.Errorf(codes.InvalidArgument, "negative page size: %d", req.PageSize)
	}
	if req.PageSize > 0 || req.PageToken != "" {
		return nil, status.Error(codes.Unimplemented, "pagination is not implemented")
	}
	allStores, err := s.storeRepo.List(ctx)
	if err != nil {
		return nil, internalError
	}
	return &pb.ListStoresResponse{
		Stores:        allStores,
		NextPageToken: "",
	}, nil
}

func (s *Service) CreateStore(ctx context.Context, req *pb.CreateStoreRequest) (*pb.Store, error) {
	store := req.Store
	store.Name = stores.GenerateName()
	if err := stores.Validate(store); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid store: %v", err)
	}
	if err := s.storeRepo.Create(ctx, store); err != nil {
		if exists := new(repositories.Exists); errors.As(err, &exists) {
			return nil, status.Error(codes.AlreadyExists, exists.Error())
		}
		return nil, internalError
	}
	return store, nil
}

func (s *Service) UpdateStore(ctx context.Context, req *pb.UpdateStoreRequest) (*pb.Store, error) {
	src := req.Store
	dst, err := s.GetStore(ctx, &pb.GetStoreRequest{Name: src.Name})
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
			case "display_name":
				dst.DisplayName = src.DisplayName
			default:
				return nil, status.Errorf(codes.Internal, "update not implemented for path %q", path)
			}
		}
	}
	if err := stores.Validate(dst); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid store: %v", err)
	}
	if err := s.storeRepo.Update(ctx, dst); err != nil {
		return nil, internalError
	}
	return dst, nil
}

func (s *Service) DeleteStore(ctx context.Context, req *pb.DeleteStoreRequest) (*emptypb.Empty, error) {
	if err := stores.ValidateName(req.Name); err != nil {
		switch err {
		case stores.ErrNameInvalidFormat:
			return nil, status.Errorf(codes.InvalidArgument, "invalid name: %v", err)
		default:
			return nil, internalError
		}
	}
	if err := s.storeRepo.Delete(ctx, req.Name); err != nil {
		if notFound := new(repositories.NotFound); errors.As(err, &notFound) {
			return nil, status.Error(codes.NotFound, notFound.Error())
		}
		return nil, internalError
	}
	return new(emptypb.Empty), nil
}
