package memberships

import (
	"testing"

	pb "github.com/Saser/strecku/api/v1"
	"github.com/Saser/strecku/resources/stores"
	"github.com/Saser/strecku/resources/users"
)

func TestValidate(t *testing.T) {
	for _, test := range []struct {
		membership *pb.Membership
		want       error
	}{
		{
			membership: &pb.Membership{
				Name:          "memberships/6f2d193c-1460-491d-8157-7dd9535526c6",
				User:          "users/6f2d193c-1460-491d-8157-7dd9535526c6",
				Store:         "stores/6f2d193c-1460-491d-8157-7dd9535526c6",
				Administrator: false,
			},
			want: nil,
		},
		{
			membership: &pb.Membership{
				Name:          "",
				User:          "users/6f2d193c-1460-491d-8157-7dd9535526c6",
				Store:         "stores/6f2d193c-1460-491d-8157-7dd9535526c6",
				Administrator: false,
			},
			want: ErrNameEmpty,
		},
		{
			membership: &pb.Membership{
				Name:          "invalidprefix/6f2d193c-1460-491d-8157-7dd9535526c6",
				User:          "users/6f2d193c-1460-491d-8157-7dd9535526c6",
				Store:         "stores/6f2d193c-1460-491d-8157-7dd9535526c6",
				Administrator: false,
			},
			want: ErrNameInvalidFormat,
		},
		{
			membership: &pb.Membership{
				Name:          "6f2d193c-1460-491d-8157-7dd9535526c6",
				User:          "users/6f2d193c-1460-491d-8157-7dd9535526c6",
				Store:         "stores/6f2d193c-1460-491d-8157-7dd9535526c6",
				Administrator: false,
			},
			want: ErrNameInvalidFormat,
		},
		{
			membership: &pb.Membership{
				Name:          "memberships/not a UUID",
				User:          "users/6f2d193c-1460-491d-8157-7dd9535526c6",
				Store:         "stores/6f2d193c-1460-491d-8157-7dd9535526c6",
				Administrator: false,
			},
			want: ErrNameInvalidFormat,
		},
		{
			membership: &pb.Membership{
				Name:          "memberships/6f2d193c-1460-491d-8157-7dd9535526c6",
				User:          "",
				Store:         "stores/6f2d193c-1460-491d-8157-7dd9535526c6",
				Administrator: false,
			},
			want: users.ErrNameEmpty,
		},
		{
			membership: &pb.Membership{
				Name:          "memberships/6f2d193c-1460-491d-8157-7dd9535526c6",
				User:          "invalidprefix/6f2d193c-1460-491d-8157-7dd9535526c6",
				Store:         "stores/6f2d193c-1460-491d-8157-7dd9535526c6",
				Administrator: false,
			},
			want: users.ErrNameInvalidFormat,
		},
		{
			membership: &pb.Membership{
				Name:          "memberships/6f2d193c-1460-491d-8157-7dd9535526c6",
				User:          "users/not a UUID",
				Store:         "stores/6f2d193c-1460-491d-8157-7dd9535526c6",
				Administrator: false,
			},
			want: users.ErrNameInvalidFormat,
		},
		{
			membership: &pb.Membership{
				Name:          "memberships/6f2d193c-1460-491d-8157-7dd9535526c6",
				User:          "users/6f2d193c-1460-491d-8157-7dd9535526c6",
				Store:         "",
				Administrator: false,
			},
			want: stores.ErrNameEmpty,
		},
		{
			membership: &pb.Membership{
				Name:          "memberships/6f2d193c-1460-491d-8157-7dd9535526c6",
				User:          "users/6f2d193c-1460-491d-8157-7dd9535526c6",
				Store:         "invalidprefix/6f2d193c-1460-491d-8157-7dd9535526c6",
				Administrator: false,
			},
			want: stores.ErrNameInvalidFormat,
		},
		{
			membership: &pb.Membership{
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
