package service

import (
	"context"
	"errors"

	pb "github.com/Saser/strecku/api/v1"
	"github.com/Saser/strecku/resourcename"
	"github.com/Saser/strecku/resources/stores"
	"github.com/Saser/strecku/resources/stores/products"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *Service) GetProduct(ctx context.Context, req *pb.GetProductRequest) (*pb.Product, error) {
	name := req.Name
	if err := products.ValidateName(name); err != nil {
		switch {
		case errors.Is(err, resourcename.ErrInvalidName):
			return nil, status.Error(codes.InvalidArgument, err.Error())
		default:
			return nil, internalError
		}
	}
	product, err := s.productRepo.LookupProduct(ctx, name)
	if err != nil {
		if notFound := new(products.NotFoundError); errors.As(err, &notFound) {
			return nil, status.Error(codes.NotFound, notFound.Error())
		}
		return nil, internalError
	}
	return product, nil
}

func (s *Service) ListProducts(ctx context.Context, req *pb.ListProductsRequest) (*pb.ListProductsResponse, error) {
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
	predicate := func(product *pb.Product) bool {
		parent, err := products.Parent(product.Name)
		if err != nil {
			return false
		}
		return parent == req.Parent
	}
	filtered, err := s.productRepo.FilterProducts(ctx, predicate)
	if err != nil {
		return nil, internalError
	}
	return &pb.ListProductsResponse{
		Products:      filtered,
		NextPageToken: "",
	}, nil
}

func (s *Service) CreateProduct(ctx context.Context, req *pb.CreateProductRequest) (*pb.Product, error) {
	if err := stores.ValidateName(req.Parent); err != nil {
		switch {
		case errors.Is(err, resourcename.ErrInvalidName):
			return nil, status.Errorf(codes.InvalidArgument, "invalid parent: %v", err)
		default:
			return nil, internalError
		}
	}
	product := req.Product
	product.Name = products.GenerateName(req.Parent)
	if err := products.Validate(product); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid product: %v", err)
	}
	if err := s.productRepo.CreateProduct(ctx, product); err != nil {
		return nil, internalError
	}
	return product, nil
}

func (s *Service) UpdateProduct(ctx context.Context, req *pb.UpdateProductRequest) (*pb.Product, error) {
	src := req.Product
	dst, err := s.GetProduct(ctx, &pb.GetProductRequest{Name: src.Name})
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
			case "full_price_cents":
				dst.FullPriceCents = src.FullPriceCents
			case "discount_price_cents":
				dst.DiscountPriceCents = src.DiscountPriceCents
			default:
				return nil, status.Errorf(codes.Internal, "update not implemented for path %q", path)
			}
		}
	}
	if err := products.Validate(dst); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid product: %v", err)
	}
	if err := s.productRepo.UpdateProduct(ctx, dst); err != nil {
		return nil, internalError
	}
	return dst, nil
}

func (s *Service) DeleteProduct(ctx context.Context, req *pb.DeleteProductRequest) (*emptypb.Empty, error) {
	if err := products.ValidateName(req.Name); err != nil {
		switch {
		case errors.Is(err, resourcename.ErrInvalidName):
			return nil, status.Error(codes.InvalidArgument, err.Error())
		default:
			return nil, internalError
		}
	}
	if err := s.productRepo.DeleteProduct(ctx, req.Name); err != nil {
		if notFound := new(products.NotFoundError); errors.As(err, &notFound) {
			return nil, status.Error(codes.NotFound, notFound.Error())
		}
		return nil, internalError
	}
	return new(emptypb.Empty), nil
}
