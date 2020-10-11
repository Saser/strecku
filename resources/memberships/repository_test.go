package memberships

import (
	"context"
	"fmt"
	"testing"

	pb "github.com/Saser/strecku/api/v1"
	"github.com/Saser/strecku/resources/stores"
	"github.com/Saser/strecku/resources/testresources"
	"github.com/Saser/strecku/resources/users"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"google.golang.org/protobuf/testing/protocmp"
)

func membershipLess(m1, m2 *pb.Membership) bool {
	return m1.Name < m2.Name
}

func TestNotFoundError_Error(t *testing.T) {
	for _, test := range []struct {
		err  *NotFoundError
		want string
	}{
		{
			err: &NotFoundError{
				Name: testresources.Alice_Bar.Name,
			},
			want: fmt.Sprintf("membership not found: %q", testresources.Alice_Bar.Name),
		},
		{
			err: &NotFoundError{
				User:  testresources.Alice.Name,
				Store: testresources.Bar.Name,
			},
			want: fmt.Sprintf("membership not found: between %q and %q", testresources.Alice.Name, testresources.Bar.Name),
		},
	} {
		if got := test.err.Error(); got != test.want {
			t.Errorf("test.err.Error() = %q; want %q", got, test.want)
		}
	}
}

func TestNotFoundError_Is(t *testing.T) {
	for _, test := range []struct {
		err    *NotFoundError
		target error
		want   bool
	}{
		{
			err: &NotFoundError{
				Name: testresources.Alice_Bar.Name,
			},
			target: &NotFoundError{
				Name: testresources.Alice_Bar.Name,
			},
			want: true,
		},
		{
			err: &NotFoundError{
				Name: testresources.Alice_Bar.Name,
			},
			target: &NotFoundError{
				User:  testresources.Alice.Name,
				Store: testresources.Bar.Name,
			},
			want: false,
		},
		{
			err: &NotFoundError{
				User:  testresources.Alice.Name,
				Store: testresources.Bar.Name,
			},
			target: &NotFoundError{
				Name: testresources.Alice_Bar.Name,
			},
			want: false,
		},
		{
			err: &NotFoundError{
				User:  testresources.Alice.Name,
				Store: testresources.Bar.Name,
			},
			target: &NotFoundError{
				User:  testresources.Alice.Name,
				Store: testresources.Bar.Name,
			},
			want: true,
		},
		{
			err: &NotFoundError{
				User:  testresources.Alice.Name,
				Store: testresources.Bar.Name,
			},
			target: fmt.Errorf("membership not found: between %q and %q", testresources.Alice.Name, testresources.Bar.Name),
			want:   false,
		},
	} {
		if got := test.err.Is(test.target); got != test.want {
			t.Errorf("test.err.Is(%v) = %v; want %v", test.target, got, test.want)
		}
	}
}

func TestExistsError_Error(t *testing.T) {
	for _, test := range []struct {
		err  *ExistsError
		want string
	}{
		{
			err: &ExistsError{
				Name: testresources.Alice_Bar.Name,
			},
			want: fmt.Sprintf("membership exists: %q", testresources.Alice_Bar.Name),
		},
		{
			err: &ExistsError{
				User:  testresources.Alice.Name,
				Store: testresources.Bar.Name,
			},
			want: fmt.Sprintf("membership exists: between %q and %q", testresources.Alice.Name, testresources.Bar.Name),
		},
	} {
		if got := test.err.Error(); got != test.want {
			t.Errorf("test.err.Error() = %q; want %q", got, test.want)
		}
	}
}

