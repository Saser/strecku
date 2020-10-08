package stores

import (
	"testing"

	pb "github.com/Saser/strecku/api/v1"
)

func TestValidate(t *testing.T) {
	for _, test := range []struct {
		store *pb.Store
		want  error
	}{
		{store: &pb.Store{Name: "stores/6f2d193c-1460-491d-8157-7dd9535526c6", DisplayName: "Store"}, want: nil},
		{store: &pb.Store{Name: "", DisplayName: "Store"}, want: ErrNameEmpty},
		{store: &pb.Store{Name: "invalidprefix/6f2d193c-1460-491d-8157-7dd9535526c6", DisplayName: "Store"}, want: ErrNameInvalidFormat},
		{store: &pb.Store{Name: "6f2d193c-1460-491d-8157-7dd9535526c6", DisplayName: "Store"}, want: ErrNameInvalidFormat},
		{store: &pb.Store{Name: "stores/6f2d193c-1460-491d-8157-7dd9535526c6", DisplayName: ""}, want: ErrDisplayNameEmpty},
	} {
		if got := Validate(test.store); got != test.want {
			t.Errorf("Validate(%v) = %v; want %v", test.store, got, test.want)
		}
	}
}
