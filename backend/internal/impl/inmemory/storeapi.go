package inmemory

import (
	"context"
	"fmt"

	streckuv1 "github.com/Saser/strecku/backend/gen/api/v1"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (i *Impl) GetStore(_ context.Context, req *streckuv1.GetStoreRequest) (*streckuv1.GetStoreResponse, error) {
	i.mu.Lock()
	defer i.mu.Unlock()
	store, ok := i.stores[req.Name]
	if !ok {
		return nil, status.Error(codes.NotFound, "Store resource not found")
	}
	return &streckuv1.GetStoreResponse{
		Store: store,
	}, nil
}

func (i *Impl) ListStores(context.Context, *streckuv1.ListStoresRequest) (*streckuv1.ListStoresResponse, error) {
	i.mu.Lock()
	defer i.mu.Unlock()
	stores := make([]*streckuv1.Store, 0, len(i.stores))
	for _, store := range i.stores {
		stores = append(stores, store)
	}
	return &streckuv1.ListStoresResponse{
		Stores: stores,
	}, nil
}

func (i *Impl) CreateStore(_ context.Context, req *streckuv1.CreateStoreRequest) (*streckuv1.CreateStoreResponse, error) {
	i.mu.Lock()
	defer i.mu.Unlock()
	newStore := req.Store
	newStore.Name = fmt.Sprintf("stores/%s", uuid.New().String())
	i.stores[newStore.Name] = newStore
	return &streckuv1.CreateStoreResponse{
		Store: newStore,
	}, nil
}