func TestExistsError_Is(t *testing.T) {
	for _, test := range []struct {
		err    *ExistsError
		target error
		want   bool
	}{
		{
			err:    &ExistsError{Name: testresources.Alice_Bar.Name},
			target: &ExistsError{Name: testresources.Alice_Bar.Name},
			want:   true,
		},
		{
			err:    &ExistsError{Name: testresources.Alice_Bar.Name},
			target: &ExistsError{Name: testresources.Alice_Mall.Name},
			want:   false,
		},
		{
			err:    &ExistsError{Name: testresources.Alice_Bar.Name},
			target: &ExistsError{User: testresources.Alice.Name, Store: testresources.Bar.Name},
			want:   false,
		},
		{
			err:    &ExistsError{User: testresources.Alice.Name, Store: testresources.Bar.Name},
			target: &ExistsError{Name: testresources.Alice_Bar.Name},
			want:   false,
		},
		{
			err:    &ExistsError{User: testresources.Alice.Name, Store: testresources.Bar.Name},
			target: &ExistsError{User: testresources.Alice.Name, Store: testresources.Bar.Name},
			want:   true,
		},
		{
			err:    &ExistsError{User: testresources.Alice.Name, Store: testresources.Bar.Name},
			target: &ExistsError{User: testresources.Alice.Name, Store: testresources.Mall.Name},
			want:   false,
		},
		{
			err:    &ExistsError{User: testresources.Alice.Name, Store: testresources.Bar.Name},
			target: &ExistsError{User: testresources.Bob.Name, Store: testresources.Bar.Name},
			want:   false,
		},
		{
			err:    &ExistsError{User: testresources.Alice.Name, Store: testresources.Mall.Name},
			target: &ExistsError{User: testresources.Alice.Name, Store: testresources.Bar.Name},
			want:   false,
		},
		{
			err:    &ExistsError{User: testresources.Bob.Name, Store: testresources.Bar.Name},
			target: &ExistsError{User: testresources.Alice.Name, Store: testresources.Bar.Name},
			want:   false,
		},
		{
			err:    &ExistsError{Name: testresources.Alice_Bar.Name},
			target: fmt.Errorf("membership exists: %q", testresources.Alice_Bar.Name),
			want:   false,
		},
		{
			err:    &ExistsError{User: testresources.Alice.Name, Store: testresources.Bar.Name},
			target: fmt.Errorf("membership exists: between %q and %q", testresources.Alice.Name, testresources.Bar.Name),
			want:   false,
		},
	} {
		if got := test.err.Is(test.target); got != test.want {
			t.Errorf("test.err.Is(%v) = %v; want %v", test.target, got, test.want)
		}
	}
}

func TestRepository_LookupMembership(t *testing.T) {
	ctx := context.Background()
	r := SeedRepository(t, []*pb.Membership{testresources.Alice_Bar})
	for _, test := range []struct {
		desc           string
		name           string
		wantMembership *pb.Membership
		wantErr        error
	}{
		{
			desc:           "OK",
			name:           testresources.Alice_Bar.Name,
			wantMembership: testresources.Alice_Bar,
			wantErr:        nil,
		},
		{
			desc:           "EmptyName",
			name:           "",
			wantMembership: nil,
			wantErr:        ErrNameEmpty,
		},
		{
			desc:           "NotFound",
			name:           testresources.Alice_Mall.Name,
			wantMembership: nil,
			wantErr:        &NotFoundError{Name: testresources.Alice_Mall.Name},
		},
	} {
		t.Run(test.desc, func(t *testing.T) {
			membership, err := r.LookupMembership(ctx, test.name)
			if diff := cmp.Diff(membership, test.wantMembership, protocmp.Transform()); diff != "" {
				t.Errorf("r.LookupMembership(%v, %q) membership != test.wantMembership (-got +want)\n%s", ctx, test.name, diff)
			}
			if !cmp.Equal(err, test.wantErr, cmpopts.EquateErrors()) {
				t.Errorf("r.LookupMembership(%v, %q) err = %v; want %v", ctx, test.name, err, test.wantErr)
			}
		})
	}
}

