package payments

import (
	"context"
	"fmt"
	"testing"

	pb "github.com/Saser/strecku/api/v1"
	"github.com/Saser/strecku/resourcename"
	"github.com/Saser/strecku/testresources"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"google.golang.org/protobuf/testing/protocmp"
)

func paymentLess(p1, p2 *pb.Payment) bool {
	return p1.Name < p2.Name
}

func TestNotFoundError_Error(t *testing.T) {
	err := &NotFoundError{Name: testresources.Bar_Alice_Payment.Name}
	want := fmt.Sprintf("payment not found: %q", testresources.Bar_Alice_Payment.Name)
	if got := err.Error(); !cmp.Equal(got, want) {
		t.Errorf("err.Error() = %q; want %q", got, want)
	}
}

func TestNotFoundError_Is(t *testing.T) {
	for _, test := range []struct {
		err    *NotFoundError
		target error
		want   bool
	}{
		{
			err:    &NotFoundError{Name: testresources.Bar_Alice_Payment.Name},
			target: &NotFoundError{Name: testresources.Bar_Alice_Payment.Name},
			want:   true,
		},
		{
			err:    &NotFoundError{Name: testresources.Bar_Alice_Payment.Name},
			target: &NotFoundError{Name: testresources.Bar_Bob_Payment.Name},
			want:   false,
		},
		{
			err:    &NotFoundError{Name: testresources.Bar_Alice_Payment.Name},
			target: fmt.Errorf("payment not found: %q", testresources.Bar_Alice_Payment.Name),
			want:   false,
		},
	} {
		if got := test.err.Is(test.target); got != test.want {
			t.Errorf("test.err.Is(%v) = %v; want %v", test.target, got, test.want)
		}
	}
}

func TestExistsError_Error(t *testing.T) {
	err := &ExistsError{Name: testresources.Bar_Alice_Payment.Name}
	want := fmt.Sprintf("payment exists: %q", testresources.Bar_Alice_Payment.Name)
	if got := err.Error(); !cmp.Equal(got, want) {
		t.Errorf("err.Error() = %q; want %q", got, want)
	}
}

func TestExistsError_Is(t *testing.T) {
	for _, test := range []struct {
		err    *ExistsError
		target error
		want   bool
	}{
		{
			err:    &ExistsError{Name: testresources.Bar_Alice_Payment.Name},
			target: &ExistsError{Name: testresources.Bar_Alice_Payment.Name},
			want:   true,
		},
		{
			err:    &ExistsError{Name: testresources.Bar_Alice_Payment.Name},
			target: &ExistsError{Name: testresources.Bar_Bob_Payment.Name},
			want:   false,
		},
		{
			err:    &ExistsError{Name: testresources.Bar_Alice_Payment.Name},
			target: fmt.Errorf("payment exists: %q", testresources.Bar_Alice_Payment.Name),
			want:   false,
		},
	} {
		if got := test.err.Is(test.target); got != test.want {
			t.Errorf("test.err.Is(%v) = %v; want %v", test.target, got, test.want)
		}
	}
}

func TestNewRepository(t *testing.T) {
	NewRepository()
}

func TestRepository_LookupPayment(t *testing.T) {
	ctx := context.Background()
	r := SeedRepository(t, []*pb.Payment{testresources.Bar_Alice_Payment})
	for _, test := range []struct {
		desc        string
		name        string
		wantPayment *pb.Payment
		wantErr     error
	}{
		{
			desc:        "OK",
			name:        testresources.Bar_Alice_Payment.Name,
			wantPayment: testresources.Bar_Alice_Payment,
			wantErr:     nil,
		},
		{
			desc:        "EmptyName",
			name:        "",
			wantPayment: nil,
			wantErr:     resourcename.ErrInvalidName,
		},
		{
			desc:        "NotFound",
			name:        testresources.Bar_Bob_Payment.Name,
			wantPayment: nil,
			wantErr:     &NotFoundError{Name: testresources.Bar_Bob_Payment.Name},
		},
	} {
		t.Run(test.desc, func(t *testing.T) {
			payment, err := r.LookupPayment(ctx, test.name)
			if diff := cmp.Diff(payment, test.wantPayment, protocmp.Transform()); diff != "" {
				t.Errorf("r.LookupPayment(%v, %q) payment != test.wantPayment (-got +want)\n%s", ctx, test.name, diff)
			}
			if !cmp.Equal(err, test.wantErr, cmpopts.EquateErrors()) {
				t.Errorf("r.LookupPayment(%v, %q) err = %v; want %v", ctx, test.name, err, test.wantErr)
			}
		})
	}
}

