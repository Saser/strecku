package stores

import (
	"context"
	"fmt"
	"testing"

	pb "github.com/Saser/strecku/api/v1"
	"github.com/Saser/strecku/resources/stores/teststores"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"google.golang.org/protobuf/testing/protocmp"
)

func storeLess(u1, u2 *pb.Store) bool {
	return u1.Name < u2.Name
}

func seedBar(t *testing.T) *Repository {
	return SeedRepository(t, []*pb.Store{teststores.Bar})
}

func seedBarMall(t *testing.T) *Repository {
	return SeedRepository(
		t,
		[]*pb.Store{
			teststores.Bar,
			teststores.Mall,
		})
}

func seedBarMallPharmacy(t *testing.T) *Repository {
	return SeedRepository(
		t,
		[]*pb.Store{
			teststores.Bar,
			teststores.Mall,
			teststores.Pharmacy,
		})
}

func TestStoreNotFoundError_Error(t *testing.T) {
	err := &StoreNotFoundError{Name: teststores.Bar.Name}
	want := fmt.Sprintf("store not found: %q", teststores.Bar.Name)
	if got := err.Error(); !cmp.Equal(got, want) {
		t.Errorf("err.Error() = %q; want %q", got, want)
	}
}

func TestStoreNotFoundError_Is(t *testing.T) {
	for _, test := range []struct {
		err    *StoreNotFoundError
		target error
		want   bool
	}{
		{
			err:    &StoreNotFoundError{Name: teststores.Bar.Name},
			target: &StoreNotFoundError{Name: teststores.Bar.Name},
			want:   true,
		},
		{
			err:    &StoreNotFoundError{Name: teststores.Bar.Name},
			target: &StoreNotFoundError{Name: teststores.Pharmacy.Name},
			want:   false,
		},
		{
			err:    &StoreNotFoundError{Name: teststores.Bar.Name},
			target: fmt.Errorf("store not found: %q", teststores.Bar.Name),
			want:   false,
		},
	} {
		if got := test.err.Is(test.target); got != test.want {
			t.Errorf("test.err.Is(%v) = %v; want %v", test.target, got, test.want)
		}
	}
}

func TestStoreExistsError_Error(t *testing.T) {
	err := &StoreExistsError{Name: teststores.Bar.Name}
	want := fmt.Sprintf("store exists: %q", teststores.Bar.Name)
	if got := err.Error(); !cmp.Equal(got, want) {
		t.Errorf("err.Error() = %q; want %q", got, want)
	}
}

func TestStoreExistsError_Is(t *testing.T) {
	for _, test := range []struct {
		err    *StoreExistsError
		target error
		want   bool
	}{
		{
			err:    &StoreExistsError{Name: teststores.Bar.Name},
			target: &StoreExistsError{Name: teststores.Bar.Name},
			want:   true,
		},
		{
			err:    &StoreExistsError{Name: teststores.Bar.Name},
			target: &StoreExistsError{Name: teststores.Pharmacy.Name},
			want:   false,
		},
		{
			err:    &StoreExistsError{Name: teststores.Bar.Name},
			target: fmt.Errorf("store exists: %q", teststores.Bar.Name),
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
			name:      teststores.Bar.Name,
			wantStore: teststores.Bar,
			wantErr:   nil,
		},
		{
			desc:      "EmptyName",
			name:      "",
			wantStore: nil,
			wantErr:   ErrNameEmpty,
		},
		{
			desc:      "NotFound",
			name:      teststores.Pharmacy.Name,
			wantStore: nil,
			wantErr:   &StoreNotFoundError{Name: teststores.Pharmacy.Name},
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
		teststores.Bar,
		teststores.Mall,
		teststores.Pharmacy,
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
			predicate: func(store *pb.Store) bool { return store.Name == teststores.Bar.Name },
			want: []*pb.Store{
				teststores.Bar,
			},
		},
		{
			name: "MultipleMatching",
			predicate: func(store *pb.Store) bool {
				switch store.Name {
				case teststores.Bar.Name, teststores.Mall.Name:
					return true
				default:
					return false
				}
			},
			want: []*pb.Store{
				teststores.Bar,
				teststores.Mall,
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
			store: teststores.Mall,
			want:  nil,
		},
		{
			name:  "DuplicateName",
			store: &pb.Store{Name: teststores.Bar.Name, DisplayName: teststores.Mall.DisplayName},
			want:  &StoreExistsError{Name: teststores.Bar.Name},
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
		oldBar := Clone(teststores.Bar)
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
				want:   ErrNameEmpty,
			},
			{
				desc:   "EmptyDisplayName",
				modify: func(bar *pb.Store) { bar.DisplayName = "" },
				want:   ErrDisplayNameEmpty,
			},
			{
				desc:   "NotFound",
				modify: func(bar *pb.Store) { bar.Name = teststores.Mall.Name },
				want:   &StoreNotFoundError{Name: teststores.Mall.Name},
			},
		} {
			t.Run(test.desc, func(t *testing.T) {
				updated := Clone(teststores.Bar)
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
		if err := r.DeleteStore(ctx, teststores.Bar.Name); err != nil {
			t.Errorf("r.DeleteStore(%v, %q) = %v; want nil", ctx, teststores.Bar.Name, err)
		}
		for _, test := range []struct {
			desc      string
			name      string
			wantStore *pb.Store
			wantErr   error
		}{
			{
				desc:      "LookupDeleted",
				name:      teststores.Bar.Name,
				wantStore: nil,
				wantErr:   &StoreNotFoundError{Name: teststores.Bar.Name},
			},
			{
				desc:      "LookupExisting",
				name:      teststores.Mall.Name,
				wantStore: teststores.Mall,
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
				want: ErrNameEmpty,
			},
			{
				desc: "NotFound",
				name: teststores.Mall.Name,
				want: &StoreNotFoundError{Name: teststores.Mall.Name},
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