func TestRepository_LookupMembershipBetween(t *testing.T) {
	ctx := context.Background()
	r := SeedRepository(t, []*pb.Membership{testresources.Alice_Bar})
	for _, test := range []struct {
		desc           string
		user           string
		store          string
		wantMembership *pb.Membership
		wantErr        error
	}{
		{
			desc:           "OK",
			user:           testresources.Alice.Name,
			store:          testresources.Bar.Name,
			wantMembership: testresources.Alice_Bar,
			wantErr:        nil,
		},
		{
			desc:           "EmptyUser",
			user:           "",
			store:          testresources.Bar.Name,
			wantMembership: nil,
			wantErr:        users.ErrNameEmpty,
		},
		{
			desc:           "EmptyStore",
			user:           testresources.Alice.Name,
			store:          "",
			wantMembership: nil,
			wantErr:        stores.ErrNameEmpty,
		},
		{
			desc:           "WrongUser",
			user:           testresources.Bob.Name,
			store:          testresources.Bar.Name,
			wantMembership: nil,
			wantErr:        &NotFoundError{User: testresources.Bob.Name, Store: testresources.Bar.Name},
		},
		{
			desc:           "WrongStore",
			user:           testresources.Alice.Name,
			store:          testresources.Mall.Name,
			wantMembership: nil,
			wantErr:        &NotFoundError{User: testresources.Alice.Name, Store: testresources.Mall.Name},
		},
	} {
		t.Run(test.desc, func(t *testing.T) {
			membership, err := r.LookupMembershipBetween(ctx, test.user, test.store)
			if diff := cmp.Diff(membership, test.wantMembership, protocmp.Transform()); diff != "" {
				t.Errorf("r.LookupMembership(%v, %q, %q) membership != test.wantMembership (-got +want)\n%s", ctx, test.user, test.store, diff)
			}
			if !cmp.Equal(err, test.wantErr, cmpopts.EquateErrors()) {
				t.Errorf("r.LookupMembership(%v, %q, %q) err = %v; want %v", ctx, test.user, test.store, err, test.wantErr)
			}
		})
	}
}

func TestRepository_ListMemberships(t *testing.T) {
	ctx := context.Background()
	allMemberships := []*pb.Membership{
		testresources.Alice_Bar,
		testresources.Alice_Mall,
		testresources.Bob_Bar,
		testresources.Bob_Mall,
	}
	r := SeedRepository(t, allMemberships)
	memberships, err := r.ListMemberships(ctx)
	if diff := cmp.Diff(memberships, allMemberships, protocmp.Transform(), cmpopts.SortSlices(membershipLess)); diff != "" {
		t.Errorf("r.ListMemberships(%v) memberships != allMemberships (-got +want)\n%s", ctx, diff)
	}
	if err != nil {
		t.Errorf("r.ListMemberships(%v) err = %v; want nil", ctx, err)
	}
}

func TestRepository_FilterMemberships(t *testing.T) {
	ctx := context.Background()
	r := SeedRepository(t, []*pb.Membership{
		testresources.Alice_Bar,
		testresources.Alice_Mall,
		testresources.Bob_Bar,
		testresources.Bob_Mall,
	})
	for _, test := range []struct {
		desc      string
		predicate func(*pb.Membership) bool
		want      []*pb.Membership
	}{
		{
			desc:      "NoneMatching",
			predicate: func(*pb.Membership) bool { return false },
			want:      nil,
		},
		{
			desc:      "OneMatching",
			predicate: func(membership *pb.Membership) bool { return membership.Name == testresources.Alice_Bar.Name },
			want: []*pb.Membership{
				testresources.Alice_Bar,
			},
		},
		{
			desc:      "MultipleMatching",
			predicate: func(membership *pb.Membership) bool { return membership.User == testresources.Alice.Name },
			want: []*pb.Membership{
				testresources.Alice_Bar,
				testresources.Alice_Mall,
			},
		},
	} {
		t.Run(test.desc, func(t *testing.T) {
			filtered, err := r.FilterMemberships(ctx, test.predicate)
			if diff := cmp.Diff(
				filtered, test.want, protocmp.Transform(),
				cmpopts.EquateEmpty(),
				cmpopts.SortSlices(membershipLess),
			); diff != "" {
				t.Errorf("r.FilterMemberships(%v, test.predicate) filtered != test.want (-got +want)\n%s", ctx, diff)
			}
			if err != nil {
				t.Errorf("r.FilterMemberships(%v) err = %v; want nil", ctx, err)
			}
		})
	}
}

