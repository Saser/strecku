package stores

import (
	"context"
	"fmt"
	"strings"
	"testing"

	streckuv1 "github.com/Saser/strecku/saser/strecku/v1"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"google.golang.org/protobuf/testing/protocmp"
)

var storeNames = map[string]string{
	"foobar": "stores/6f2d193c-1460-491d-8157-7dd9535526c6",
	"barbaz": "stores/d8bbf79e-8c59-4fae-aef9-634fcac00e07",
	"quux":   "stores/9cd3ec05-e7af-418c-bd50-80a7c39a18cc",
	"cookie": "stores/7f8f2c29-3860-49fa-923d-896a53f0ca26",
}

func storeLess(u1, u2 *streckuv1.Store) bool {
	return u1.Name < u2.Name
}

func seed(t *testing.T, stores []*streckuv1.Store) *Repository {
	t.Helper()
	mStores := make(map[string]*streckuv1.Store, len(stores))
	for _, store := range stores {
		if got := Validate(store); got != nil {
			t.Errorf("Validate(%v) = %v; want %v", store, got, nil)
		}
		mStores[store.Name] = store
	}
	return newStores(mStores)
}

func TestStoreNotFoundError_Error(t *testing.T) {
	for _, test := range []struct {
		name string
		want string
	}{
		{name: storeNames["foobar"], want: fmt.Sprintf("store not found: %q", storeNames["foobar"])},
		{name: "some name", want: `store not found: "some name"`},
	} {
		err := &StoreNotFoundError{
			Name: test.name,
		}
		if got := err.Error(); got != test.want {
			t.Errorf("err.Error() = %q; want %q", got, test.want)
		}
	}
}

func TestStoreNotFoundError_Is(t *testing.T) {
	for _, test := range []struct {
		err    *StoreNotFoundError
		target error
		want   bool
	}{
		{
			err:    &StoreNotFoundError{Name: storeNames["foobar"]},
			target: &StoreNotFoundError{Name: storeNames["foobar"]},
			want:   true,
		},
		{
			err:    &StoreNotFoundError{Name: storeNames["foobar"]},
			target: &StoreNotFoundError{Name: storeNames["barbaz"]},
			want:   false,
		},
		{
			err:    &StoreNotFoundError{Name: storeNames["foobar"]},
			target: fmt.Errorf("store not found: %q", storeNames["foobar"]),
			want:   false,
		},
	} {
		if got := test.err.Is(test.target); got != test.want {
			t.Errorf("test.err.Is(%v) = %v; want %v", test.target, got, test.want)
		}
	}
}

func TestStoreExistsError_Error(t *testing.T) {
	for _, test := range []struct {
		name string
		want string
	}{
		{name: storeNames["foobar"], want: fmt.Sprintf("store exists: %q", storeNames["foobar"])},
		{name: "some name", want: `store exists: "some name"`},
	} {
		err := &StoreExistsError{
			Name: test.name,
		}
		if got := err.Error(); got != test.want {
			t.Errorf("err.Error() = %q; want %q", got, test.want)
		}
	}
}

func TestStoreExistsError_Is(t *testing.T) {
	for _, test := range []struct {
		err    *StoreExistsError
		target error
		want   bool
	}{
		{
			err:    &StoreExistsError{Name: storeNames["foobar"]},
			target: &StoreExistsError{Name: storeNames["foobar"]},
			want:   true,
		},
		{
			err:    &StoreExistsError{Name: storeNames["foobar"]},
			target: &StoreExistsError{Name: storeNames["barbaz"]},
			want:   false,
		},
		{
			err:    &StoreExistsError{Name: storeNames["foobar"]},
			target: fmt.Errorf("store exists: %q", storeNames["foobar"]),
			want:   false,
		},
	} {
		if got := test.err.Is(test.target); got != test.want {
			t.Errorf("test.err.Is(%v) = %v; want %v", test.target, got, test.want)
		}
	}
}

