package memberships

import (
	"fmt"
	"testing"
)

// Valid resource names of users, stores, and membership relations between them.
const (
	Alice = "users/6f2d193c-1460-491d-8157-7dd9535526c6"

	GroceryStore = "stores/d8bbf79e-8c59-4fae-aef9-634fcac00e07"

	Alice_GroceryStore = "memberships/9cd3ec05-e7af-418c-bd50-80a7c39a18cc"
)

func TestMembershipNotFoundError_Error(t *testing.T) {
	for _, test := range []struct {
		err  *MembershipNotFoundError
		want string
	}{
		{
			err: &MembershipNotFoundError{
				Name: Alice_GroceryStore,
			},
			want: fmt.Sprintf("membership not found: %q", Alice_GroceryStore),
		},
		{
			err: &MembershipNotFoundError{
				User:  Alice,
				Store: GroceryStore,
			},
			want: fmt.Sprintf("membership not found: between %q and %q", Alice, GroceryStore),
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
				Name: Alice_GroceryStore,
			},
			target: &MembershipNotFoundError{
				Name: Alice_GroceryStore,
			},
			want: true,
		},
		{
			err: &MembershipNotFoundError{
				Name: Alice_GroceryStore,
			},
			target: &MembershipNotFoundError{
				User:  Alice,
				Store: GroceryStore,
			},
			want: false,
		},
		{
			err: &MembershipNotFoundError{
				User:  Alice,
				Store: GroceryStore,
			},
			target: &MembershipNotFoundError{
				Name: Alice_GroceryStore,
			},
			want: false,
		},
		{
			err: &MembershipNotFoundError{
				User:  Alice,
				Store: GroceryStore,
			},
			target: &MembershipNotFoundError{
				User:  Alice,
				Store: GroceryStore,
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
				Name: Alice_GroceryStore,
			},
			want: fmt.Sprintf("membership exists: %q", Alice_GroceryStore),
		},
		{
			err: &MembershipExistsError{
				User:  Alice,
				Store: GroceryStore,
			},
			want: fmt.Sprintf("membership exists: between %q and %q", Alice, GroceryStore),
		},
	} {
		if got := test.err.Error(); got != test.want {
			t.Errorf("test.err.Error() = %q; want %q", got, test.want)
		}
	}
}