func TestRepository_ListPayments(t *testing.T) {
	ctx := context.Background()
	allPayments := []*pb.Payment{
		testresources.Bar_Alice_Payment,
		testresources.Bar_Bob_Payment,
		testresources.Bar_Carol_Payment,
	}
	r := SeedRepository(t, allPayments)
	payments, err := r.ListPayments(ctx)
	if diff := cmp.Diff(
		payments, allPayments, protocmp.Transform(),
		cmpopts.EquateEmpty(),
		cmpopts.SortSlices(paymentLess),
	); diff != "" {
		t.Errorf("r.ListPayments(%v) payments != allPayments (-got +want)\n%s", ctx, diff)
	}
	if err != nil {
		t.Errorf("r.ListPayments(%v) err = %v; want nil", ctx, err)
	}
}

func TestRepository_FilterPayments(t *testing.T) {
	ctx := context.Background()
	r := SeedRepository(t, []*pb.Payment{
		testresources.Bar_Alice_Payment,
		testresources.Bar_Bob_Payment,
		testresources.Bar_Carol_Payment,
	})
	for _, test := range []struct {
		desc      string
		predicate func(*pb.Payment) bool
		want      []*pb.Payment
	}{
		{
			desc:      "NoneMatching",
			predicate: func(*pb.Payment) bool { return false },
			want:      nil,
		},
		{
			desc:      "OneMatching",
			predicate: func(payment *pb.Payment) bool { return payment.Name == testresources.Bar_Alice_Payment.Name },
			want: []*pb.Payment{
				testresources.Bar_Alice_Payment,
			},
		},
		{
			desc: "MultipleMatching",
			predicate: func(payment *pb.Payment) bool {
				return payment.Name == testresources.Bar_Alice_Payment.Name || payment.Name == testresources.Bar_Bob_Payment.Name
			},
			want: []*pb.Payment{
				testresources.Bar_Alice_Payment,
				testresources.Bar_Bob_Payment,
			},
		},
	} {
		t.Run(test.desc, func(t *testing.T) {
			filtered, err := r.FilterPayments(ctx, test.predicate)
			if diff := cmp.Diff(
				filtered, test.want, protocmp.Transform(),
				cmpopts.EquateEmpty(),
				cmpopts.SortSlices(paymentLess),
			); diff != "" {
				t.Errorf("r.FilterPayments(%v, test.predicate) filtered != test.want (-got +want)\n%s", ctx, diff)
			}
			if err != nil {
				t.Errorf("r.FilterPayments(%v, test.predicate) err = %v; want nil", ctx, err)
			}
		})
	}
}

func TestRepository_CreatePayment(t *testing.T) {
	ctx := context.Background()
	for _, test := range []struct {
		desc    string
		payment *pb.Payment
		want    error
	}{
		{
			desc:    "OK",
			payment: testresources.Bar_Bob_Payment,
			want:    nil,
		},
		{
			desc: "DuplicateName",
			payment: func() *pb.Payment {
				payment := Clone(testresources.Bar_Bob_Payment)
				payment.Name = testresources.Bar_Alice_Payment.Name
				return payment
			}(),
			want: &ExistsError{Name: testresources.Bar_Alice_Payment.Name},
		},
		{
			desc: "InvalidPayment",
			payment: func() *pb.Payment {
				// Create a payment with negative amount.
				// This type of invalidity was chosen arbitrarily.
				payment := Clone(testresources.Bar_Bob_Payment)
				payment.AmountCents = -10000
				return payment
			}(),
			want: ErrAmountNegative,
		},
	} {
		t.Run(test.desc, func(t *testing.T) {
			r := SeedRepository(t, []*pb.Payment{testresources.Bar_Alice_Payment})
			if got := r.CreatePayment(ctx, test.payment); !cmp.Equal(got, test.want, cmpopts.EquateErrors()) {
				t.Errorf("r.CreatePayment(%v, %v) = %v; want %v", ctx, test.payment, got, test.want)
			}
		})
	}
}

