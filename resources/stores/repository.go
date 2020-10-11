package stores

import (
	"context"
	"fmt"
	"testing"

	pb "github.com/Saser/strecku/api/v1"
	"google.golang.org/protobuf/proto"
)

type Repository struct {
	stores map[string]*pb.Store // name -> store
}

type NotFoundError struct {
	Name string
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("store not found: %q", e.Name)
}

func (e *NotFoundError) Is(target error) bool {
	other, ok := target.(*NotFoundError)
	if !ok {
		return false
	}
	return e.Name == other.Name
}

type ExistsError struct {
	Name string
}

func (e *ExistsError) Error() string {
	return fmt.Sprintf("store exists: %q", e.Name)
}

func (e *ExistsError) Is(target error) bool {
	other, ok := target.(*ExistsError)
	if !ok {
		return false
	}
	return e.Name == other.Name
}

func Clone(store *pb.Store) *pb.Store {
	return proto.Clone(store).(*pb.Store)
}

func NewRepository() *Repository {
	return newRepository(make(map[string]*pb.Store))
}

func SeedRepository(t *testing.T, stores []*pb.Store) *Repository {
	t.Helper()
	mStores := make(map[string]*pb.Store, len(stores))
	for _, store := range stores {
		if got := Validate(store); got != nil {
			t.Errorf("Validate(%v) = %v; want %v", store, got, nil)
		}
		mStores[store.Name] = store
	}
	if t.Failed() {
		t.FailNow()
	}
	return newRepository(mStores)
}

func newRepository(stores map[string]*pb.Store) *Repository {
	return &Repository{
		stores: stores,
	}
}

func (r *Repository) LookupStore(_ context.Context, name string) (*pb.Store, error) {
	if err := ValidateName(name); err != nil {
		return nil, err
	}
	store, ok := r.stores[name]
	if !ok {
		return nil, &NotFoundError{Name: name}
	}
	return store, nil
}

func (r *Repository) ListStores(ctx context.Context) ([]*pb.Store, error) {
	return r.FilterStores(ctx, func(*pb.Store) bool { return true })
}

func (r *Repository) FilterStores(_ context.Context, predicate func(*pb.Store) bool) ([]*pb.Store, error) {
	var filtered []*pb.Store
	for _, store := range r.stores {
		if predicate(store) {
			filtered = append(filtered, store)
		}
	}
	return filtered, nil
}

func (r *Repository) CreateStore(_ context.Context, store *pb.Store) error {
	if _, exists := r.stores[store.Name]; exists {
		return &ExistsError{Name: store.Name}
	}
	r.stores[store.Name] = store
	return nil
}

func (r *Repository) UpdateStore(_ context.Context, updated *pb.Store) error {
	if err := Validate(updated); err != nil {
		return err
	}
	if _, exists := r.stores[updated.Name]; !exists {
		return &NotFoundError{Name: updated.Name}
	}
	r.stores[updated.Name] = updated
	return nil
}

func (r *Repository) DeleteStore(_ context.Context, name string) error {
	if err := ValidateName(name); err != nil {
		return err
	}
	if _, exists := r.stores[name]; !exists {
		return &NotFoundError{Name: name}
	}
	delete(r.stores, name)
	return nil
}
