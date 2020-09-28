package memberships

import (
	"testing"

	"github.com/Saser/strecku/resources/stores"
	"github.com/Saser/strecku/resources/users"
	streckuv1 "github.com/Saser/strecku/saser/strecku/v1"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestValidate(t *testing.T) {
	for _, test := range []struct {
		membership *streckuv1.Membership
		want       error
	}{
		{
			membership: &streckuv1.Membership{
				Name:          "memberships/6f2d193c-1460-491d-8157-7dd9535526c6",
				User:          "users/6f2d193c-1460-491d-8157-7dd9535526c6",
				Store:         "stores/6f2d193c-1460-491d-8157-7dd9535526c6",
				Administrator: false,
			},
			want: nil,
		},
		{
			membership: &streckuv1.Membership{
				Name:          "",
				User:          "users/6f2d193c-1460-491d-8157-7dd9535526c6",
				Store:         "stores/6f2d193c-1460-491d-8157-7dd9535526c6",
				Administrator: false,
			},
			want: ErrNameEmpty,
		},
		{
			membership: &streckuv1.Membership{
				Name:          "invalidprefix/6f2d193c-1460-491d-8157-7dd9535526c6",
				User:          "users/6f2d193c-1460-491d-8157-7dd9535526c6",
				Store:         "stores/6f2d193c-1460-491d-8157-7dd9535526c6",
				Administrator: false,
			},
			want: ErrNameInvalidFormat,
		},
		{
			membership: &streckuv1.Membership{
				Name:          "6f2d193c-1460-491d-8157-7dd9535526c6",
				User:          "users/6f2d193c-1460-491d-8157-7dd9535526c6",
				Store:         "stores/6f2d193c-1460-491d-8157-7dd9535526c6",
				Administrator: false,
			},
			want: ErrNameInvalidFormat,
		},
		{
			membership: &streckuv1.Membership{
				Name:          "memberships/not a UUID",
				User:          "users/6f2d193c-1460-491d-8157-7dd9535526c6",
				Store:         "stores/6f2d193c-1460-491d-8157-7dd9535526c6",
				Administrator: false,
			},
			want: ErrNameInvalidFormat,
		},
		{
			membership: &streckuv1.Membership{
				Name:          "memberships/6f2d193c-1460-491d-8157-7dd9535526c6",
				User:          "",
				Store:         "stores/6f2d193c-1460-491d-8157-7dd9535526c6",
				Administrator: false,
			},
			want: users.ErrNameEmpty,
		},
		{
			membership: &streckuv1.Membership{
				Name:          "memberships/6f2d193c-1460-491d-8157-7dd9535526c6",
				User:          "invalidprefix/6f2d193c-1460-491d-8157-7dd9535526c6",
				Store:         "stores/6f2d193c-1460-491d-8157-7dd9535526c6",
				Administrator: false,
			},
			want: users.ErrNameInvalidFormat,
		},
		{
			membership: &streckuv1.Membership{
				Name:          "memberships/6f2d193c-1460-491d-8157-7dd9535526c6",
				User:          "users/not a UUID",
				Store:         "stores/6f2d193c-1460-491d-8157-7dd9535526c6",
				Administrator: false,
			},
			want: users.ErrNameInvalidFormat,
		},
		{
			membership: &streckuv1.Membership{
				Name:          "memberships/6f2d193c-1460-491d-8157-7dd9535526c6",
				User:          "users/6f2d193c-1460-491d-8157-7dd9535526c6",
				Store:         "",
				Administrator: false,
			},
			want: stores.ErrNameEmpty,
		},
		{
			membership: &streckuv1.Membership{
				Name:          "memberships/6f2d193c-1460-491d-8157-7dd9535526c6",
				User:          "users/6f2d193c-1460-491d-8157-7dd9535526c6",
				Store:         "invalidprefix/6f2d193c-1460-491d-8157-7dd9535526c6",
				Administrator: false,
			},
			want: stores.ErrNameInvalidFormat,
		},
		{
			membership: &streckuv1.Membership{
				Name:          "memberships/6f2d193c-1460-491d-8157-7dd9535526c6",
				User:          "users/6f2d193c-1460-491d-8157-7dd9535526c6",
				Store:         "stores/not a UUID",
				Administrator: false,
			},
			want: stores.ErrNameInvalidFormat,
		},
	} {
		if got := Validate(test.membership); got != test.want {
			t.Errorf("Validate(%v) = %v; want %v", test.membership, got, test.want)
		}
	}
}

func TestValidateName(t *testing.T) {
	for _, test := range []struct {
		name string
		want error
	}{
		{name: "memberships/6f2d193c-1460-491d-8157-7dd9535526c6", want: nil},
		{name: "", want: ErrNameEmpty},
		{name: "invalidprefix/6f2d193c-1460-491d-8157-7dd9535526c6", want: ErrNameInvalidFormat},
		{name: "memberships/not a UUID", want: ErrNameInvalidFormat},
		{name: "6f2d193c-1460-491d-8157-7dd9535526c6", want: ErrNameInvalidFormat},
	} {
		if got := ValidateName(test.name); !cmp.Equal(got, test.want, cmpopts.EquateErrors()) {
			t.Errorf("ValidateName(%q) = %v; want %v", test.name, got, test.want)
		}
	}
}
