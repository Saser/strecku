package inmemory

import (
	"context"
	"fmt"
	"strings"

	streckuv1 "github.com/Saser/strecku/backend/gen/api/v1"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (i *Impl) ListProducts(_ context.Context, req *streckuv1.ListProductsRequest) (*streckuv1.ListProductsResponse, error) {
	i.mu.Lock()
	defer i.mu.Unlock()
	products := make([]*streckuv1.Product, 0)
	for _, product := range i.products {
		if strings.HasPrefix(product.Name, req.Parent) {
			products = append(products, product)
		}
	}
	return &streckuv1.ListProductsResponse{
		Products: products,
	}, nil
}

func (i *Impl) GetProduct(_ context.Context, req *streckuv1.GetProductRequest) (*streckuv1.GetProductResponse, error) {
	i.mu.Lock()
	defer i.mu.Unlock()
	product, ok := i.products[req.Name]
	if !ok {
		return nil, status.Error(codes.NotFound, "Product resource not found")
	}
	return &streckuv1.GetProductResponse{
		Product: product,
	}, nil
}

func (i *Impl) CreateProduct(_ context.Context, req *streckuv1.CreateProductRequest) (*streckuv1.CreateProductResponse, error) {
	i.mu.Lock()
	defer i.mu.Unlock()
	if _, ok := i.stores[req.Parent]; !ok {
		return nil, status.Error(codes.InvalidArgument, "invalid parent resource name")
	}
	newProduct := req.Product
	newProduct.Name = fmt.Sprintf("%s/products/%s", req.Parent, uuid.New().String())
	i.products[newProduct.Name] = newProduct
	return &streckuv1.CreateProductResponse{
		Product: newProduct,
	}, nil
}
