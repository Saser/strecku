package testing

import (
	"context"
	"strings"

	streckuv1 "github.com/Saser/strecku/backend/gen/api/strecku/v1"
	"github.com/golang/protobuf/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (i *IntegrationTestSuite) TestUserAPI_ListUsers() {
	ctx := context.Background()
	c := streckuv1.NewUserAPIClient(i.cc)
	res, err := c.ListUsers(ctx, &streckuv1.ListUsersRequest{})
	i.Require().NoError(err)
	i.Assert().NotEmpty(res.Users)
}

func (i *IntegrationTestSuite) TestUserAPI_GetUser() {
	ctx := context.Background()
	c := streckuv1.NewUserAPIClient(i.cc)
	listRes, err := c.ListUsers(ctx, &streckuv1.ListUsersRequest{})
	i.Require().NoError(err)
	i.Assert().GreaterOrEqual(len(listRes.Users), 1)
	for _, user := range listRes.Users {
		getRes, err := c.GetUser(ctx, &streckuv1.GetUserRequest{
			Name: user.Name,
		})
		i.Require().NoError(err)
		i.Assert().Truef(proto.Equal(user, getRes.User), "user=%v,getRes.User=%v", user, getRes.User)
	}
}

func (i *IntegrationTestSuite) TestUserAPI_CreateUser() {
	ctx := context.Background()
	c := streckuv1.NewUserAPIClient(i.cc)
	newUser := &streckuv1.User{
		DisplayName:  "Saser Again",
		EmailAddress: "saseragain@saser.com",
	}
	i.Run("create", func() {
		res, err := c.CreateUser(ctx, &streckuv1.CreateUserRequest{
			User: newUser,
		})
		i.Require().NoError(err)
		i.Truef(strings.HasPrefix(res.User.Name, "users/"), "res.User.Name=%v", res.User.Name)
		newUser.Name = res.User.Name
	})
	i.Run("get", func() {
		res, err := c.GetUser(ctx, &streckuv1.GetUserRequest{
			Name: newUser.Name,
		})
		i.Require().NoError(err)
		i.Assert().Truef(proto.Equal(newUser, res.User), "newUser=%v,res.User=%v", newUser, res.User)
	})
	i.Run("list", func() {
		res, err := c.ListUsers(ctx, &streckuv1.ListUsersRequest{})
		i.Require().NoError(err)
		i.Assert().Equal(2, len(res.Users))
		ok := false
		for _, user := range res.Users {
			ok = ok || proto.Equal(newUser, user)
			if ok {
				break
			}
		}
		i.Assert().Truef(ok, "newUser=%v,res.Users=%v", newUser, res.Users)
	})
}

func (i *IntegrationTestSuite) TestUserAPI_CreateUser_Duplicate() {
	ctx := context.Background()
	c := streckuv1.NewUserAPIClient(i.cc)
	_, err := c.CreateUser(ctx, &streckuv1.CreateUserRequest{
		User: &streckuv1.User{
			DisplayName:  "Saser (Duplicate email address)",
			EmailAddress: "saser@saser.com",
		},
	})
	i.Require().Error(err)
	i.Assert().Equal(codes.AlreadyExists, status.Code(err))
}
