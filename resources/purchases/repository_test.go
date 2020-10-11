package purchases

import (
	"fmt"
	"testing"

	"github.com/Saser/strecku/resources/testresources"
	"github.com/google/go-cmp/cmp"
)

func TestNotFoundError_Error(t *testing.T) {
	err := &NotFoundError{Name: testresources.Alice_Beer1.Name}
	want := fmt.Sprintf("purchase not found: %q", testresources.Alice_Beer1.Name)
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
			err:    &NotFoundError{Name: testresources.Alice_Beer1.Name},
			target: &NotFoundError{Name: testresources.Alice_Beer1.Name},
			want:   true,
		},
		{
			err:    &NotFoundError{Name: testresources.Alice_Beer1.Name},
			target: &NotFoundError{Name: testresources.Alice_Cocktail1.Name},
			want:   false,
		},
		{
			err:    &NotFoundError{Name: testresources.Alice_Beer1.Name},
			target: fmt.Errorf("purchase not found: %q", testresources.Alice_Beer1.Name),
			want:   false,
		},
	} {
		if got := test.err.Is(test.target); got != test.want {
			t.Errorf("test.err.Is(%v) = %v; want %v", test.target, got, test.want)
		}
	}
}

func TestExistsError_Error(t *testing.T) {
	err := &ExistsError{Name: testresources.Beer.Name}
	want := fmt.Sprintf("purchase exists: %q", testresources.Beer.Name)
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
			err:    &ExistsError{Name: testresources.Alice_Beer1.Name},
			target: &ExistsError{Name: testresources.Alice_Beer1.Name},
			want:   true,
		},
		{
			err:    &ExistsError{Name: testresources.Alice_Beer1.Name},
			target: &ExistsError{Name: testresources.Alice_Cocktail1.Name},
			want:   false,
		},
		{
			err:    &ExistsError{Name: testresources.Alice_Beer1.Name},
			target: fmt.Errorf("purchase exists: %q", testresources.Alice_Beer1.Name),
			want:   false,
		},
	} {
		if got := test.err.Is(test.target); got != test.want {
			t.Errorf("test.err.Is(%v) = %v; want %v", test.target, got, test.want)
		}
	}
}
