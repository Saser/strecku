package testing

import (
	"context"
	"strings"

	streckuv1 "github.com/Saser/strecku/backend/gen/api/strecku/v1"
	"github.com/golang/protobuf/proto"
)

func (i *IntegrationTestSuite) TestStoreAPI_ListStores() {
	ctx := context.Background()
	c := streckuv1.NewStoreAPIClient(i.cc)
	res, err := c.ListStores(ctx, &streckuv1.ListStoresRequest{})
	i.Require().NoError(err)
	i.Assert().NotEmpty(res.Stores)
}

func (i *IntegrationTestSuite) TestStoreAPI_GetStore() {
	ctx := context.Background()
	c := streckuv1.NewStoreAPIClient(i.cc)
	listRes, err := c.ListStores(ctx, &streckuv1.ListStoresRequest{})
	i.Require().NoError(err)
	i.Assert().GreaterOrEqual(len(listRes.Stores), 1)
	for _, store := range listRes.Stores {
		getRes, err := c.GetStore(ctx, &streckuv1.GetStoreRequest{
			Name: store.Name,
		})
		i.Require().NoError(err)
		i.Assert().Truef(proto.Equal(store, getRes.Store), "store=%v,getRes.Store=%v", store, getRes.Store)
	}
}

func (i *IntegrationTestSuite) TestStoreAPI_CreateStore() {
	ctx := context.Background()
	c := streckuv1.NewStoreAPIClient(i.cc)
	newStore := &streckuv1.Store{
		DisplayName: "Another Store",
	}
	i.Run("create", func() {
		res, err := c.CreateStore(ctx, &streckuv1.CreateStoreRequest{
			Store: newStore,
		})
		i.Require().NoError(err)
		i.Assert().Truef(strings.HasPrefix(res.Store.Name, "stores/"), "res.Store.Name=%v", res.Store.Name)
		newStore.Name = res.Store.Name
	})
	i.Run("get", func() {
		res, err := c.GetStore(ctx, &streckuv1.GetStoreRequest{
			Name: newStore.Name,
		})
		i.Require().NoError(err)
		i.Assert().Truef(proto.Equal(newStore, res.Store), "newStore=%v,res.Store=%v", newStore, res.Store)
	})
	i.Run("list", func() {
		res, err := c.ListStores(ctx, &streckuv1.ListStoresRequest{})
		i.Require().NoError(err)
		i.Assert().Equal(2, len(res.Stores))
		ok := false
		for _, store := range res.Stores {
			ok = ok || proto.Equal(newStore, store)
			if ok {
				break
			}
		}
		i.Assert().Truef(ok, "newStore=%v,res.Stores=%v", newStore, res.Stores)
	})
}
