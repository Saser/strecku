package memberships

import (
	"context"
	"fmt"
	"testing"

	pb "github.com/Saser/strecku/api/v1"
	"github.com/Saser/strecku/resources/memberships/testmemberships"
	"github.com/Saser/strecku/resources/stores"
	"github.com/Saser/strecku/resources/stores/teststores"
	"github.com/Saser/strecku/resources/users"
	"github.com/Saser/strecku/resources/users/testusers"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"google.golang.org/protobuf/testing/protocmp"
)

func membershipLess(m1, m2 *pb.Membership) bool {
	return m1.Name < m2.Name
}

func TestMembershipNotFoundError_Error(t *testing.T) {
	for _, test := range []struct {
		err  *MembershipNotFoundError
		want string
	}{
		{
			err: &MembershipNotFoundError{
				Name: testmemberships.Alice_Bar.Name,
			},
			want: fmt.Sprintf("membership not found: %q", testmemberships.Alice_Bar.Name),
		},
		{
			err: &MembershipNotFoundError{
				User:  testusers.Alice.Name,
				Store: teststores.Bar.Name,
			},
			want: fmt.Sprintf("membership not found: between %q and %q", testusers.Alice.Name, teststores.Bar.Name),
		},
	} {
		if got := test.err.Error(); got != test.want {
			t.Errorf("test.err.Error() = %q; want %q", got, test.want)
		}
	}
}

func TestMembershipNotFoundError_Is(t *testing.T) {
	for _, test := range []struct {
		err    *MembershipNotFoundError
		target error
		want   bool
	}{
		{
			err: &MembershipNotFoundError{
				Name: testmemberships.Alice_Bar.Name,
			},
			target: &MembershipNotFoundError{
				Name: testmemberships.Alice_Bar.Name,
			},
			want: true,
		},
		{
			err: &MembershipNotFoundError{
				Name: testmemberships.Alice_Bar.Name,
			},
			target: &MembershipNotFoundError{
				User:  testusers.Alice.Name,
				Store: teststores.Bar.Name,
			},
			want: false,
		},
		{
			err: &MembershipNotFoundError{
				User:  testusers.Alice.Name,
				Store: teststores.Bar.Name,
			},
			target: &MembershipNotFoundError{
				Name: testmemberships.Alice_Bar.Name,
			},
			want: false,
		},
		{
			err: &MembershipNotFoundError{
				User:  testusers.Alice.Name,
				Store: teststores.Bar.Name,
			},
			target: &MembershipNotFoundError{
				User:  testusers.Alice.Name,
				Store: teststores.Bar.Name,
			},
			want: true,
		},
	} {
		if got := test.err.Is(test.target); got != test.want {
			t.Errorf("test.err.Is(%v) = %v; want %v", test.target, got, test.want)
		}
	}
}

func TestMembershipExistsError_Error(t *testing.T) {
	for _, test := range []struct {
		err  *MembershipExistsError
		want string
	}{
		{
			err: &MembershipExistsError{
				Name: testmemberships.Alice_Bar.Name,
			},
			want: fmt.Sprintf("membership exists: %q", testmemberships.Alice_Bar.Name),
		},
		{
			err: &MembershipExistsError{
				User:  testusers.Alice.Name,
				Store: teststores.Bar.Name,
			},
			want: fmt.Sprintf("membership exists: between %q and %q", testusers.Alice.Name, teststores.Bar.Name),
		},
	} {
		if got := test.err.Error(); got != test.want {
			t.Errorf("test.err.Error() = %q; want %q", got, test.want)
		}
	}
}

