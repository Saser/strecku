package testresources

import (
	"testing"

	pb "github.com/Saser/strecku/api/v1"
	"github.com/Saser/strecku/resources/products"
)

func TestProductsValid(t *testing.T) {
	for _, product := range []*pb.Product{
		Beer,
		Cocktail,
		Pills,
		Lotion,
	} {
		if err := products.Validate(product); err != nil {
			t.Errorf("products.Validate(%v) = %v; want nil", product, err)
		}
	}
}
