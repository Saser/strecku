package stores

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

func storeLess(u1, u2 *pb.Store) bool {
	return u1.Name < u2.Name
}

func seedBar(t *testing.T) *Repository {
	return SeedRepository(t, []*pb.Store{testresources.Bar})
}

func seedBarMall(t *testing.T) *Repository {
	return SeedRepository(
		t,
		[]*pb.Store{
			testresources.Bar,
			testresources.Mall,
		})
}

func seedBarMallPharmacy(t *testing.T) *Repository {
	return SeedRepository(
		t,
		[]*pb.Store{
			testresources.Bar,
			testresources.Mall,
			testresources.Pharmacy,
		})
}

func TestNotFoundError_Error(t *testing.T) {
	err := &NotFoundError{Name: testresources.Bar.Name}
	want := fmt.Sprintf("store not found: %q", testresources.Bar.Name)
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
			err:    &NotFoundError{Name: testresources.Bar.Name},
			target: &NotFoundError{Name: testresources.Bar.Name},
			want:   true,
		},
		{
			err:    &NotFoundError{Name: testresources.Bar.Name},
			target: &NotFoundError{Name: testresources.Pharmacy.Name},
			want:   false,
		},
		{
			err:    &NotFoundError{Name: testresources.Bar.Name},
			target: fmt.Errorf("store not found: %q", testresources.Bar.Name),
			want:   false,
		},
	} {
		if got := test.err.Is(test.target); got != test.want {
			t.Errorf("test.err.Is(%v) = %v; want %v", test.target, got, test.want)
		}
	}
}

func TestExistsError_Error(t *testing.T) {
	err := &ExistsError{Name: testresources.Bar.Name}
	want := fmt.Sprintf("store exists: %q", testresources.Bar.Name)
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
			err:    &ExistsError{Name: testresources.Bar.Name},
			target: &ExistsError{Name: testresources.Bar.Name},
			want:   true,
		},
		{
			err:    &ExistsError{Name: testresources.Bar.Name},
			target: &ExistsError{Name: testresources.Pharmacy.Name},
			want:   false,
		},
		{
			err:    &ExistsError{Name: testresources.Bar.Name},
			target: fmt.Errorf("store exists: %q", testresources.Bar.Name),
			want:   false,
		},
	} {
		if got := test.err.Is(test.target); got != test.want {
			t.Errorf("test.err.Is(%v) = %v; want %v", test.target, got, test.want)
		}
	}
}

func TestRepository_LookupStore(t *testing.T) {
	ctx := context.Background()
	r := seedBar(t)
	for _, test := range []struct {
		desc      string
		name      string
		wantStore *pb.Store
		wantErr   error
	}{
		{
			desc:      "OK",
			name:      testresources.Bar.Name,
			wantStore: testresources.Bar,
			wantErr:   nil,
		},
		{
			desc:      "EmptyName",
			name:      "",
			wantStore: nil,
			wantErr:   ErrNameInvalidFormat,
		},
		{
			desc:      "NotFound",
			name:      testresources.Pharmacy.Name,
			wantStore: nil,
			wantErr:   &NotFoundError{Name: testresources.Pharmacy.Name},
		},
	} {
		t.Run(test.desc, func(t *testing.T) {
			store, err := r.LookupStore(ctx, test.name)
			if diff := cmp.Diff(store, test.wantStore, protocmp.Transform()); diff != "" {
				t.Errorf("r.LookupStore(%v, %q) store != test.wantStore (-got +want)\n%s", ctx, test.name, diff)
			}
			if got, want := err, test.wantErr; !cmp.Equal(got, want, cmpopts.EquateErrors()) {
				t.Errorf("r.LookupStore(%v, %q) err = %v; want %v", ctx, test.name, got, want)
			}
		})
	}
}

func TestRepository_ListStores(t *testing.T) {
	ctx := context.Background()
	r := seedBarMallPharmacy(t)
	want := []*pb.Store{
		testresources.Bar,
		testresources.Mall,
		testresources.Pharmacy,
	}
	stores, err := r.ListStores(ctx)
	if diff := cmp.Diff(
		stores, want, protocmp.Transform(),
		cmpopts.SortSlices(storeLess),
	); diff != "" {
		t.Errorf("r.ListStores(%v) stores != want (-got +want)\n%s", ctx, diff)
	}
	if err != nil {
		t.Errorf("r.ListStores(%v) err = %v; want nil", ctx, err)
	}
}

func TestRepository_FilterStores(t *testing.T) {
	ctx := context.Background()
	r := seedBarMallPharmacy(t)
	for _, test := range []struct {
		name      string
		predicate func(*pb.Store) bool
		want      []*pb.Store
	}{
		{
			name:      "NoneMatching",
			predicate: func(*pb.Store) bool { return false },
			want:      nil,
		},
		{
			name:      "OneMatching",
			predicate: func(store *pb.Store) bool { return store.Name == testresources.Bar.Name },
			want: []*pb.Store{
				testresources.Bar,
			},
		},
		{
			name: "MultipleMatching",
			predicate: func(store *pb.Store) bool {
				switch store.Name {
				case testresources.Bar.Name, testresources.Mall.Name:
					return true
				default:
					return false
				}
			},
			want: []*pb.Store{
				testresources.Bar,
				testresources.Mall,
			},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			stores, err := r.FilterStores(ctx, test.predicate)
			if diff := cmp.Diff(
				stores, test.want, protocmp.Transform(),
				cmpopts.EquateEmpty(),
				cmpopts.SortSlices(storeLess),
			); diff != "" {
				t.Errorf("r.FilterStores(%v, test.predicate) stores != test.want (-got +want)\n%s", ctx, diff)
			}
			if got, want := err, error(nil); !cmp.Equal(got, want) {
				t.Errorf("r.FilterStores(%v, test.predicate) err = %v; want %v", ctx, got, want)
			}
		})
	}
}

