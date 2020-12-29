package repositories

import (
	"context"

	pb "github.com/Saser/strecku/api/v1"
	"github.com/Saser/strecku/resources/stores"
)

type InMemoryStores struct {
	stores map[string]*pb.Store // name -> store
}

var _ Stores = (*InMemoryStores)(nil)

func NewInMemoryStores() *InMemoryStores {
	return &InMemoryStores{
		stores: make(map[string]*pb.Store),
	}
}

func (u *InMemoryStores) Lookup(ctx context.Context, name string) (*pb.Store, error) {
	if err := stores.ValidateName(name); err != nil {
		return nil, err
	}
	store, ok := u.stores[name]
	if !ok {
		return nil, &NotFound{Name: name}
	}
	return stores.Clone(store), nil
}

func (u *InMemoryStores) List(ctx context.Context) ([]*pb.Store, error) {
	allStores := make([]*pb.Store, 0, len(u.stores))
	for _, store := range u.stores {
		allStores = append(allStores, stores.Clone(store))
	}
	return allStores, nil
}

func (u *InMemoryStores) Create(ctx context.Context, store *pb.Store) error {
	if err := stores.Validate(store); err != nil {
		return err
	}
	if _, exists := u.stores[store.Name]; exists {
		return &Exists{Name: store.Name}
	}
	u.stores[store.Name] = stores.Clone(store)
	return nil
}

func (u *InMemoryStores) Update(ctx context.Context, store *pb.Store) error {
	if err := stores.Validate(store); err != nil {
		return err
	}
	if _, exists := u.stores[store.Name]; !exists {
		return &NotFound{Name: store.Name}
	}
	u.stores[store.Name] = stores.Clone(store)
	return nil
}

func (u *InMemoryStores) Delete(ctx context.Context, name string) error {
	store, err := u.Lookup(ctx, name)
	if err != nil {
		return err
	}
	delete(u.stores, store.Name)
	return nil
}
