package purchases

import (
	"errors"

	pb "github.com/Saser/strecku/api/v1"
	"github.com/Saser/strecku/resources/stores/products"
	"github.com/Saser/strecku/resources/users"
	"google.golang.org/protobuf/proto"
)

var (
	ErrLinesEmpty              = errors.New("lines are empty")
	ErrLineDescriptionEmpty    = errors.New("line description is empty")
	ErrLineQuantityNonPositive = errors.New("line quantity is non-positive")
	ErrLinePricePositive       = errors.New("line price is positive")
	ErrLineProductWrongStore   = errors.New("line product belongs to another store")
)

func Clone(purchase *pb.Purchase) *pb.Purchase {
	return proto.Clone(purchase).(*pb.Purchase)
}

func Validate(purchase *pb.Purchase) error {
	if err := ValidateName(purchase.Name); err != nil {
		return err
	}
	if err := users.ValidateName(purchase.User); err != nil {
		return err
	}
	if len(purchase.Lines) == 0 {
		return ErrLinesEmpty
	}
	for _, line := range purchase.Lines {
		if err := ValidateLine(purchase, line); err != nil {
			return err
		}
	}
	return nil
}

func ValidateLine(purchase *pb.Purchase, line *pb.Purchase_Line) error {
	if line.Description == "" {
		return ErrLineDescriptionEmpty
	}
	if line.Quantity <= 0 {
		return ErrLineQuantityNonPositive
	}
	if line.PriceCents > 0 {
		return ErrLinePricePositive
	}
	product := line.Product
	if product != "" {
		if err := products.ValidateName(product); err != nil {
			return err
		}
		store, _ := products.Parent(product)
		parent, _ := Parent(purchase.Name)
		if parent != store {
			return ErrLineProductWrongStore
		}
	}
	return nil
}
