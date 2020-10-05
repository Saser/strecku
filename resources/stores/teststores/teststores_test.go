package teststores

import (
	"testing"

	"github.com/Saser/strecku/resources/stores"
	streckuv1 "github.com/Saser/strecku/saser/strecku/v1"
)

func TestValid(t *testing.T) {
	for _, store := range []*streckuv1.Store{
		Pharmacy,
		Bar,
	} {
		if err := stores.Validate(store); err != nil {
			t.Errorf("stores.Validate(%v) = %v; want nil", store, err)
		}
	}
}