func TestStores_LookupStore(t *testing.T) {
	ctx := context.Background()
	for _, test := range []struct {
		desc      string
		stores    []*streckuv1.Store
		name      string
		wantStore *streckuv1.Store
		wantErr   error
	}{
		{
			desc:      "EmptyDatabaseEmptyName",
			stores:    nil,
			name:      "",
			wantStore: nil,
			wantErr:   &StoreNotFoundError{Name: ""},
		},
		{
			desc:      "EmptyDatabaseNonEmptyName",
			stores:    nil,
			name:      storeNames["foobar"],
			wantStore: nil,
			wantErr:   &StoreNotFoundError{Name: storeNames["foobar"]},
		},
		{
			desc: "OneStoreOK",
			stores: []*streckuv1.Store{
				{Name: storeNames["foobar"], DisplayName: "Store"},
			},
			name:      storeNames["foobar"],
			wantStore: &streckuv1.Store{Name: storeNames["foobar"], DisplayName: "Store"},
			wantErr:   nil,
		},
		{
			desc: "MultipleStoresOK",
			stores: []*streckuv1.Store{
				{Name: storeNames["foobar"], DisplayName: "Foo Bar"},
				{Name: storeNames["barbaz"], DisplayName: "Barba Z."},
				{Name: storeNames["quux"], DisplayName: "Qu Ux"},
			},
			name:      storeNames["barbaz"],
			wantStore: &streckuv1.Store{Name: storeNames["barbaz"], DisplayName: "Barba Z."},
			wantErr:   nil,
		},
		{
			desc: "OneStoreNotFound",
			stores: []*streckuv1.Store{
				{Name: storeNames["foobar"], DisplayName: "Store"},
			},
			name:      storeNames["barbaz"],
			wantStore: nil,
			wantErr:   &StoreNotFoundError{Name: storeNames["barbaz"]},
		},
		{
			desc: "MultipleStoresNotFound",
			stores: []*streckuv1.Store{
				{Name: storeNames["foobar"], DisplayName: "Foo Bar"},
				{Name: storeNames["barbaz"], DisplayName: "Barba Z."},
				{Name: storeNames["quux"], DisplayName: "Qu Ux"},
			},
			name:      storeNames["cookie"],
			wantStore: nil,
			wantErr:   &StoreNotFoundError{Name: storeNames["cookie"]},
		},
	} {
		t.Run(test.desc, func(t *testing.T) {
			r := seed(t, test.stores)
			store, err := r.LookupStore(ctx, test.name)
			if diff := cmp.Diff(store, test.wantStore, protocmp.Transform()); diff != "" {
				t.Errorf("r.LookupStore(%v, %q) store != test.wantStore (-got +want)\n%s", ctx, test.name, diff)
			}
			if got, want := err, test.wantErr; !cmp.Equal(got, want) {
				t.Errorf("r.LookupStore(%v, %q) err = %v; want %v", ctx, test.name, got, want)
			}
		})
	}
}

func TestStores_ListStores(t *testing.T) {
	ctx := context.Background()
	for _, test := range []struct {
		name   string
		stores []*streckuv1.Store
	}{
		{name: "Empty", stores: nil},
		{
			name: "OneStore",
			stores: []*streckuv1.Store{
				{Name: storeNames["foobar"], DisplayName: "Foo Bar"},
			},
		},
		{
			name: "ThreeStores",
			stores: []*streckuv1.Store{
				{Name: storeNames["foobar"], DisplayName: "Foo Bar"},
				{Name: storeNames["barbaz"], DisplayName: "Barba Z."},
				{Name: storeNames["quux"], DisplayName: "Qu Ux"},
			},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			r := seed(t, test.stores)
			stores, err := r.ListStores(ctx)
			if diff := cmp.Diff(
				stores, test.stores, protocmp.Transform(),
				cmpopts.EquateEmpty(),
				cmpopts.SortSlices(storeLess),
			); diff != "" {
				t.Errorf("r.ListStores(%v) stores != test.stores (-got +want)\n%s", ctx, diff)
			}
			if got, want := err, error(nil); !cmp.Equal(got, want) {
				t.Errorf("r.ListStores(%v) err = %v; want %v", ctx, got, want)
			}
		})
	}
}

