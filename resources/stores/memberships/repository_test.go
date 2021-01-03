package memberships

import (
	"context"
	"fmt"
	"testing"

	pb "github.com/Saser/strecku/api/v1"
	"github.com/Saser/strecku/resourcename"
	"github.com/Saser/strecku/testresources"
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
				Name: testresources.Bar_Alice.Name,
			},
			want: fmt.Sprintf("membership not found: %q", testresources.Bar_Alice.Name),
		},
		{
			err: &NotFoundError{
				User:   testresources.Alice.Name,
				Parent: testresources.Bar.Name,
			},
			want: fmt.Sprintf("membership not found: in %q for %q", testresources.Bar.Name, testresources.Alice.Name),
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
				Name: testresources.Bar_Alice.Name,
			},
			target: &NotFoundError{
				Name: testresources.Bar_Alice.Name,
			},
			want: true,
		},
		{
			err: &NotFoundError{
				Name: testresources.Bar_Alice.Name,
			},
			target: &NotFoundError{
				User:   testresources.Alice.Name,
				Parent: testresources.Bar.Name,
			},
			want: false,
		},
		{
			err: &NotFoundError{
				User:   testresources.Alice.Name,
				Parent: testresources.Bar.Name,
			},
			target: &NotFoundError{
				Name: testresources.Bar_Alice.Name,
			},
			want: false,
		},
		{
			err: &NotFoundError{
				User:   testresources.Alice.Name,
				Parent: testresources.Bar.Name,
			},
			target: &NotFoundError{
				User:   testresources.Alice.Name,
				Parent: testresources.Bar.Name,
			},
			want: true,
		},
		{
			err: &NotFoundError{
				User:   testresources.Alice.Name,
				Parent: testresources.Bar.Name,
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
				Name: testresources.Bar_Alice.Name,
			},
			want: fmt.Sprintf("membership exists: %q", testresources.Bar_Alice.Name),
		},
		{
			err: &ExistsError{
				User:   testresources.Alice.Name,
				Parent: testresources.Bar.Name,
			},
			want: fmt.Sprintf("membership exists: in %q for %q", testresources.Bar.Name, testresources.Alice.Name),
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
			err:    &ExistsError{Name: testresources.Bar_Alice.Name},
			target: &ExistsError{Name: testresources.Bar_Alice.Name},
			want:   true,
		},
		{
			err:    &ExistsError{Name: testresources.Bar_Alice.Name},
			target: &ExistsError{Name: testresources.Mall_Alice.Name},
			want:   false,
		},
		{
			err:    &ExistsError{Name: testresources.Bar_Alice.Name},
			target: &ExistsError{User: testresources.Alice.Name, Parent: testresources.Bar.Name},
			want:   false,
		},
		{
			err:    &ExistsError{User: testresources.Alice.Name, Parent: testresources.Bar.Name},
			target: &ExistsError{Name: testresources.Bar_Alice.Name},
			want:   false,
		},
		{
			err:    &ExistsError{User: testresources.Alice.Name, Parent: testresources.Bar.Name},
			target: &ExistsError{User: testresources.Alice.Name, Parent: testresources.Bar.Name},
			want:   true,
		},
		{
			err:    &ExistsError{User: testresources.Alice.Name, Parent: testresources.Bar.Name},
			target: &ExistsError{User: testresources.Alice.Name, Parent: testresources.Mall.Name},
			want:   false,
		},
		{
			err:    &ExistsError{User: testresources.Alice.Name, Parent: testresources.Bar.Name},
			target: &ExistsError{User: testresources.Bob.Name, Parent: testresources.Bar.Name},
			want:   false,
		},
		{
			err:    &ExistsError{User: testresources.Alice.Name, Parent: testresources.Mall.Name},
			target: &ExistsError{User: testresources.Alice.Name, Parent: testresources.Bar.Name},
			want:   false,
		},
		{
			err:    &ExistsError{User: testresources.Bob.Name, Parent: testresources.Bar.Name},
			target: &ExistsError{User: testresources.Alice.Name, Parent: testresources.Bar.Name},
			want:   false,
		},
		{
			err:    &ExistsError{Name: testresources.Bar_Alice.Name},
			target: fmt.Errorf("membership exists: %q", testresources.Bar_Alice.Name),
			want:   false,
		},
		{
			err:    &ExistsError{User: testresources.Alice.Name, Parent: testresources.Bar.Name},
			target: fmt.Errorf("membership exists: between %q and %q", testresources.Alice.Name, testresources.Bar.Name),
			want:   false,
		},
	} {
		if got := test.err.Is(test.target); got != test.want {
			t.Errorf("test.err.Is(%v) = %v; want %v", test.target, got, test.want)
		}
	}
}

func TestNewRepository(t *testing.T) {
	NewRepository()
}

