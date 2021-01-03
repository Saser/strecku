package service

import (
	"context"
	"errors"

	pb "github.com/Saser/strecku/api/v1"
	"github.com/Saser/strecku/resourcename"
	"github.com/Saser/strecku/resources/stores"
	"github.com/Saser/strecku/resources/stores/payments"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *Service) GetPayment(ctx context.Context, req *pb.GetPaymentRequest) (*pb.Payment, error) {
	name := req.Name
	if err := payments.ValidateName(name); err != nil {
		switch {
		case errors.Is(err, resourcename.ErrInvalidName):
			return nil, status.Error(codes.InvalidArgument, err.Error())
		default:
			return nil, internalError
		}
	}
	payment, err := s.paymentRepo.LookupPayment(ctx, name)
	if err != nil {
		if notFound := new(payments.NotFoundError); errors.As(err, &notFound) {
			return nil, status.Error(codes.NotFound, notFound.Error())
		}
		return nil, internalError
	}
	return payment, nil
}

func (s *Service) ListPayments(ctx context.Context, req *pb.ListPaymentsRequest) (*pb.ListPaymentsResponse, error) {
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
	predicate := func(payment *pb.Payment) bool {
		parent, err := payments.Parent(payment.Name)
		if err != nil {
			return false
		}
		return parent == req.Parent
	}
	filtered, err := s.paymentRepo.FilterPayments(ctx, predicate)
	if err != nil {
		return nil, internalError
	}
	return &pb.ListPaymentsResponse{
		Payments:      filtered,
		NextPageToken: "",
	}, nil
}

func (s *Service) CreatePayment(ctx context.Context, req *pb.CreatePaymentRequest) (*pb.Payment, error) {
	if err := stores.ValidateName(req.Parent); err != nil {
		switch {
		case errors.Is(err, resourcename.ErrInvalidName):
			return nil, status.Errorf(codes.InvalidArgument, "invalid parent: %v", err)
		default:
			return nil, internalError
		}
	}
	payment := req.Payment
	payment.Name = payments.GenerateName(req.Parent)
	if err := payments.Validate(payment); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid payment: %v", err)
	}
	if err := s.paymentRepo.CreatePayment(ctx, payment); err != nil {
		return nil, internalError
	}
	return payment, nil
}

func (s *Service) UpdatePayment(ctx context.Context, req *pb.UpdatePaymentRequest) (*pb.Payment, error) {
	src := req.Payment
	dst, err := s.GetPayment(ctx, &pb.GetPaymentRequest{Name: src.Name})
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
			case "description":
				dst.Description = src.Description
			case "amount_cents":
				dst.AmountCents = src.AmountCents
			default:
				return nil, status.Errorf(codes.InvalidArgument, "update not implemented for path %q", path)
			}
		}
	}
	if err := payments.Validate(dst); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid payment: %v", err)
	}
	if err := s.paymentRepo.UpdatePayment(ctx, dst); err != nil {
		switch err {
		case payments.ErrUpdateUser:
			return nil, status.Errorf(codes.InvalidArgument, "invalid update: %v", err)
		default:
			return nil, internalError
		}
	}
	return dst, nil
}

func (s *Service) DeletePayment(ctx context.Context, req *pb.DeletePaymentRequest) (*emptypb.Empty, error) {
	if err := payments.ValidateName(req.Name); err != nil {
		switch {
		case errors.Is(err, resourcename.ErrInvalidName):
			return nil, status.Error(codes.InvalidArgument, err.Error())
		default:
			return nil, internalError
		}
	}
	if err := s.paymentRepo.DeletePayment(ctx, req.Name); err != nil {
		if notFound := new(payments.NotFoundError); errors.As(err, &notFound) {
			return nil, status.Error(codes.NotFound, notFound.Error())
		}
		return nil, internalError
	}
	return new(emptypb.Empty), nil
}
