package payments

import (
	"context"
	"errors"
	"fmt"
	"testing"

	pb "github.com/Saser/strecku/api/v1"
)

var ErrUpdateUser = errors.New("user cannot be updated")

type NotFoundError struct {
	Name string
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("payment not found: %q", e.Name)
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
	return fmt.Sprintf("payment exists: %q", e.Name)
}

func (e *ExistsError) Is(target error) bool {
	other, ok := target.(*ExistsError)
	if !ok {
		return false
	}
	return other.Name == e.Name
}

type Repository struct {
	payments map[string]*pb.Payment // name -> payment
}

func NewRepository() *Repository {
	return newRepository(make(map[string]*pb.Payment))
}

func SeedRepository(t *testing.T, payments []*pb.Payment) *Repository {
	mPayments := make(map[string]*pb.Payment)
	for _, payment := range payments {
		if err := Validate(payment); err != nil {
			t.Errorf("Validate(%v) = %v; want nil", payment, err)
		}
		mPayments[payment.Name] = payment
	}
	if t.Failed() {
		t.FailNow()
	}
	return newRepository(mPayments)
}

func newRepository(payments map[string]*pb.Payment) *Repository {
	return &Repository{
		payments: payments,
	}
}

func (r *Repository) LookupPayment(_ context.Context, name string) (*pb.Payment, error) {
	if err := ValidateName(name); err != nil {
		return nil, err
	}
	payment, ok := r.payments[name]
	if !ok {
		return nil, &NotFoundError{Name: name}
	}
	return payment, nil
}

func (r *Repository) ListPayments(ctx context.Context) ([]*pb.Payment, error) {
	return r.FilterPayments(ctx, func(*pb.Payment) bool { return true })
}

func (r *Repository) FilterPayments(_ context.Context, predicate func(*pb.Payment) bool) ([]*pb.Payment, error) {
	var filtered []*pb.Payment
	for _, payment := range r.payments {
		if predicate(payment) {
			filtered = append(filtered, payment)
		}
	}
	return filtered, nil
}

func (r *Repository) CreatePayment(_ context.Context, payment *pb.Payment) error {
	if err := Validate(payment); err != nil {
		return err
	}
	name := payment.Name
	if _, exists := r.payments[name]; exists {
		return &ExistsError{Name: name}
	}
	r.payments[name] = payment
	return nil
}

func (r *Repository) UpdatePayment(_ context.Context, updated *pb.Payment) error {
	if err := Validate(updated); err != nil {
		return err
	}
	name := updated.Name
	payment, ok := r.payments[name]
	if !ok {
		return &NotFoundError{Name: name}
	}
	if updated.User != payment.User {
		return ErrUpdateUser
	}
	r.payments[name] = updated
	return nil
}

func (r *Repository) DeletePayment(_ context.Context, name string) error {
	if err := ValidateName(name); err != nil {
		return err
	}
	if _, exists := r.payments[name]; !exists {
		return &NotFoundError{Name: name}
	}
	delete(r.payments, name)
	return nil
}
