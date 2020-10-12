package purchases

import (
	"context"
	"fmt"
	"testing"

	pb "github.com/Saser/strecku/api/v1"
	"github.com/Saser/strecku/resources/testresources"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"google.golang.org/protobuf/testing/protocmp"
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

func TestRepository_LookupPurchase(t *testing.T) {
	ctx := context.Background()
	r := SeedRepository(t, []*pb.Purchase{testresources.Alice_Beer1})
	for _, test := range []struct {
		desc         string
		name         string
		wantPurchase *pb.Purchase
		wantErr      error
	}{
		{
			desc:         "OK",
			name:         testresources.Alice_Beer1.Name,
			wantPurchase: testresources.Alice_Beer1,
			wantErr:      nil,
		},
		{
			desc:         "EmptyName",
			name:         "",
			wantPurchase: nil,
			wantErr:      ErrNameEmpty,
		},
		{
			desc:         "NotFound",
			name:         testresources.Alice_Cocktail1.Name,
			wantPurchase: nil,
			wantErr:      &NotFoundError{Name: testresources.Alice_Cocktail1.Name},
		},
	} {
		t.Run(test.desc, func(t *testing.T) {
			purchase, err := r.LookupPurchase(ctx, test.name)
			if diff := cmp.Diff(purchase, test.wantPurchase, protocmp.Transform()); diff != "" {
				t.Errorf("r.LookupPurchase(%v, %q) purchase != test.wantPurchase (-got +want)\n%s", ctx, test.name, diff)
			}
			if !cmp.Equal(err, test.wantErr, cmpopts.EquateErrors()) {
				t.Errorf("r.LookupPurchase(%v, %q) err = %v; want %v", ctx, test.name, err, test.wantErr)
			}
		})
	}
}
