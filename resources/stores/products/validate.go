package products

import (
	"errors"

	pb "github.com/Saser/strecku/api/v1"
)

var (
	ErrDisplayNameEmpty                 = errors.New("display name is empty")
	ErrFullPricePositive                = errors.New("full price is positive")
	ErrDiscountPricePositive            = errors.New("discount price is positive")
	ErrDiscountPriceHigherThanFullPrice = errors.New("discount price is higher than the full price")
)

func Validate(product *pb.Product) error {
	if err := ValidateName(product.Name); err != nil {
		return err
	}
	if product.DisplayName == "" {
		return ErrDisplayNameEmpty
	}
	if product.FullPriceCents > 0 {
		return ErrFullPricePositive
	}
	if product.DiscountPriceCents > 0 {
		return ErrDiscountPricePositive
	}
	if product.DiscountPriceCents < product.FullPriceCents {
		return ErrDiscountPriceHigherThanFullPrice
	}
	return nil
}
