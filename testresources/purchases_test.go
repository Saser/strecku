package testresources

import (
	"testing"

	pb "github.com/Saser/strecku/api/v1"
	"github.com/Saser/strecku/resources/stores/purchases"
)

func TestPurchasesValid(t *testing.T) {
	for _, purchase := range []*pb.Purchase{
		Bar_Alice_Beer1,
		Bar_Alice_Cocktail1,
		Bar_Alice_Beer2_Cocktail2,
		Mall_Alice_Jeans1,
	} {
		if err := purchases.Validate(purchase); err != nil {
			t.Errorf("purchases.Validate(%v) = %v; want nil", purchase, err)
		}
	}
}
