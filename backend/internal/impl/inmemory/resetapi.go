package inmemory

import (
	"context"
	"os"

	testingv1 "github.com/Saser/strecku/backend/gen/api/testing/v1"
	streckuv1 "github.com/Saser/strecku/backend/gen/api/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (i *Impl) Reset(context.Context, *testingv1.ResetRequest) (*testingv1.ResetResponse, error) {
	if os.Getenv("STRECKU_INTEGRATION_TESTING") != "1" {
		return nil, status.Error(codes.Unimplemented, "not implemented")
	}
	i.mu.Lock()
	defer i.mu.Unlock()
	i.users = make(map[string]*streckuv1.User)
	i.stores = make(map[string]*streckuv1.Store)
	i.roles = make(map[string]*streckuv1.Role)
	i.products = make(map[string]*streckuv1.Product)
	return &testingv1.ResetResponse{}, nil
}
