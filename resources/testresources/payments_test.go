package testresources

import (
	"testing"

	pb "github.com/Saser/strecku/api/v1"
	"github.com/Saser/strecku/resources/stores/payments"
)

func TestPaymentsValid(t *testing.T) {
	for _, payment := range []*pb.Payment{
		Bar_Alice_Payment,
		Bar_Bob_Payment,
		Bar_Carol_Payment,
	} {
		if err := payments.Validate(payment); err != nil {
			t.Errorf("payments.Validate(%v) = %v; want nil", payment, err)
		}
	}
}
