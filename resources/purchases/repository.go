package purchases

import (
	"context"
	"errors"
	"fmt"
	"testing"

	pb "github.com/Saser/strecku/api/v1"
)

var (
	ErrUpdateUser  = errors.New("user cannot be updated")
	ErrUpdateStore = errors.New("store cannot be updated")
)

type NotFoundError struct {
	Name string
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("purchase not found: %q", e.Name)
}

func (e *NotFoundError) Is(target error) bool {
	other, ok := target.(*NotFoundError)
	if !ok {
		return false
	}
	return other.Name == e.Name
}

type ExistsError struct {
	Name string
}

func (e *ExistsError) Error() string {
	return fmt.Sprintf("purchase exists: %q", e.Name)
}

func (e *ExistsError) Is(target error) bool {
	other, ok := target.(*ExistsError)
	if !ok {
		return false
	}
	return other.Name == e.Name
}

type Repository struct {
	purchases map[string]*pb.Purchase // name -> purchase
}

func NewRepository() *Repository {
	return newRepository(make(map[string]*pb.Purchase))
}

func SeedRepository(t *testing.T, purchases []*pb.Purchase) *Repository {
	mPurchases := make(map[string]*pb.Purchase)
	for _, purchase := range purchases {
		if err := Validate(purchase); err != nil {
			t.Errorf("Validate(%v) = %v; want nil", purchase, err)
		}
		mPurchases[purchase.Name] = purchase
	}
	if t.Failed() {
		t.FailNow()
	}
	return newRepository(mPurchases)
}

func newRepository(purchases map[string]*pb.Purchase) *Repository {
	return &Repository{
		purchases: purchases,
	}
}

func (r *Repository) LookupPurchase(ctx context.Context, name string) (*pb.Purchase, error) {
	if err := ValidateName(name); err != nil {
		return nil, err
	}
	purchase, ok := r.purchases[name]
	if !ok {
		return nil, &NotFoundError{Name: name}
	}
	return purchase, nil
}

func (r *Repository) ListPurchases(ctx context.Context) ([]*pb.Purchase, error) {
	return r.FilterPurchases(ctx, func(*pb.Purchase) bool { return true })
}

func (r *Repository) FilterPurchases(ctx context.Context, predicate func(*pb.Purchase) bool) ([]*pb.Purchase, error) {
	var filtered []*pb.Purchase
	for _, purchase := range r.purchases {
		if predicate(purchase) {
			filtered = append(filtered, purchase)
		}
	}
	return filtered, nil
}

func (r *Repository) CreatePurchase(ctx context.Context, purchase *pb.Purchase) error {
	if err := Validate(purchase); err != nil {
		return err
	}
	name := purchase.Name
	if _, exists := r.purchases[name]; exists {
		return &ExistsError{Name: name}
	}
	r.purchases[name] = purchase
	return nil
}

func (r *Repository) UpdatePurchase(ctx context.Context, updated *pb.Purchase) error {
	if err := Validate(updated); err != nil {
		return err
	}
	name := updated.Name
	purchase, ok := r.purchases[name]
	if !ok {
		return &NotFoundError{Name: name}
	}
	if updated.User != purchase.User {
		return ErrUpdateUser
	}
	if updated.Store != purchase.Store {
		return ErrUpdateStore
	}
	r.purchases[name] = updated
	return nil
}

func (r *Repository) DeletePurchase(ctx context.Context, name string) error {
	if err := ValidateName(name); err != nil {
		return err
	}
	if _, exists := r.purchases[name]; !exists {
		return &NotFoundError{Name: name}
	}
	delete(r.purchases, name)
	return nil
}
