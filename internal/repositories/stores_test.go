package repositories

import (
	"context"
	"testing"

	pb "github.com/Saser/strecku/api/v1"
	"github.com/Saser/strecku/resourcename"
	"github.com/Saser/strecku/resources/stores"
	"github.com/Saser/strecku/testresources"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/stretchr/testify/suite"
	"google.golang.org/protobuf/testing/protocmp"
)

func storeLess(u1, u2 *pb.Store) bool {
	return u1.Name < u2.Name
}

type StoresTestSuite struct {
	suite.Suite
	newStores func() Stores
}

func (s *StoresTestSuite) seedStores(ctx context.Context, t *testing.T, stores []*pb.Store) Stores {
	t.Helper()
	r := s.newStores()
	SeedStores(ctx, t, r, stores)
	return r
}

func (s *StoresTestSuite) seedBar(ctx context.Context, t *testing.T) Stores {
	t.Helper()
	return s.seedStores(
		ctx,
		t,
		[]*pb.Store{
			testresources.Bar,
		},
	)
}

func (s *StoresTestSuite) seedBarMall(ctx context.Context, t *testing.T) Stores {
	t.Helper()
	return s.seedStores(
		ctx,
		t,
		[]*pb.Store{
			testresources.Bar,
			testresources.Mall,
		},
	)
}

func (s *StoresTestSuite) seedBarMallPharmacy(ctx context.Context, t *testing.T) Stores {
	t.Helper()
	return s.seedStores(
		ctx,
		t,
		[]*pb.Store{
			testresources.Bar,
			testresources.Mall,
			testresources.Pharmacy,
		},
	)
}

func (s *StoresTestSuite) TestLookup() {
	t := s.T()
	ctx := context.Background()
	r := s.seedBar(ctx, t)
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
			wantErr:   resourcename.ErrInvalidName,
		},
		{
			desc:      "NotFound",
			name:      testresources.Pharmacy.Name,
			wantStore: nil,
			wantErr:   &NotFound{Name: testresources.Pharmacy.Name},
		},
	} {
		t.Run(test.desc, func(t *testing.T) {
			store, err := r.Lookup(ctx, test.name)
			if diff := cmp.Diff(store, test.wantStore, protocmp.Transform()); diff != "" {
				t.Errorf("r.Lookup(%v, %q) store != test.wantStore (-got +want)\n%s", ctx, test.name, diff)
			}
			if got, want := err, test.wantErr; !cmp.Equal(got, want, cmpopts.EquateErrors()) {
				t.Errorf("r.Lookup(%v, %q) err = %v; want %v", ctx, test.name, got, want)
			}
		})
	}
}

func (s *StoresTestSuite) TestList() {
	t := s.T()
	ctx := context.Background()
	r := s.seedBarMallPharmacy(ctx, t)
	want := []*pb.Store{
		testresources.Bar,
		testresources.Mall,
		testresources.Pharmacy,
	}
	stores, err := r.List(ctx)
	if diff := cmp.Diff(
		stores, want, protocmp.Transform(),
		cmpopts.SortSlices(storeLess),
	); diff != "" {
		t.Errorf("r.List(%v) stores != want (-got +want)\n%s", ctx, diff)
	}
	if err != nil {
		t.Errorf("r.List(%v) err = %v; want nil", ctx, err)
	}
}

func (s *StoresTestSuite) TestCreate() {
	t := s.T()
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
			want:  &Exists{Name: testresources.Bar.Name},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			r := s.seedBar(ctx, t)
			if got := r.Create(ctx, test.store); !cmp.Equal(got, test.want) {
				t.Errorf("r.Create(%v, %v) = %v; want %v", ctx, test.store, got, test.want)
			}
		})
	}
}

