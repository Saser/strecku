package inmemory

import (
	"context"
	"fmt"
	"sync"

	streckuv1 "github.com/Saser/strecku/backend/gen/api/v1"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type StoreAPI struct {
	mu     sync.Mutex
	stores map[string]*streckuv1.Store
}

func NewStoreAPI() *StoreAPI {
	return &StoreAPI{
		stores: make(map[string]*streckuv1.Store),
	}
}

func (s *StoreAPI) GetStore(_ context.Context, req *streckuv1.GetStoreRequest) (*streckuv1.GetStoreResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	store, ok := s.stores[req.Name]
	if !ok {
		return nil, status.Error(codes.NotFound, "Store resource not found")
	}
	return &streckuv1.GetStoreResponse{
		Store: store,
	}, nil
}

func (s *StoreAPI) ListStores(context.Context, *streckuv1.ListStoresRequest) (*streckuv1.ListStoresResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	stores := make([]*streckuv1.Store, 0, len(s.stores))
	for _, store := range s.stores {
		stores = append(stores, store)
	}
	return &streckuv1.ListStoresResponse{
		Stores: stores,
	}, nil
}

func (s *StoreAPI) CreateStore(_ context.Context, req *streckuv1.CreateStoreRequest) (*streckuv1.CreateStoreResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	newStore := req.Store
	newStore.Name = fmt.Sprintf("stores/%s", uuid.New().String())
	s.stores[newStore.Name] = newStore
	return &streckuv1.CreateStoreResponse{
		Store: newStore,
	}, nil
}
