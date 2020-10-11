package products

import (
	"fmt"
	"testing"

	"github.com/Saser/strecku/resources/products/testproducts"
	"github.com/google/go-cmp/cmp"
)

func TestNotFoundError_Error(t *testing.T) {
	err := &NotFoundError{Name: testproducts.Bar_Beer.Name}
	want := fmt.Sprintf("product not found: %q", testproducts.Bar_Beer.Name)
	if got := err.Error(); !cmp.Equal(got, want) {
		t.Errorf("err.Error() = %q; want %q", got, want)
	}
}

func TestNotFoundError_Is(t *testing.T) {
	for _, test := range []struct {
		err    *NotFoundError
		target error
		want   bool
	}{
		{
			err:    &NotFoundError{Name: testproducts.Bar_Beer.Name},
			target: &NotFoundError{Name: testproducts.Bar_Beer.Name},
			want:   true,
		},
		{
			err:    &NotFoundError{Name: testproducts.Bar_Beer.Name},
			target: &NotFoundError{Name: testproducts.Bar_Cocktail.Name},
			want:   false,
		},
		{
			err:    &NotFoundError{Name: testproducts.Bar_Beer.Name},
			target: fmt.Errorf("product not found: %q", testproducts.Bar_Beer.Name),
			want:   false,
		},
	} {
		if got := test.err.Is(test.target); got != test.want {
			t.Errorf("test.err.Is(%v) = %v; want %v", test.target, got, test.want)
		}
	}
}

func TestExistsError_Error(t *testing.T) {
	err := &ExistsError{Name: testproducts.Bar_Beer.Name}
	want := fmt.Sprintf("product exists: %q", testproducts.Bar_Beer.Name)
	if got := err.Error(); !cmp.Equal(got, want) {
		t.Errorf("err.Error() = %q; want %q", got, want)
	}
}

func TestExistsError_Is(t *testing.T) {
	for _, test := range []struct {
		err    *ExistsError
		target error
		want   bool
	}{
		{
			err:    &ExistsError{Name: testproducts.Bar_Beer.Name},
			target: &ExistsError{Name: testproducts.Bar_Beer.Name},
			want:   true,
		},
		{
			err:    &ExistsError{Name: testproducts.Bar_Beer.Name},
			target: &ExistsError{Name: testproducts.Bar_Cocktail.Name},
			want:   false,
		},
		{
			err:    &ExistsError{Name: testproducts.Bar_Beer.Name},
			target: fmt.Errorf("product exists: %q", testproducts.Bar_Beer.Name),
			want:   false,
		},
	} {
		if got := test.err.Is(test.target); got != test.want {
			t.Errorf("test.err.Is(%v) = %v; want %v", test.target, got, test.want)
		}
	}
}
