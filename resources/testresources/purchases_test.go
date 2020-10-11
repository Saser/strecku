package testresources

import (
	"testing"

	pb "github.com/Saser/strecku/api/v1"
	"github.com/Saser/strecku/resources/purchases"
)

func TestPurchasesValid(t *testing.T) {
	for _, purchase := range []*pb.Purchase{
		Alice_Beer1,
	} {
		if err := purchases.Validate(purchase); err != nil {
			t.Errorf("purchases.Validate(%v) = %v; want nil", purchase, err)
		}
	}
}
