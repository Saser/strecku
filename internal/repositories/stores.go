package repositories

import (
	"context"

	pb "github.com/Saser/strecku/api/v1"
)

type Stores interface {
	// Lookup returns the store corresponding to the given name, or returns a
	// non-nil error otherwise. The name will be validated using package
	// stores. If no store is found, a NotFound error will be returned.
	Lookup(ctx context.Context, name string) (*pb.Store, error)

	// List returns a list of all stores.
	List(ctx context.Context) ([]*pb.Store, error)

	// Create creates a new store resource based on the given store. The
	// given store will be validated using package stores. If a store
	// already If a store already exists with the given name, an Exists
	// error will be returned.
	Create(ctx context.Context, store *pb.Store) error

	// Update updates an existing store to the version specified by the given
	// store. The given store will be validated using package store. The name
	// of the given store is used to identify which store to update. If no
	// store with that name exists, a NotFound error will be returned.
	Update(ctx context.Context, store *pb.Store) error

	// Delete deletes the store corresponding to the given name. The name
	// will be validated using package stores. If no store with that name
	// exists, a NotFound error will be returned.
	Delete(ctx context.Context, name string) error
}
