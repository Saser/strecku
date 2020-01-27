package inmemory

import (
	"context"
	"fmt"

	streckuv1 "github.com/Saser/strecku/backend/gen/api/v1"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (i *Impl) ListRoles(context.Context, *streckuv1.ListRolesRequest) (*streckuv1.ListRolesResponse, error) {
	i.mu.Lock()
	defer i.mu.Unlock()
	roles := make([]*streckuv1.Role, 0, len(i.roles))
	for _, role := range i.roles {
		roles = append(roles, role)
	}
	return &streckuv1.ListRolesResponse{
		Roles: roles,
	}, nil
}

func (i *Impl) GetRole(_ context.Context, req *streckuv1.GetRoleRequest) (*streckuv1.GetRoleResponse, error) {
	i.mu.Lock()
	defer i.mu.Unlock()
	role, ok := i.roles[req.Name]
	if !ok {
		return nil, status.Error(codes.NotFound, "Role resource not found")
	}
	return &streckuv1.GetRoleResponse{
		Role: role,
	}, nil
}

func (i *Impl) CreateRole(_ context.Context, req *streckuv1.CreateRoleRequest) (*streckuv1.CreateRoleResponse, error) {
	newRole := req.Role
	userName := newRole.UserName
	if _, ok := i.users[userName]; !ok {
		return nil, status.Error(codes.InvalidArgument, "invalid User resource name")
	}
	storeName := newRole.StoreName
	if _, ok := i.stores[storeName]; !ok {
		return nil, status.Error(codes.InvalidArgument, "invalid Store resource name")
	}
	newRole.Name = fmt.Sprintf("roles/%s", uuid.New().String())
	i.roles[newRole.Name] = newRole
	return &streckuv1.CreateRoleResponse{
		Role: newRole,
	}, nil
}
