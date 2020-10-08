package stores

import (
	"testing"

	pb "github.com/Saser/strecku/api/v1"
	"github.com/Saser/strecku/resources/stores/teststores"
)

func TestValidate(t *testing.T) {
	for _, test := range []struct {
		store *pb.Store
		want  error
	}{
		{
			store: &pb.Store{
				Name:        teststores.Bar.Name,
				DisplayName: "",
			},
			want: ErrDisplayNameEmpty,
		},
	} {
		if got := Validate(test.store); got != test.want {
			t.Errorf("Validate(%v) = %v; want %v", test.store, got, test.want)
		}
	}
}