func TestStores_FilterStores(t *testing.T) {
	ctx := context.Background()
	for _, test := range []struct {
		name      string
		stores    []*streckuv1.Store
		predicate func(*streckuv1.Store) bool
		want      []*streckuv1.Store
	}{
		{
			name:      "Empty",
			stores:    nil,
			predicate: func(*streckuv1.Store) bool { return true },
			want:      nil,
		},
		{
			name: "OneStoreNoneMatching",
			stores: []*streckuv1.Store{
				{Name: storeNames["foobar"], DisplayName: "Foo Bar"},
			},
			predicate: func(store *streckuv1.Store) bool { return false },
			want:      nil,
		},
		{
			name: "MultipleStoresNoneMatching",
			stores: []*streckuv1.Store{
				{Name: storeNames["foobar"], DisplayName: "Foo Bar"},
				{Name: storeNames["barbaz"], DisplayName: "Barba Z."},
				{Name: storeNames["quux"], DisplayName: "Qu Ux"},
			},
			predicate: func(store *streckuv1.Store) bool { return false },
			want:      nil,
		},
		{
			name: "OneStoreOneMatching",
			stores: []*streckuv1.Store{
				{Name: storeNames["foobar"], DisplayName: "Foo Bar"},
			},
			predicate: func(store *streckuv1.Store) bool { return strings.HasPrefix(store.DisplayName, "Foo") },
			want: []*streckuv1.Store{
				{Name: storeNames["foobar"], DisplayName: "Foo Bar"},
			},
		},
		{
			name: "MultipleStoresOneMatching",
			stores: []*streckuv1.Store{
				{Name: storeNames["foobar"], DisplayName: "Foo Bar"},
				{Name: storeNames["barbaz"], DisplayName: "Barba Z."},
				{Name: storeNames["quux"], DisplayName: "Qu Ux"},
			},
			predicate: func(store *streckuv1.Store) bool { return strings.HasPrefix(store.DisplayName, "Foo") },
			want: []*streckuv1.Store{
				{Name: storeNames["foobar"], DisplayName: "Foo Bar"},
			},
		},
		{
			name: "MultipleStoresMultipleMatching",
			stores: []*streckuv1.Store{
				{Name: storeNames["foobar"], DisplayName: "Foo Bar"},
				{Name: storeNames["barbaz"], DisplayName: "Barba Z."},
				{Name: storeNames["quux"], DisplayName: "Qu Ux"},
			},
			predicate: func(store *streckuv1.Store) bool { return strings.Contains(store.DisplayName, "Bar") },
			want: []*streckuv1.Store{
				{Name: storeNames["foobar"], DisplayName: "Foo Bar"},
				{Name: storeNames["barbaz"], DisplayName: "Barba Z."},
			},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			r := seed(t, test.stores)
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

func TestStores_CreateStore(t *testing.T) {
	ctx := context.Background()
	for _, test := range []struct {
		name   string
		stores []*streckuv1.Store
		store  *streckuv1.Store
		want   error
	}{
		{
			name:   "Empty",
			stores: nil,
			store:  &streckuv1.Store{Name: storeNames["foobar"], DisplayName: "Foo Bar"},
			want:   nil,
		},
		{
			name: "OneStoreOK",
			stores: []*streckuv1.Store{
				{Name: storeNames["foobar"], DisplayName: "Foo Bar"},
			},
			store: &streckuv1.Store{Name: storeNames["barbaz"], DisplayName: "Foo Bar"},
			want:  nil,
		},
		{
			name: "MultipleStoresOK",
			stores: []*streckuv1.Store{
				{Name: storeNames["foobar"], DisplayName: "Foo Bar"},
				{Name: storeNames["barbaz"], DisplayName: "Barba Z."},
				{Name: storeNames["quux"], DisplayName: "Qu Ux"},
			},
			store: &streckuv1.Store{Name: storeNames["cookie"], DisplayName: "Cookie"},
			want:  nil,
		},
		{
			name: "OneStoreDuplicateName",
			stores: []*streckuv1.Store{
				{Name: storeNames["foobar"], DisplayName: "Foo Bar"},
			},
			store: &streckuv1.Store{Name: storeNames["foobar"], DisplayName: "New Foo Bar"},
			want:  &StoreExistsError{Name: storeNames["foobar"]},
		},
		{
			name: "MultipleStoresDuplicateName",
			stores: []*streckuv1.Store{
				{Name: storeNames["foobar"], DisplayName: "Foo Bar"},
				{Name: storeNames["barbaz"], DisplayName: "Barba Z."},
				{Name: storeNames["quux"], DisplayName: "Qu Ux"},
			},
			store: &streckuv1.Store{Name: storeNames["foobar"], DisplayName: "Another Foo Bar"},
			want:  &StoreExistsError{Name: storeNames["foobar"]},
		},
		{
			name: "DuplicateDisplayName",
			stores: []*streckuv1.Store{
				{Name: storeNames["foobar"], DisplayName: "Foo Bar"},
			},
			store: &streckuv1.Store{Name: storeNames["barbaz"], DisplayName: "Foo Bar"},
			want:  nil,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			r := seed(t, test.stores)
			if got := r.CreateStore(ctx, test.store); !cmp.Equal(got, test.want) {
				t.Errorf("r.CreateStore(%v, %v) = %v; want %v", ctx, test.store, got, test.want)
			}
		})
	}
}

func TestStores_UpdateStore(t *testing.T) {
	ctx := context.Background()
	for _, test := range []struct {
		name          string
		stores        []*streckuv1.Store
		updated       *streckuv1.Store
		wantUpdateErr error
		lookupName    string
		wantStore     *streckuv1.Store
		wantLookupErr error
	}{
		{
			name: "OneStoreDisplayName",
			stores: []*streckuv1.Store{
				{Name: storeNames["foobar"], DisplayName: "Foo Bar"},
			},
			updated:       &streckuv1.Store{Name: storeNames["foobar"], DisplayName: "New Foo Bar"},
			wantUpdateErr: nil,
			lookupName:    storeNames["foobar"],
			wantStore:     &streckuv1.Store{Name: storeNames["foobar"], DisplayName: "New Foo Bar"},
			wantLookupErr: nil,
		},
		{
			name: "OneStoreNotFound",
			stores: []*streckuv1.Store{
				{Name: storeNames["foobar"], DisplayName: "Foo Bar"},
			},
			updated:       &streckuv1.Store{Name: storeNames["barbaz"], DisplayName: "Foo Bar"},
			wantUpdateErr: &StoreNotFoundError{Name: storeNames["barbaz"]},
			lookupName:    storeNames["barbaz"],
			wantStore:     nil,
			wantLookupErr: &StoreNotFoundError{Name: storeNames["barbaz"]},
		},
		{
			name: "MultipleStoresDisplayName",
			stores: []*streckuv1.Store{
				{Name: storeNames["foobar"], DisplayName: "Foo Bar"},
				{Name: storeNames["barbaz"], DisplayName: "Barba Z."},
				{Name: storeNames["quux"], DisplayName: "Qu Ux"},
			},
			updated:       &streckuv1.Store{Name: storeNames["barbaz"], DisplayName: "All New Barba Z."},
			wantUpdateErr: nil,
			lookupName:    storeNames["barbaz"],
			wantStore:     &streckuv1.Store{Name: storeNames["barbaz"], DisplayName: "All New Barba Z."},
			wantLookupErr: nil,
		},
		{
			name: "MultipleStoresNotFound",
			stores: []*streckuv1.Store{
				{Name: storeNames["foobar"], DisplayName: "Foo Bar"},
				{Name: storeNames["barbaz"], DisplayName: "Barba Z."},
				{Name: storeNames["quux"], DisplayName: "Qu Ux"},
			},
			updated:       &streckuv1.Store{Name: storeNames["cookie"], DisplayName: "Barba Z."},
			wantUpdateErr: &StoreNotFoundError{Name: storeNames["cookie"]},
			lookupName:    storeNames["barbaz"],
			wantStore:     &streckuv1.Store{Name: storeNames["barbaz"], DisplayName: "Barba Z."},
			wantLookupErr: nil,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			r := seed(t, test.stores)
			if got := r.UpdateStore(ctx, test.updated); !cmp.Equal(got, test.wantUpdateErr) {
				t.Errorf("r.UpdateStore(%v, %v) = %v; want %v", ctx, test.updated, got, test.wantUpdateErr)
			}
			store, err := r.LookupStore(ctx, test.lookupName)
			if diff := cmp.Diff(store, test.wantStore, protocmp.Transform()); diff != "" {
				t.Errorf("r.LookupStore(%v, %q) store != test.wantStore (-got +want)\n%s", ctx, test.lookupName, diff)
			}
			if got, want := err, test.wantLookupErr; !cmp.Equal(got, want) {
				t.Errorf("r.LookupStore(%v, %q) err = %v; want %v", ctx, test.lookupName, got, want)
			}
		})
	}
}

func TestStores_DeleteStore(t *testing.T) {
	ctx := context.Background()
	for _, test := range []struct {
		desc          string
		stores        []*streckuv1.Store
		name          string
		want          error
		lookupName    string
		wantStore     *streckuv1.Store
		wantLookupErr error
	}{
		{
			desc:          "Empty",
			stores:        nil,
			name:          storeNames["foobar"],
			want:          &StoreNotFoundError{Name: storeNames["foobar"]},
			lookupName:    storeNames["barbaz"],
			wantStore:     nil,
			wantLookupErr: &StoreNotFoundError{Name: storeNames["barbaz"]},
		},
		{
			desc: "OneStoreOK",
			stores: []*streckuv1.Store{
				{Name: storeNames["foobar"], DisplayName: "Foo Bar"},
			},
			name:          storeNames["foobar"],
			want:          nil,
			lookupName:    storeNames["foobar"],
			wantStore:     nil,
			wantLookupErr: &StoreNotFoundError{Name: storeNames["foobar"]},
		},
		{
			desc: "MultipleStoresLookupDeleted",
			stores: []*streckuv1.Store{
				{Name: storeNames["foobar"], DisplayName: "Foo Bar"},
				{Name: storeNames["barbaz"], DisplayName: "Barba Z."},
				{Name: storeNames["quux"], DisplayName: "Qu Ux"},
			},
			name:          storeNames["barbaz"],
			want:          nil,
			lookupName:    storeNames["barbaz"],
			wantStore:     nil,
			wantLookupErr: &StoreNotFoundError{Name: storeNames["barbaz"]},
		},
		{
			desc: "MultipleStoresLookupExisting",
			stores: []*streckuv1.Store{
				{Name: storeNames["foobar"], DisplayName: "Foo Bar"},
				{Name: storeNames["barbaz"], DisplayName: "Barba Z."},
				{Name: storeNames["quux"], DisplayName: "Qu Ux"},
			},
			name:          storeNames["barbaz"],
			want:          nil,
			lookupName:    storeNames["foobar"],
			wantStore:     &streckuv1.Store{Name: storeNames["foobar"], DisplayName: "Foo Bar"},
			wantLookupErr: nil,
		},
		{
			desc: "OneStoreNotFound",
			stores: []*streckuv1.Store{
				{Name: storeNames["foobar"], DisplayName: "Foo Bar"},
			},
			name:          storeNames["barbaz"],
			want:          &StoreNotFoundError{Name: storeNames["barbaz"]},
			lookupName:    storeNames["foobar"],
			wantStore:     &streckuv1.Store{Name: storeNames["foobar"], DisplayName: "Foo Bar"},
			wantLookupErr: nil,
		},
		{
			desc: "MultipleStoresNotFound",
			stores: []*streckuv1.Store{
				{Name: storeNames["foobar"], DisplayName: "Foo Bar"},
				{Name: storeNames["barbaz"], DisplayName: "Barba Z."},
				{Name: storeNames["quux"], DisplayName: "Qu Ux"},
			},
			name:          storeNames["cookie"],
			want:          &StoreNotFoundError{Name: storeNames["cookie"]},
			lookupName:    storeNames["foobar"],
			wantStore:     &streckuv1.Store{Name: storeNames["foobar"], DisplayName: "Foo Bar"},
			wantLookupErr: nil,
		},
	} {
		t.Run(test.desc, func(t *testing.T) {
			r := seed(t, test.stores)
			err := r.DeleteStore(ctx, test.name)
			if got, want := err, test.want; !cmp.Equal(got, want) {
				t.Errorf("r.DeleteStore(%v, %q) = %v; want %v", ctx, test.name, got, want)
			}
			store, err := r.LookupStore(ctx, test.lookupName)
			if diff := cmp.Diff(store, test.wantStore, protocmp.Transform()); diff != "" {
				t.Errorf("r.LookupStore(%v, %q) store != test.wantStore (-got +want)\n%s", ctx, test.lookupName, diff)
			}
			if got, want := err, test.wantLookupErr; !cmp.Equal(got, want) {
				t.Errorf("r.LookupStore(%v, %q) err = %v; want %v", ctx, test.lookupName, got, want)
			}
		})
	}
}
