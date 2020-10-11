package products

import (
	"fmt"
	"testing"

	"github.com/Saser/strecku/resources/products/testproducts"
	"github.com/google/go-cmp/cmp"
)

func TestProductNotFoundError_Error(t *testing.T) {
	err := &ProductNotFoundError{Name: testproducts.Bar_Beer.Name}
	want := fmt.Sprintf("product not found: %q", testproducts.Bar_Beer.Name)
	if got := err.Error(); !cmp.Equal(got, want) {
		t.Errorf("err.Error() = %q; want %q", got, want)
	}
}

func TestProductNotFoundError_Is(t *testing.T) {
	for _, test := range []struct {
		err    *ProductNotFoundError
		target error
		want   bool
	}{
		{
			err:    &ProductNotFoundError{Name: testproducts.Bar_Beer.Name},
			target: &ProductNotFoundError{Name: testproducts.Bar_Beer.Name},
			want:   true,
		},
		{
			err:    &ProductNotFoundError{Name: testproducts.Bar_Beer.Name},
			target: &ProductNotFoundError{Name: testproducts.Bar_Cocktail.Name},
			want:   false,
		},
		{
			err:    &ProductNotFoundError{Name: testproducts.Bar_Beer.Name},
			target: fmt.Errorf("product not found: %q", testproducts.Bar_Beer.Name),
			want:   false,
		},
	} {
		if got := test.err.Is(test.target); got != test.want {
			t.Errorf("test.err.Is(%v) = %v; want %v", test.target, got, test.want)
		}
	}
}

func TestProductExistsError_Error(t *testing.T) {
	err := &ProductExistsError{Name: testproducts.Bar_Beer.Name}
	want := fmt.Sprintf("product exists: %q", testproducts.Bar_Beer.Name)
	if got := err.Error(); !cmp.Equal(got, want) {
		t.Errorf("err.Error() = %q; want %q", got, want)
	}
}

func TestProductExistsError_Is(t *testing.T) {
	for _, test := range []struct {
		err    *ProductExistsError
		target error
		want   bool
	}{
		{
			err:    &ProductExistsError{Name: testproducts.Bar_Beer.Name},
			target: &ProductExistsError{Name: testproducts.Bar_Beer.Name},
			want:   true,
		},
		{
			err:    &ProductExistsError{Name: testproducts.Bar_Beer.Name},
			target: &ProductExistsError{Name: testproducts.Bar_Cocktail.Name},
			want:   false,
		},
		{
			err:    &ProductExistsError{Name: testproducts.Bar_Beer.Name},
			target: fmt.Errorf("product exists: %q", testproducts.Bar_Beer.Name),
			want:   false,
		},
	} {
		if got := test.err.Is(test.target); got != test.want {
			t.Errorf("test.err.Is(%v) = %v; want %v", test.target, got, test.want)
		}
	}
}
