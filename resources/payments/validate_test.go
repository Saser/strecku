package payments

import (
	"testing"

	pb "github.com/Saser/strecku/api/v1"
	"github.com/Saser/strecku/resources/testresources"
	"github.com/Saser/strecku/resources/users"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestValidate(t *testing.T) {
	for _, test := range []struct {
		payment *pb.Payment
		want    error
	}{
		{
			payment: func() *pb.Payment {
				payment := Clone(testresources.Bar_Alice_Payment)
				payment.Name = ""
				return payment
			}(),
			want: ErrNameInvalidFormat,
		},
		{
			payment: func() *pb.Payment {
				payment := Clone(testresources.Bar_Alice_Payment)
				payment.Name = testresources.Bar.Name // not a name of a payment
				return payment
			}(),
			want: ErrNameInvalidFormat,
		},
		{
			payment: func() *pb.Payment {
				payment := Clone(testresources.Bar_Alice_Payment)
				payment.User = ""
				return payment
			}(),
			want: users.ErrNameInvalidFormat,
		},
		{
			payment: func() *pb.Payment {
				payment := Clone(testresources.Bar_Alice_Payment)
				payment.AmountCents = -10000
				return payment
			}(),
			want: ErrAmountNegative,
		},
	} {
		if got := Validate(test.payment); !cmp.Equal(got, test.want, cmpopts.EquateErrors()) {
			t.Errorf("Validate(%v) = %v; want %v", test.payment, got, test.want)
		}
	}
}
