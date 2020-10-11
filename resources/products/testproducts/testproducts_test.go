package testproducts

import (
	"testing"

	pb "github.com/Saser/strecku/api/v1"
	"github.com/Saser/strecku/resources/products"
)

func TestValid(t *testing.T) {
	for _, product := range []*pb.Product{
		Bar_Beer,
		Bar_Cocktail,
	} {
		if err := products.Validate(product); err != nil {
			t.Errorf("products.Validate(%v) = %v; want nil", product, err)
		}
	}
}
