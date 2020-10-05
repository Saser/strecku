package memberships

import (
	"fmt"
	"testing"

	"github.com/Saser/strecku/resources/memberships/testmemberships"
	"github.com/Saser/strecku/resources/stores/teststores"
	"github.com/Saser/strecku/resources/users/testusers"
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
