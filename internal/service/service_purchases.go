package service

import (
	"context"
	"errors"

	pb "github.com/Saser/strecku/api/v1"
	"github.com/Saser/strecku/resources/stores"
	"github.com/Saser/strecku/resources/stores/purchases"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *Service) GetPurchase(ctx context.Context, req *pb.GetPurchaseRequest) (*pb.Purchase, error) {
	name := req.Name
	if err := purchases.ValidateName(name); err != nil {
		switch err {
		case purchases.ErrNameInvalidFormat:
			return nil, status.Errorf(codes.InvalidArgument, "invalid name: %v", err)
		default:
			return nil, internalError
		}
	}
	purchase, err := s.purchaseRepo.LookupPurchase(ctx, name)
	if err != nil {
		if notFound := new(purchases.NotFoundError); errors.As(err, &notFound) {
			return nil, status.Error(codes.NotFound, notFound.Error())
		}
		return nil, internalError
	}
	return purchase, nil
}

func (s *Service) ListPurchases(ctx context.Context, req *pb.ListPurchasesRequest) (*pb.ListPurchasesResponse, error) {
	if err := stores.ValidateName(req.Parent); err != nil {
		switch err {
		case stores.ErrNameInvalidFormat:
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
	predicate := func(purchase *pb.Purchase) bool {
		parent, err := purchases.Parent(purchase.Name)
		if err != nil {
			return false
		}
		return parent == req.Parent
	}
	filtered, err := s.purchaseRepo.FilterPurchases(ctx, predicate)
	if err != nil {
		return nil, internalError
	}
	return &pb.ListPurchasesResponse{
		Purchases:     filtered,
		NextPageToken: "",
	}, nil
}

func (s *Service) CreatePurchase(ctx context.Context, req *pb.CreatePurchaseRequest) (*pb.Purchase, error) {
	if err := stores.ValidateName(req.Parent); err != nil {
		switch err {
		case stores.ErrNameInvalidFormat:
			return nil, status.Errorf(codes.InvalidArgument, "invalid parent: %v", err)
		default:
			return nil, internalError
		}
	}
	purchase := req.Purchase
	purchase.Name = purchases.GenerateName(req.Parent)
	if err := purchases.Validate(purchase); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid purchase: %v", err)
	}
	if err := s.purchaseRepo.CreatePurchase(ctx, purchase); err != nil {
		return nil, internalError
	}
	return purchase, nil
}

func (s *Service) UpdatePurchase(ctx context.Context, req *pb.UpdatePurchaseRequest) (*pb.Purchase, error) {
	src := req.Purchase
	dst, err := s.GetPurchase(ctx, &pb.GetPurchaseRequest{Name: src.Name})
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
			case "lines":
				dst.Lines = src.Lines
			default:
				return nil, status.Errorf(codes.InvalidArgument, "update not implemented for path %q", path)
			}
		}
	}
	if err := purchases.Validate(dst); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid purchase: %v", err)
	}
	if err := s.purchaseRepo.UpdatePurchase(ctx, dst); err != nil {
		switch err {
		case purchases.ErrUpdateUser:
			return nil, status.Errorf(codes.InvalidArgument, "invalid update: %v", err)
		default:
			return nil, internalError
		}
	}
	return dst, nil
}

func (s *Service) DeletePurchase(ctx context.Context, req *pb.DeletePurchaseRequest) (*emptypb.Empty, error) {
	if err := purchases.ValidateName(req.Name); err != nil {
		switch err {
		case purchases.ErrNameInvalidFormat:
			return nil, status.Errorf(codes.InvalidArgument, "invalid name: %v", err)
		default:
			return nil, internalError
		}
	}
	if err := s.purchaseRepo.DeletePurchase(ctx, req.Name); err != nil {
		if notFound := new(purchases.NotFoundError); errors.As(err, &notFound) {
			return nil, status.Error(codes.NotFound, notFound.Error())
		}
		return nil, internalError
	}
	return new(emptypb.Empty), nil
}