func TestRepository_UpdatePayment(t *testing.T) {
	ctx := context.Background()
	// Test scenario(s) where the update is successful.
	t.Run("OK", func(t *testing.T) {
		r := SeedRepository(t, []*pb.Payment{testresources.Bar_Alice_Payment})
		oldPayment := Clone(testresources.Bar_Alice_Payment)
		newPayment := Clone(oldPayment)
		newPayment.AmountCents = 50000
		newPayment.Description = "Alice's new payment"
		if err := r.UpdatePayment(ctx, newPayment); err != nil {
			t.Errorf("r.UpdatePayment(%v, %v) = %v; want nil", ctx, newPayment, err)
		}
		payment, err := r.LookupPayment(ctx, newPayment.Name)
		if diff := cmp.Diff(payment, newPayment, protocmp.Transform()); diff != "" {
			t.Errorf("r.LookupPayment(%v, %q) payment != newPayment (-got +want)\n%s", ctx, newPayment.Name, diff)
		}
		if err != nil {
			t.Errorf("r.LookupPayment(%v, %q) err = %v; want nil", ctx, newPayment.Name, err)
		}
	})
	// Test scenario(s) where the update failed.
	t.Run("Errors", func(t *testing.T) {
		r := SeedRepository(t, []*pb.Payment{testresources.Bar_Alice_Payment})
		for _, test := range []struct {
			desc   string
			modify func(payment *pb.Payment)
			want   error
		}{
			{
				desc:   "UpdateUser",
				modify: func(payment *pb.Payment) { payment.User = testresources.Bob.Name },
				want:   ErrUpdateUser,
			},
			{
				desc:   "NotFound",
				modify: func(payment *pb.Payment) { payment.Name = testresources.Bar_Bob_Payment.Name },
				want:   &NotFoundError{Name: testresources.Bar_Bob_Payment.Name},
			},
		} {
			t.Run(test.desc, func(t *testing.T) {
				oldPayment := Clone(testresources.Bar_Alice_Payment)
				newPayment := Clone(oldPayment)
				test.modify(newPayment)
				if got := r.UpdatePayment(ctx, newPayment); !cmp.Equal(got, test.want, cmpopts.EquateErrors()) {
					t.Errorf("r.UpdatePayment(%v, %v) = %v; want %v", ctx, newPayment, got, test.want)
				}
				payment, err := r.LookupPayment(ctx, oldPayment.Name)
				if diff := cmp.Diff(payment, oldPayment, protocmp.Transform()); diff != "" {
					t.Errorf("r.LookupPayment(%v, %q) payment != oldPayment (-got +want)\n%s", ctx, oldPayment.Name, diff)
				}
				if err != nil {
					t.Errorf("r.LookupPayment(%v, %q) err = %v; want nil", ctx, oldPayment.Name, err)
				}
			})
		}
	})
}

func TestRepository_DeletePayment(t *testing.T) {
	ctx := context.Background()
	t.Run("OK", func(t *testing.T) {
		r := SeedRepository(t, []*pb.Payment{
			testresources.Bar_Alice_Payment,
			testresources.Bar_Bob_Payment,
		})
		if err := r.DeletePayment(ctx, testresources.Bar_Alice_Payment.Name); err != nil {
			t.Errorf("r.DeletePayment(%v, %q) = %v; want nil", ctx, testresources.Bar_Alice_Payment.Name, err)
		}
		for _, test := range []struct {
			desc        string
			name        string
			wantPayment *pb.Payment
			wantErr     error
		}{
			{
				desc:        "LookupDeleted",
				name:        testresources.Bar_Alice_Payment.Name,
				wantPayment: nil,
				wantErr:     &NotFoundError{Name: testresources.Bar_Alice_Payment.Name},
			},
			{
				desc:        "LookupExisting",
				name:        testresources.Bar_Bob_Payment.Name,
				wantPayment: testresources.Bar_Bob_Payment,
				wantErr:     nil,
			},
		} {
			t.Run(test.desc, func(t *testing.T) {
				store, err := r.LookupPayment(ctx, test.name)
				if diff := cmp.Diff(store, test.wantPayment, protocmp.Transform()); diff != "" {
					t.Errorf("r.LookupPayment(%v, %q) store != test.wantPayment (-got +want)\n%s", ctx, test.name, diff)
				}
				if !cmp.Equal(err, test.wantErr, cmpopts.EquateErrors()) {
					t.Errorf("r.LookupPayment(%v, %q) err = %v; want %v", ctx, test.name, err, test.wantErr)
				}
			})
		}
	})
	t.Run("Errors", func(t *testing.T) {
		r := SeedRepository(t, []*pb.Payment{testresources.Bar_Alice_Payment})
		for _, test := range []struct {
			desc string
			name string
			want error
		}{
			{
				desc: "EmptyName",
				name: "",
				want: resourcename.ErrInvalidName,
			},
			{
				desc: "NotFound",
				name: testresources.Bar_Bob_Payment.Name,
				want: &NotFoundError{Name: testresources.Bar_Bob_Payment.Name},
			},
		} {
			t.Run(test.desc, func(t *testing.T) {
				if got := r.DeletePayment(ctx, test.name); !cmp.Equal(got, test.want, cmpopts.EquateErrors()) {
					t.Errorf("r.DeletePayment(%v, %q) = %v; want %v", ctx, test.name, got, test.want)
				}
			})
		}
	})
}
