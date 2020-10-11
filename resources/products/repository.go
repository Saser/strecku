package products

import (
	"context"
	"errors"
	"fmt"
	"testing"

	pb "github.com/Saser/strecku/api/v1"
	"google.golang.org/protobuf/proto"
)

var ErrUpdateParent = errors.New("parent cannot be updated")

type NotFoundError struct {
	Name string
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("product not found: %q", e.Name)
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
	return fmt.Sprintf("product exists: %q", e.Name)
}

func (e *ExistsError) Is(target error) bool {
	other, ok := target.(*ExistsError)
	if !ok {
		return false
	}
	return other.Name == e.Name
}

func Clone(product *pb.Product) *pb.Product {
	return proto.Clone(product).(*pb.Product)
}

type Repository struct {
	products map[string]*pb.Product // name -> product
}

func NewRepository() *Repository {
	return newRepository(make(map[string]*pb.Product))
}

func SeedRepository(t *testing.T, products []*pb.Product) *Repository {
	mProducts := make(map[string]*pb.Product)
	for _, product := range products {
		if err := Validate(product); err != nil {
			t.Errorf("Validate(%v) = %v; want nil", product, err)
		}
		mProducts[product.Name] = product
	}
	if t.Failed() {
		t.FailNow()
	}
	return newRepository(mProducts)
}

func newRepository(products map[string]*pb.Product) *Repository {
	return &Repository{
		products: products,
	}
}

func (r *Repository) LookupProduct(ctx context.Context, name string) (*pb.Product, error) {
	if err := ValidateName(name); err != nil {
		return nil, err
	}
	product, ok := r.products[name]
	if !ok {
		return nil, &NotFoundError{Name: name}
	}
	return product, nil
}

func (r *Repository) ListProducts(ctx context.Context) ([]*pb.Product, error) {
	return r.FilterProducts(ctx, func(*pb.Product) bool { return true })
}

func (r *Repository) FilterProducts(ctx context.Context, predicate func(*pb.Product) bool) ([]*pb.Product, error) {
	var filtered []*pb.Product
	for _, product := range r.products {
		if predicate(product) {
			filtered = append(filtered, product)
		}
	}
	return filtered, nil
}

func (r *Repository) CreateProduct(ctx context.Context, product *pb.Product) error {
	if err := Validate(product); err != nil {
		return err
	}
	name := product.Name
	if _, exists := r.products[name]; exists {
		return &ExistsError{Name: name}
	}
	r.products[name] = product
	return nil
}

func (r *Repository) UpdateProduct(ctx context.Context, updated *pb.Product) error {
	if err := Validate(updated); err != nil {
		return err
	}
	name := updated.Name
	product, ok := r.products[name]
	if !ok {
		return &NotFoundError{Name: name}
	}
	if updated.Parent != product.Parent {
		return ErrUpdateParent
	}
	r.products[name] = updated
	return nil
}

func (r *Repository) DeleteProduct(ctx context.Context, name string) error {
	if err := ValidateName(name); err != nil {
		return err
	}
	if _, exists := r.products[name]; !exists {
		return &NotFoundError{Name: name}
	}
	delete(r.products, name)
	return nil
}