func TestRepository_LookupMembership(t *testing.T) {
	ctx := context.Background()
	r := SeedRepository(t, []*pb.Membership{testresources.Bar_Alice})
	for _, test := range []struct {
		desc           string
		name           string
		wantMembership *pb.Membership
		wantErr        error
	}{
		{
			desc:           "OK",
			name:           testresources.Bar_Alice.Name,
			wantMembership: testresources.Bar_Alice,
			wantErr:        nil,
		},
		{
			desc:           "EmptyName",
			name:           "",
			wantMembership: nil,
			wantErr:        resourcename.ErrInvalidName,
		},
		{
			desc:           "NotFound",
			name:           testresources.Mall_Alice.Name,
			wantMembership: nil,
			wantErr:        &NotFoundError{Name: testresources.Mall_Alice.Name},
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
	r := SeedRepository(t, []*pb.Membership{testresources.Bar_Alice})
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
			wantMembership: testresources.Bar_Alice,
			wantErr:        nil,
		},
		{
			desc:           "EmptyUser",
			user:           "",
			store:          testresources.Bar.Name,
			wantMembership: nil,
			wantErr:        resourcename.ErrInvalidName,
		},
		{
			desc:           "EmptyStore",
			user:           testresources.Alice.Name,
			store:          "",
			wantMembership: nil,
			wantErr:        resourcename.ErrInvalidName,
		},
		{
			desc:           "WrongUser",
			user:           testresources.Bob.Name,
			store:          testresources.Bar.Name,
			wantMembership: nil,
			wantErr:        &NotFoundError{User: testresources.Bob.Name, Parent: testresources.Bar.Name},
		},
		{
			desc:           "WrongStore",
			user:           testresources.Alice.Name,
			store:          testresources.Mall.Name,
			wantMembership: nil,
			wantErr:        &NotFoundError{User: testresources.Alice.Name, Parent: testresources.Mall.Name},
		},
	} {
		t.Run(test.desc, func(t *testing.T) {
			membership, err := r.LookupMembershipIn(ctx, test.store, test.user)
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
		testresources.Bar_Alice,
		testresources.Mall_Alice,
		testresources.Bar_Bob,
		testresources.Mall_Bob,
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
		testresources.Bar_Alice,
		testresources.Mall_Alice,
		testresources.Bar_Bob,
		testresources.Mall_Bob,
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
			predicate: func(membership *pb.Membership) bool { return membership.Name == testresources.Bar_Alice.Name },
			want: []*pb.Membership{
				testresources.Bar_Alice,
			},
		},
		{
			desc:      "MultipleMatching",
			predicate: func(membership *pb.Membership) bool { return membership.User == testresources.Alice.Name },
			want: []*pb.Membership{
				testresources.Bar_Alice,
				testresources.Mall_Alice,
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
			desc:       "OK_SameParent",
			membership: testresources.Bar_Bob,
			want:       nil,
		},
		{
			desc:       "OK_SameUser",
			membership: testresources.Mall_Alice,
			want:       nil,
		},
		{
			desc: "DuplicateName",
			membership: &pb.Membership{
				Name:          testresources.Bar_Alice.Name,
				User:          testresources.Bob.Name, // chosen arbitrarily
				Administrator: false,
				Discount:      false,
			},
			want: &ExistsError{Name: testresources.Bar_Alice.Name},
		},
		{
			desc: "DuplicateUserAndStore",
			membership: &pb.Membership{
				Name:          testresources.Bar.Name + "/" + CollectionID + "/2a1f364b-1a1f-400b-a2da-aa2e14e40eae", // chosen arbitrarily
				User:          testresources.Alice.Name,
				Administrator: false,
				Discount:      false,
			},
			want: &ExistsError{
				Parent: testresources.Bar.Name,
				User:   testresources.Alice.Name,
			},
		},
	} {
		t.Run(test.desc, func(t *testing.T) {
			r := SeedRepository(t, []*pb.Membership{testresources.Bar_Alice})
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
				r := SeedRepository(t, []*pb.Membership{testresources.Bar_Alice})
				updated := Clone(testresources.Bar_Alice)
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
			modify func(barAlice *pb.Membership)
			want   error
		}{
			{
				desc:   "UpdateUser",
				modify: func(barAlice *pb.Membership) { barAlice.User = testresources.Bob.Name },
				want:   ErrUpdateUser,
			},
			{
				desc:   "NotFound",
				modify: func(barAlice *pb.Membership) { barAlice.Name = testresources.Mall_Alice.Name },
				want:   &NotFoundError{Name: testresources.Mall_Alice.Name},
			},
		} {
			t.Run(test.desc, func(t *testing.T) {
				r := SeedRepository(t, []*pb.Membership{testresources.Bar_Alice})
				oldBarAlice := testresources.Bar_Alice
				newBarAlice := Clone(oldBarAlice)
				test.modify(newBarAlice)
				if got := r.UpdateMembership(ctx, newBarAlice); !cmp.Equal(got, test.want, cmpopts.EquateErrors()) {
					t.Errorf("r.UpdateMembership(%v, %v) = %v; want %v", ctx, newBarAlice, got, test.want)
				}
				membership, err := r.LookupMembership(ctx, oldBarAlice.Name)
				if diff := cmp.Diff(membership, oldBarAlice, protocmp.Transform()); diff != "" {
					t.Errorf("r.LookupMembership(%v, %q) membership != testmemberships.Bar_Alice (-got +want)\n%s", ctx, oldBarAlice.Name, diff)
				}
				if err != nil {
					t.Errorf("r.LookupMembership(%v, %q) err = %v; want nil", ctx, oldBarAlice.Name, err)
				}
			})
		}
	})
}

func TestRepository_DeleteMembership(t *testing.T) {
	ctx := context.Background()
	// Test scenario where the delete succeeds.
	t.Run("OK", func(t *testing.T) {
		r := SeedRepository(t, []*pb.Membership{testresources.Bar_Alice})
		// First, delete the membership.
		if err := r.DeleteMembership(ctx, testresources.Bar_Alice.Name); err != nil {
			t.Errorf("r.DeleteMembership(%v, %q) = %v; want nil", ctx, testresources.Bar_Alice.Name, err)
		}
		// Then, verify that looking it up by name fails.
		wantMembership := (*pb.Membership)(nil)
		wantErr := &NotFoundError{Name: testresources.Bar_Alice.Name}
		membership, err := r.LookupMembership(ctx, testresources.Bar_Alice.Name)
		if diff := cmp.Diff(membership, wantMembership, protocmp.Transform()); diff != "" {
			t.Errorf("r.LookupMembership(%v, %q) membership != wantMembership (-got +want)\n%s", ctx, testresources.Bar_Alice.Name, diff)
		}
		if !cmp.Equal(err, wantErr, cmpopts.EquateErrors()) {
			t.Errorf("r.LookupMembership(%v, %q) err = %v; want %v", ctx, testresources.Bar_Alice.Name, err, wantErr)
		}
		// Finally, verify that looking it up by user and parent fails also.
		wantMembership = (*pb.Membership)(nil)
		wantErr = &NotFoundError{
			Parent: testresources.Bar.Name,
			User:   testresources.Alice.Name,
		}
		membership, err = r.LookupMembershipIn(ctx, testresources.Bar.Name, testresources.Alice.Name)
		if diff := cmp.Diff(membership, wantMembership, protocmp.Transform()); diff != "" {
			t.Errorf("r.LookupMembershipIn(%v, %q, %q) membership != wantMembership (-got +want)\n%s", ctx, testresources.Bar.Name, testresources.Alice.Name, diff)
		}
		if !cmp.Equal(err, wantErr, cmpopts.EquateErrors()) {
			t.Errorf("r.LookupMembershipIn(%v, %q, %q) err = %v; want %v", ctx, testresources.Bar.Name, testresources.Alice.Name, err, wantErr)
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
				want: resourcename.ErrInvalidName,
			},
			{
				desc: "NotFound",
				name: testresources.Mall_Alice.Name,
				want: &NotFoundError{Name: testresources.Mall_Alice.Name},
			},
		} {
			t.Run(test.desc, func(t *testing.T) {
				r := SeedRepository(t, []*pb.Membership{testresources.Bar_Alice})
				// First, try and fail to delete the membership.
				if got := r.DeleteMembership(ctx, test.name); !cmp.Equal(got, test.want, cmpopts.EquateErrors()) {
					t.Errorf("r.DeleteMembership(%v, %q) = %v; want %v", ctx, test.name, got, test.want)
				}
				wantMembership := testresources.Bar_Alice
				// Then, verify that a lookup by name succeeds.
				membership, err := r.LookupMembership(ctx, testresources.Bar_Alice.Name)
				if diff := cmp.Diff(membership, wantMembership, protocmp.Transform()); diff != "" {
					t.Errorf("r.LookupMembership(%v, %q) membership != wantMembership (-got +want)\n%s", ctx, testresources.Bar_Alice.Name, diff)
				}
				if err != nil {
					t.Errorf("r.LookupMembership(%v, %q) err = %v; want nil", ctx, testresources.Bar_Alice.Name, err)
				}
				// Finally, verify that a lookup by user and parent succeeds.
				membership, err = r.LookupMembershipIn(ctx, testresources.Bar.Name, testresources.Alice.Name)
				if diff := cmp.Diff(membership, wantMembership, protocmp.Transform()); diff != "" {
					t.Errorf("r.LookupMembershipIn(%v, %q, %q) membership != wantMembership (-got +want)\n%s", ctx, testresources.Bar.Name, testresources.Alice.Name, diff)
				}
				if err != nil {
					t.Errorf("r.LookupMembershipIn(%v, %q, %q) err = %v; want nil", ctx, testresources.Bar.Name, testresources.Alice.Name, err)
				}
			})
		}
	})
}
