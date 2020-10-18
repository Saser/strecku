package payments

import (
	"errors"

	pb "github.com/Saser/strecku/api/v1"
	"github.com/Saser/strecku/resources/users"
	"google.golang.org/protobuf/proto"
)

var ErrAmountNegative = errors.New("amount is negative")

func Clone(payment *pb.Payment) *pb.Payment {
	return proto.Clone(payment).(*pb.Payment)
}

func Validate(payment *pb.Payment) error {
	if err := ValidateName(payment.Name); err != nil {
		return err
	}
	if err := users.ValidateName(payment.User); err != nil {
		return err
	}
	if payment.AmountCents < 0 {
		return ErrAmountNegative
	}
	return nil
}