func TestRepository_LookupMembership(t *testing.T) {
	ctx := context.Background()
	r := SeedRepository(t, []*pb.Membership{testmemberships.Alice_Bar})
	for _, test := range []struct {
		desc           string
		name           string
		wantMembership *pb.Membership
		wantErr        error
	}{
		{
			desc:           "OK",
			name:           testmemberships.Alice_Bar.Name,
			wantMembership: testmemberships.Alice_Bar,
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
			name:           testmemberships.Alice_Mall.Name,
			wantMembership: nil,
			wantErr:        &MembershipNotFoundError{Name: testmemberships.Alice_Mall.Name},
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
	r := SeedRepository(t, []*pb.Membership{testmemberships.Alice_Bar})
	for _, test := range []struct {
		desc           string
		user           string
		store          string
		wantMembership *pb.Membership
		wantErr        error
	}{
		{
			desc:           "OK",
			user:           testusers.Alice.Name,
			store:          teststores.Bar.Name,
			wantMembership: testmemberships.Alice_Bar,
			wantErr:        nil,
		},
		{
			desc:           "EmptyUser",
			user:           "",
			store:          teststores.Bar.Name,
			wantMembership: nil,
			wantErr:        users.ErrNameEmpty,
		},
		{
			desc:           "EmptyStore",
			user:           testusers.Alice.Name,
			store:          "",
			wantMembership: nil,
			wantErr:        stores.ErrNameEmpty,
		},
		{
			desc:           "WrongUser",
			user:           testusers.Bob.Name,
			store:          teststores.Bar.Name,
			wantMembership: nil,
			wantErr:        &MembershipNotFoundError{User: testusers.Bob.Name, Store: teststores.Bar.Name},
		},
		{
			desc:           "WrongStore",
			user:           testusers.Alice.Name,
			store:          teststores.Mall.Name,
			wantMembership: nil,
			wantErr:        &MembershipNotFoundError{User: testusers.Alice.Name, Store: teststores.Mall.Name},
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
		testmemberships.Alice_Bar,
		testmemberships.Alice_Mall,
		testmemberships.Bob_Bar,
		testmemberships.Bob_Mall,
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
		testmemberships.Alice_Bar,
		testmemberships.Alice_Mall,
		testmemberships.Bob_Bar,
		testmemberships.Bob_Mall,
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
			predicate: func(membership *pb.Membership) bool { return membership.Name == testmemberships.Alice_Bar.Name },
			want: []*pb.Membership{
				testmemberships.Alice_Bar,
			},
		},
		{
			desc:      "MultipleMatching",
			predicate: func(membership *pb.Membership) bool { return membership.User == testusers.Alice.Name },
			want: []*pb.Membership{
				testmemberships.Alice_Bar,
				testmemberships.Alice_Mall,
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
			membership: testmemberships.Alice_Mall,
			want:       nil,
		},
		{
			desc:       "OK_SameStore",
			membership: testmemberships.Bob_Bar,
			want:       nil,
		},
		{
			desc: "DuplicateName",
			membership: &pb.Membership{
				Name:          testmemberships.Alice_Bar.Name,
				User:          testusers.Bob.Name,   // chosen arbitrarily
				Store:         teststores.Mall.Name, // chosen arbitrarily
				Administrator: false,
			},
			want: &MembershipExistsError{Name: testmemberships.Alice_Bar.Name},
		},
		{
			desc: "DuplicateUserAndStore",
			membership: &pb.Membership{
				Name:          testmemberships.Bob_Mall.Name, // chosen arbitrarily
				User:          testusers.Alice.Name,
				Store:         teststores.Bar.Name,
				Administrator: false,
			},
			want: &MembershipExistsError{
				User:  testusers.Alice.Name,
				Store: teststores.Bar.Name,
			},
		},
	} {
		t.Run(test.desc, func(t *testing.T) {
			r := SeedRepository(t, []*pb.Membership{testmemberships.Alice_Bar})
			if got := r.CreateMembership(ctx, test.membership); !cmp.Equal(got, test.want, cmpopts.EquateErrors()) {
				t.Errorf("r.CreateMembership(%v, %v) = %v; want %v", ctx, test.membership, got, test.want)
			}
		})
	}
}
