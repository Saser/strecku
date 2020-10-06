package memberships

import (
	"context"
	"fmt"
	"testing"

	pb "github.com/Saser/strecku/api/v1"
	"github.com/Saser/strecku/resources/memberships/testmemberships"
	"github.com/Saser/strecku/resources/stores/teststores"
	"github.com/Saser/strecku/resources/users/testusers"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"google.golang.org/protobuf/testing/protocmp"
)

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