func TestRepository_CreateMembership(t *testing.T) {
	ctx := context.Background()
	for _, test := range []struct {
		desc       string
		membership *pb.Membership
		want       error
	}{
		{
			desc:       "OK_SameUser",
			membership: testresources.Alice_Mall,
			want:       nil,
		},
		{
			desc:       "OK_SameStore",
			membership: testresources.Bob_Bar,
			want:       nil,
		},
		{
			desc: "DuplicateName",
			membership: &pb.Membership{
				Name:          testresources.Alice_Bar.Name,
				User:          testresources.Bob.Name,  // chosen arbitrarily
				Store:         testresources.Mall.Name, // chosen arbitrarily
				Administrator: false,
			},
			want: &ExistsError{Name: testresources.Alice_Bar.Name},
		},
		{
			desc: "DuplicateUserAndStore",
			membership: &pb.Membership{
				Name:          testresources.Bob_Mall.Name, // chosen arbitrarily
				User:          testresources.Alice.Name,
				Store:         testresources.Bar.Name,
				Administrator: false,
			},
			want: &ExistsError{
				User:  testresources.Alice.Name,
				Store: testresources.Bar.Name,
			},
		},
	} {
		t.Run(test.desc, func(t *testing.T) {
			r := SeedRepository(t, []*pb.Membership{testresources.Alice_Bar})
			if got := r.CreateMembership(ctx, test.membership); !cmp.Equal(got, test.want, cmpopts.EquateErrors()) {
				t.Errorf("r.CreateMembership(%v, %v) = %v; want %v", ctx, test.membership, got, test.want)
			}
		})
	}
}

func TestRepository_UpdateMembership(t *testing.T) {
	ctx := context.Background()
	t.Run("OK", func(t *testing.T) {
		for _, test := range []struct {
			desc   string
			modify func(aliceBar *pb.Membership)
			want   error
		}{
			{
				desc:   "NoOp",
				modify: func(aliceBar *pb.Membership) {},
				want:   nil,
			},
			{
				desc:   "UpdateAdministrator",
				modify: func(aliceBar *pb.Membership) { aliceBar.Administrator = true },
				want:   nil,
			},
		} {
			t.Run(test.desc, func(t *testing.T) {
				r := SeedRepository(t, []*pb.Membership{testresources.Alice_Bar})
				updated := Clone(testresources.Alice_Bar)
				test.modify(updated)
				if got := r.UpdateMembership(ctx, updated); !cmp.Equal(got, test.want, cmpopts.EquateErrors()) {
					t.Errorf("r.UpdateMembership(%v, %v) = %v; want %v", ctx, updated, got, test.want)
				}
				membership, err := r.LookupMembership(ctx, updated.Name)
				if diff := cmp.Diff(membership, updated, protocmp.Transform()); diff != "" {
					t.Errorf("r.LookupMembership(%v, %q) membership != updated (-got +want)\n%s", ctx, updated.Name, diff)
				}
				if err != nil {
					t.Errorf("r.LookupMembership(%v, %q) err = %v; want nil", ctx, updated.Name, err)
				}
			})
		}
	})
	t.Run("Errors", func(t *testing.T) {
		for _, test := range []struct {
			desc   string
			modify func(aliceBar *pb.Membership)
			want   error
		}{
			{
				desc:   "UpdateUser",
				modify: func(aliceBar *pb.Membership) { aliceBar.User = testresources.Bob.Name },
				want:   ErrUpdateUser,
			},
			{
				desc:   "UpdateStore",
				modify: func(aliceBar *pb.Membership) { aliceBar.Store = testresources.Mall.Name },
				want:   ErrUpdateStore,
			},
			{
				desc:   "NotFound",
				modify: func(aliceBar *pb.Membership) { aliceBar.Name = testresources.Alice_Mall.Name },
				want:   &NotFoundError{Name: testresources.Alice_Mall.Name},
			},
		} {
			t.Run(test.desc, func(t *testing.T) {
				r := SeedRepository(t, []*pb.Membership{testresources.Alice_Bar})
				oldAliceBar := testresources.Alice_Bar
				newAliceBar := Clone(oldAliceBar)
				test.modify(newAliceBar)
				if got := r.UpdateMembership(ctx, newAliceBar); !cmp.Equal(got, test.want, cmpopts.EquateErrors()) {
					t.Errorf("r.UpdateMembership(%v, %v) = %v; want %v", ctx, newAliceBar, got, test.want)
				}
				membership, err := r.LookupMembership(ctx, oldAliceBar.Name)
				if diff := cmp.Diff(membership, oldAliceBar, protocmp.Transform()); diff != "" {
					t.Errorf("r.LookupMembership(%v, %q) membership != testmemberships.Alice_Bar (-got +want)\n%s", ctx, oldAliceBar.Name, diff)
				}
				if err != nil {
					t.Errorf("r.LookupMembership(%v, %q) err = %v; want nil", ctx, oldAliceBar.Name, err)
				}
			})
		}
	})
}

