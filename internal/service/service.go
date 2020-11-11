package service

import (
	"context"
	"errors"

	pb "github.com/Saser/strecku/api/v1"
	"github.com/Saser/strecku/resources/stores"
	"github.com/Saser/strecku/resources/stores/memberships"
	"github.com/Saser/strecku/resources/stores/products"
	"github.com/Saser/strecku/resources/stores/purchases"
	"github.com/Saser/strecku/resources/users"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

var internalError = status.Error(codes.Internal, "internal error")

type Service struct {
	pb.UnimplementedStreckUServer

	userRepo       *users.Repository
	storeRepo      *stores.Repository
	membershipRepo *memberships.Repository
	productRepo    *products.Repository
	purchaseRepo   *purchases.Repository
}

func New(
	userRepo *users.Repository,
	storeRepo *stores.Repository,
	membershipRepo *memberships.Repository,
	productRepo *products.Repository,
	purchaseRepo *purchases.Repository,
) *Service {
	return &Service{
		userRepo:       userRepo,
		storeRepo:      storeRepo,
		membershipRepo: membershipRepo,
		productRepo:    productRepo,
		purchaseRepo:   purchaseRepo,
	}
}

func (s *Service) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.User, error) {
	name := req.Name
	if err := users.ValidateName(name); err != nil {
		switch err {
		case users.ErrNameInvalidFormat:
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
	if err := users.Validate(dst); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid user: %v", err)
	}
	if err := s.userRepo.UpdateUser(ctx, dst); err != nil {
		if exists := new(users.ExistsError); errors.As(err, &exists) {
			return nil, status.Error(codes.AlreadyExists, exists.Error())
		}
		return nil, internalError
	}
	return dst, nil
}

func (s *Service) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*emptypb.Empty, error) {
	if err := users.ValidateName(req.Name); err != nil {
		switch err {
		case users.ErrNameInvalidFormat:
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
	store, err := s.storeRepo.LookupStore(ctx, name)
	if err != nil {
		if notFound := new(stores.NotFoundError); errors.As(err, &notFound) {
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
	allStores, err := s.storeRepo.ListStores(ctx)
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
	if err := s.storeRepo.CreateStore(ctx, store); err != nil {
		if exists := new(stores.ExistsError); errors.As(err, &exists) {
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
	if err := s.storeRepo.UpdateStore(ctx, dst); err != nil {
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
	if err := s.storeRepo.DeleteStore(ctx, req.Name); err != nil {
		if notFound := new(stores.NotFoundError); errors.As(err, &notFound) {
			return nil, status.Error(codes.NotFound, notFound.Error())
		}
		return nil, internalError
	}
	return new(emptypb.Empty), nil
}

func (s *Service) GetMembership(ctx context.Context, req *pb.GetMembershipRequest) (*pb.Membership, error) {
	name := req.Name
	if err := memberships.ValidateName(name); err != nil {
		switch err {
		case memberships.ErrNameInvalidFormat:
			return nil, status.Errorf(codes.InvalidArgument, "invalid name: %v", err)
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
		switch err {
		case stores.ErrNameInvalidFormat:
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
		switch err {
		case memberships.ErrNameInvalidFormat:
			return nil, status.Errorf(codes.InvalidArgument, "invalid name: %v", err)
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

func (s *Service) GetProduct(ctx context.Context, req *pb.GetProductRequest) (*pb.Product, error) {
	name := req.Name
	if err := products.ValidateName(name); err != nil {
		switch err {
		case products.ErrNameInvalidFormat:
			return nil, status.Errorf(codes.InvalidArgument, "invalid name: %v", err)
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
		switch err {
		case stores.ErrNameInvalidFormat:
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
		switch err {
		case products.ErrNameInvalidFormat:
			return nil, status.Errorf(codes.InvalidArgument, "invalid name: %v", err)
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