func TestRepository_CreateStore(t *testing.T) {
	ctx := context.Background()
	for _, test := range []struct {
		name  string
		store *pb.Store
		want  error
	}{
		{
			name:  "OneStoreOK",
			store: testresources.Mall,
			want:  nil,
		},
		{
			name:  "DuplicateName",
			store: &pb.Store{Name: testresources.Bar.Name, DisplayName: testresources.Mall.DisplayName},
			want:  &ExistsError{Name: testresources.Bar.Name},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			r := seedBar(t)
			if got := r.CreateStore(ctx, test.store); !cmp.Equal(got, test.want) {
				t.Errorf("r.CreateStore(%v, %v) = %v; want %v", ctx, test.store, got, test.want)
			}
		})
	}
}

func TestRepository_UpdateStore(t *testing.T) {
	ctx := context.Background()
	// Test scenario where the update is successful.
	t.Run("OK", func(t *testing.T) {
		r := seedBar(t)
		oldBar := Clone(testresources.Bar)
		newBar := Clone(oldBar)
		newBar.DisplayName = "New Bar"
		if err := r.UpdateStore(ctx, newBar); err != nil {
			t.Errorf("r.UpdateStore(%v, %v) = %v; want nil", ctx, newBar, err)
		}
		store, err := r.LookupStore(ctx, newBar.Name)
		if diff := cmp.Diff(store, newBar, protocmp.Transform()); diff != "" {
			t.Errorf("r.LookupStore(%v, %q) store != newBar (-got +want)\n%s", ctx, newBar.Name, diff)
		}
		if err != nil {
			t.Errorf("r.LookupStore(%v, %q) err = %v; want nil", ctx, newBar.Name, err)
		}
	})

	// Test scenario where the update fails.
	t.Run("Errors", func(t *testing.T) {
		r := seedBar(t)
		for _, test := range []struct {
			desc   string
			modify func(bar *pb.Store)
			want   error
		}{
			{
				desc:   "EmptyName",
				modify: func(bar *pb.Store) { bar.Name = "" },
				want:   ErrNameInvalidFormat,
			},
			{
				desc:   "EmptyDisplayName",
				modify: func(bar *pb.Store) { bar.DisplayName = "" },
				want:   ErrDisplayNameEmpty,
			},
			{
				desc:   "NotFound",
				modify: func(bar *pb.Store) { bar.Name = testresources.Mall.Name },
				want:   &NotFoundError{Name: testresources.Mall.Name},
			},
		} {
			t.Run(test.desc, func(t *testing.T) {
				updated := Clone(testresources.Bar)
				test.modify(updated)
				if got := r.UpdateStore(ctx, updated); !cmp.Equal(got, test.want, cmpopts.EquateErrors()) {
					t.Errorf("r.UpdateStore(%v, %v) = %v; want %v", ctx, updated, got, test.want)
				}
			})
		}
	})
}

func TestRepository_DeleteStore(t *testing.T) {
	ctx := context.Background()

	t.Run("OK", func(t *testing.T) {
		r := seedBarMall(t)
		if err := r.DeleteStore(ctx, testresources.Bar.Name); err != nil {
			t.Errorf("r.DeleteStore(%v, %q) = %v; want nil", ctx, testresources.Bar.Name, err)
		}
		for _, test := range []struct {
			desc      string
			name      string
			wantStore *pb.Store
			wantErr   error
		}{
			{
				desc:      "LookupDeleted",
				name:      testresources.Bar.Name,
				wantStore: nil,
				wantErr:   &NotFoundError{Name: testresources.Bar.Name},
			},
			{
				desc:      "LookupExisting",
				name:      testresources.Mall.Name,
				wantStore: testresources.Mall,
				wantErr:   nil,
			},
		} {
			t.Run(test.desc, func(t *testing.T) {
				store, err := r.LookupStore(ctx, test.name)
				if diff := cmp.Diff(store, test.wantStore, protocmp.Transform()); diff != "" {
					t.Errorf("r.LookupStore(%v, %q) store != test.wantStore (-got +want)\n%s", ctx, test.name, diff)
				}
				if !cmp.Equal(err, test.wantErr, cmpopts.EquateErrors()) {
					t.Errorf("r.LookupStore(%v, %q) err = %v; want %v", ctx, test.name, err, test.wantErr)
				}
			})
		}
	})
	t.Run("Errors", func(t *testing.T) {
		r := seedBar(t)
		for _, test := range []struct {
			desc string
			name string
			want error
		}{
			{
				desc: "EmptyName",
				name: "",
				want: ErrNameInvalidFormat,
			},
			{
				desc: "NotFound",
				name: testresources.Mall.Name,
				want: &NotFoundError{Name: testresources.Mall.Name},
			},
		} {
			t.Run(test.desc, func(t *testing.T) {
				if got := r.DeleteStore(ctx, test.name); !cmp.Equal(got, test.want, cmpopts.EquateErrors()) {
					t.Errorf("r.DeleteStore(%v, %q) = %v; want %v", ctx, test.name, got, test.want)
				}
			})
		}
	})
}