func (s *StoresTestSuite) TestUpdate() {
	t := s.T()
	ctx := context.Background()
	// Test scenario where the update is successful.
	t.Run("OK", func(t *testing.T) {
		r := s.seedBar(ctx, t)
		oldBar := stores.Clone(testresources.Bar)
		newBar := stores.Clone(oldBar)
		newBar.DisplayName = "New Bar"
		if err := r.Update(ctx, newBar); err != nil {
			t.Errorf("r.Update(%v, %v) = %v; want nil", ctx, newBar, err)
		}
		store, err := r.Lookup(ctx, newBar.Name)
		if diff := cmp.Diff(store, newBar, protocmp.Transform()); diff != "" {
			t.Errorf("r.Lookup(%v, %q) store != newBar (-got +want)\n%s", ctx, newBar.Name, diff)
		}
		if err != nil {
			t.Errorf("r.Lookup(%v, %q) err = %v; want nil", ctx, newBar.Name, err)
		}
	})

	// Test scenario where the update fails.
	t.Run("Errors", func(t *testing.T) {
		r := s.seedBar(ctx, t)
		for _, test := range []struct {
			desc   string
			modify func(bar *pb.Store)
			want   error
		}{
			{
				desc:   "EmptyName",
				modify: func(bar *pb.Store) { bar.Name = "" },
				want:   resourcename.ErrInvalidName,
			},
			{
				desc:   "EmptyDisplayName",
				modify: func(bar *pb.Store) { bar.DisplayName = "" },
				want:   stores.ErrDisplayNameEmpty,
			},
			{
				desc:   "NotFound",
				modify: func(bar *pb.Store) { bar.Name = testresources.Mall.Name },
				want:   &NotFound{Name: testresources.Mall.Name},
			},
		} {
			t.Run(test.desc, func(t *testing.T) {
				updated := stores.Clone(testresources.Bar)
				test.modify(updated)
				if got := r.Update(ctx, updated); !cmp.Equal(got, test.want, cmpopts.EquateErrors()) {
					t.Errorf("r.Update(%v, %v) = %v; want %v", ctx, updated, got, test.want)
				}
			})
		}
	})
}

func (s *StoresTestSuite) TestDelete() {
	t := s.T()
	ctx := context.Background()

	t.Run("OK", func(t *testing.T) {
		r := s.seedBarMall(ctx, t)
		if err := r.Delete(ctx, testresources.Bar.Name); err != nil {
			t.Errorf("r.Delete(%v, %q) = %v; want nil", ctx, testresources.Bar.Name, err)
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
				wantErr:   &NotFound{Name: testresources.Bar.Name},
			},
			{
				desc:      "LookupExisting",
				name:      testresources.Mall.Name,
				wantStore: testresources.Mall,
				wantErr:   nil,
			},
		} {
			t.Run(test.desc, func(t *testing.T) {
				store, err := r.Lookup(ctx, test.name)
				if diff := cmp.Diff(store, test.wantStore, protocmp.Transform()); diff != "" {
					t.Errorf("r.Lookup(%v, %q) store != test.wantStore (-got +want)\n%s", ctx, test.name, diff)
				}
				if !cmp.Equal(err, test.wantErr, cmpopts.EquateErrors()) {
					t.Errorf("r.Lookup(%v, %q) err = %v; want %v", ctx, test.name, err, test.wantErr)
				}
			})
		}
	})

	t.Run("Errors", func(t *testing.T) {
		r := s.seedBar(ctx, t)
		for _, test := range []struct {
			desc string
			name string
			want error
		}{
			{
				desc: "EmptyName",
				name: "",
				want: resourcename.ErrInvalidName,
			},
			{
				desc: "NotFound",
				name: testresources.Mall.Name,
				want: &NotFound{Name: testresources.Mall.Name},
			},
		} {
			t.Run(test.desc, func(t *testing.T) {
				if got := r.Delete(ctx, test.name); !cmp.Equal(got, test.want, cmpopts.EquateErrors()) {
					t.Errorf("r.Delete(%v, %q) = %v; want %v", ctx, test.name, got, test.want)
				}
			})
		}
	})
}