func TestRepository_DeleteMembership(t *testing.T) {
	ctx := context.Background()
	// Test scenario where the delete succeeds.
	t.Run("OK", func(t *testing.T) {
		r := SeedRepository(t, []*pb.Membership{testresources.Alice_Bar})
		// First, delete the membership.
		if err := r.DeleteMembership(ctx, testresources.Alice_Bar.Name); err != nil {
			t.Errorf("r.DeleteMembership(%v, %q) = %v; want nil", ctx, testresources.Alice_Bar.Name, err)
		}
		// Then, verify that looking it up by name fails.
		wantMembership := (*pb.Membership)(nil)
		wantErr := &NotFoundError{Name: testresources.Alice_Bar.Name}
		membership, err := r.LookupMembership(ctx, testresources.Alice_Bar.Name)
		if diff := cmp.Diff(membership, wantMembership, protocmp.Transform()); diff != "" {
			t.Errorf("r.LookupMembership(%v, %q) membership != wantMembership (-got +want)\n%s", ctx, testresources.Alice_Bar.Name, diff)
		}
		if !cmp.Equal(err, wantErr, cmpopts.EquateErrors()) {
			t.Errorf("r.LookupMembership(%v, %q) err = %v; want %v", ctx, testresources.Alice_Bar.Name, err, wantErr)
		}
		// Finally, verify that looking it up by user and store fails also.
		wantMembership = (*pb.Membership)(nil)
		wantErr = &NotFoundError{
			User:  testresources.Alice.Name,
			Store: testresources.Bar.Name,
		}
		membership, err = r.LookupMembershipBetween(ctx, testresources.Alice.Name, testresources.Bar.Name)
		if diff := cmp.Diff(membership, wantMembership, protocmp.Transform()); diff != "" {
			t.Errorf("r.LookupMembershipBetween(%v, %q, %q) membership != wantMembership (-got +want)\n%s", ctx, testresources.Alice.Name, testresources.Bar.Name, diff)
		}
		if !cmp.Equal(err, wantErr, cmpopts.EquateErrors()) {
			t.Errorf("r.LookupMembershipBetween(%v, %q, %q) err = %v; want %v", ctx, testresources.Alice.Name, testresources.Bar.Name, err, wantErr)
		}
	})
	// Test scenarios where the delete fails.
	t.Run("Errors", func(t *testing.T) {
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
				name: testresources.Alice_Mall.Name,
				want: &NotFoundError{Name: testresources.Alice_Mall.Name},
			},
		} {
			t.Run(test.desc, func(t *testing.T) {
				r := SeedRepository(t, []*pb.Membership{testresources.Alice_Bar})
				// First, try and fail to delete the membership.
				if got := r.DeleteMembership(ctx, test.name); !cmp.Equal(got, test.want, cmpopts.EquateErrors()) {
					t.Errorf("r.DeleteMembership(%v, %q) = %v; want %v", ctx, test.name, got, test.want)
				}
				wantMembership := testresources.Alice_Bar
				// Then, verify that a lookup by name succeeds.
				membership, err := r.LookupMembership(ctx, testresources.Alice_Bar.Name)
				if diff := cmp.Diff(membership, wantMembership, protocmp.Transform()); diff != "" {
					t.Errorf("r.LookupMembership(%v, %q) membership != wantMembership (-got +want)\n%s", ctx, testresources.Alice_Bar.Name, diff)
				}
				if err != nil {
					t.Errorf("r.LookupMembership(%v, %q) err = %v; want nil", ctx, testresources.Alice_Bar.Name, err)
				}
				// Finally, verify that a lookup by user and store succeeds.
				membership, err = r.LookupMembershipBetween(ctx, testresources.Alice.Name, testresources.Bar.Name)
				if diff := cmp.Diff(membership, wantMembership, protocmp.Transform()); diff != "" {
					t.Errorf("r.LookupMembershipBetween(%v, %q, %q) membership != wantMembership (-got +want)\n%s", ctx, testresources.Alice.Name, testresources.Bar.Name, diff)
				}
				if err != nil {
					t.Errorf("r.LookupMembershipBetween(%v, %q, %q) err = %v; want nil", ctx, testresources.Alice.Name, testresources.Bar.Name, err)
				}
			})
		}
	})
}
