package stores

import (
	"context"
	"fmt"

	streckuv1 "github.com/Saser/strecku/saser/strecku/v1"
)

type Repository struct {
	stores map[string]*streckuv1.Store // name -> store
}

type StoreNotFoundError struct {
	Name string
}

func (e *StoreNotFoundError) Error() string {
	return fmt.Sprintf("store not found: %q", e.Name)
}

func (e *StoreNotFoundError) Is(target error) bool {
	other, ok := target.(*StoreNotFoundError)
	if !ok {
		return false
	}
	return e.Name == other.Name
}

type StoreExistsError struct {
	Name string
}

func (e *StoreExistsError) Error() string {
	return fmt.Sprintf("store exists: %q", e.Name)
}

func (e *StoreExistsError) Is(target error) bool {
	other, ok := target.(*StoreExistsError)
	if !ok {
		return false
	}
	return e.Name == other.Name
}

func NewStores() *Repository {
	return newStores(make(map[string]*streckuv1.Store))
}

func newStores(stores map[string]*streckuv1.Store) *Repository {
	return &Repository{
		stores: stores,
	}
}

func (r *Repository) LookupStore(_ context.Context, name string) (*streckuv1.Store, error) {
	store, ok := r.stores[name]
	if !ok {
		return nil, &StoreNotFoundError{Name: name}
	}
	return store, nil
}

func (r *Repository) ListStores(ctx context.Context) ([]*streckuv1.Store, error) {
	return r.FilterStores(ctx, func(*streckuv1.Store) bool { return true })
}

func (r *Repository) FilterStores(_ context.Context, predicate func(*streckuv1.Store) bool) ([]*streckuv1.Store, error) {
	var filtered []*streckuv1.Store
	for _, store := range r.stores {
		if predicate(store) {
			filtered = append(filtered, store)
		}
	}
	return filtered, nil
}

func (r *Repository) CreateStore(_ context.Context, store *streckuv1.Store) error {
	if _, exists := r.stores[store.Name]; exists {
		return &StoreExistsError{Name: store.Name}
	}
	r.stores[store.Name] = store
	return nil
}

func (r *Repository) UpdateStore(_ context.Context, updated *streckuv1.Store) error {
	if _, exists := r.stores[updated.Name]; !exists {
		return &StoreNotFoundError{Name: updated.Name}
	}
	r.stores[updated.Name] = updated
	return nil
}

func (r *Repository) DeleteStore(_ context.Context, name string) error {
	if _, exists := r.stores[name]; !exists {
		return &StoreNotFoundError{Name: name}
	}
	delete(r.stores, name)
	return nil
}
