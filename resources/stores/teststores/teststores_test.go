package teststores

import (
	"testing"

	pb "github.com/Saser/strecku/api/v1"
	"github.com/Saser/strecku/resources/stores"
)

func TestValid(t *testing.T) {
	for _, store := range []*pb.Store{
		Pharmacy,
		Bar,
	} {
		if err := stores.Validate(store); err != nil {
			t.Errorf("stores.Validate(%v) = %v; want nil", store, err)
		}
	}
}
