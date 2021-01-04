package inmemory

import (
	"context"
	"log"

	streckuv1 "github.com/Saser/strecku/backend/gen/api/strecku/v1"
	testingv1 "github.com/Saser/strecku/backend/gen/api/testing/v1"
)

func (i *Impl) Reset(_ context.Context, req *testingv1.ResetRequest) (*testingv1.ResetResponse, error) {
	i.mu.Lock()
	defer i.mu.Unlock()
	i.users = make(map[string]*streckuv1.User)
	i.stores = make(map[string]*streckuv1.Store)
	i.roles = make(map[string]*streckuv1.Role)
	i.products = make(map[string]*streckuv1.Product)
	log.Printf("reason=%v", req.Reason)
	return &testingv1.ResetResponse{}, nil
}
