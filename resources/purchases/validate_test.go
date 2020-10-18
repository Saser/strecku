package purchases

import (
	"testing"

	pb "github.com/Saser/strecku/api/v1"
	"github.com/Saser/strecku/resources/stores"
	"github.com/Saser/strecku/resources/testresources"
	"github.com/Saser/strecku/resources/users"
)

func TestValidate(t *testing.T) {
	// In the following test cases, the valid purchase will be a clone of testresources.Alice_Beer1.
	for _, test := range []struct {
		modify func(valid *pb.Purchase)
		want   error
	}{
		{
			modify: func(valid *pb.Purchase) { valid.User = "" },
			want:   users.ErrNameInvalidFormat,
		},
		{
			modify: func(valid *pb.Purchase) { valid.Store = "" },
			want:   stores.ErrNameInvalidFormat,
		},
		{
			modify: func(valid *pb.Purchase) { valid.Lines = nil },
			want:   ErrLinesEmpty,
		},
		{
			modify: func(valid *pb.Purchase) { valid.Lines = []*pb.Purchase_Line{} },
			want:   ErrLinesEmpty,
		},
		{
			modify: func(valid *pb.Purchase) { valid.Lines = make([]*pb.Purchase_Line, 0, 10) },
			want:   ErrLinesEmpty,
		},
		{
			modify: func(valid *pb.Purchase) { valid.Lines[0].Description = "" },
			want:   ErrLineDescriptionEmpty,
		},
		{
			modify: func(valid *pb.Purchase) { valid.Lines[0].Quantity = 0 },
			want:   ErrLineQuantityNonPositive,
		},
		{
			modify: func(valid *pb.Purchase) { valid.Lines[0].Quantity = -1 },
			want:   ErrLineQuantityNonPositive,
		},
		{
			modify: func(valid *pb.Purchase) { valid.Lines[0].PriceCents = 1000 },
			want:   ErrLinePricePositive,
		},
		{
			modify: func(valid *pb.Purchase) { valid.Lines[0].Product = testresources.Pills.Name },
			want:   ErrLineProductWrongStore,
		},
	} {
		purchase := Clone(testresources.Alice_Beer1)
		test.modify(purchase)
		if got := Validate(purchase); got != test.want {
			t.Errorf("Validate(%v) = %v; want %v", purchase, got, test.want)
		}
	}
}
