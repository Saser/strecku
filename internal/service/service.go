package service

import (
	pb "github.com/Saser/strecku/api/v1"
	"github.com/Saser/strecku/internal/repositories"
	"github.com/Saser/strecku/resources/stores"
	"github.com/Saser/strecku/resources/stores/memberships"
	"github.com/Saser/strecku/resources/stores/payments"
	"github.com/Saser/strecku/resources/stores/products"
	"github.com/Saser/strecku/resources/stores/purchases"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var internalError = status.Error(codes.Internal, "internal error")

type Service struct {
	pb.UnimplementedStreckUServer

	userRepo       repositories.Users
	storeRepo      *stores.Repository
	membershipRepo *memberships.Repository
	productRepo    *products.Repository
	purchaseRepo   *purchases.Repository
	paymentRepo    *payments.Repository
}

func New(
	userRepo repositories.Users,
	storeRepo *stores.Repository,
	membershipRepo *memberships.Repository,
	productRepo *products.Repository,
	purchaseRepo *purchases.Repository,
	paymentRepo *payments.Repository,
) *Service {
	return &Service{
		userRepo:       userRepo,
		storeRepo:      storeRepo,
		membershipRepo: membershipRepo,
		productRepo:    productRepo,
		purchaseRepo:   purchaseRepo,
		paymentRepo:    paymentRepo,
	}
}
