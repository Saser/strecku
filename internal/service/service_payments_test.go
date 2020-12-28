package service

import (
	"context"
	"testing"

	pb "github.com/Saser/strecku/api/v1"
	"github.com/Saser/strecku/resources/stores/payments"
	"github.com/Saser/strecku/resources/testresources"
	"github.com/google/go-cmp/cmp"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/testing/protocmp"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
)

func paymentLess(u1, u2 *pb.Payment) bool {
	return u1.Name < u2.Name
}

func TestService_GetPayment(t *testing.T) {
	ctx := context.Background()
	c := serveAndDial(ctx, t, seed(ctx, t))
	for _, test := range []struct {
		desc        string
		req         *pb.GetPaymentRequest
		wantPayment *pb.Payment
		wantCode    codes.Code
	}{
		{
			desc:        "OK",
			req:         &pb.GetPaymentRequest{Name: testresources.Bar_Alice_Payment.Name},
			wantPayment: testresources.Bar_Alice_Payment,
			wantCode:    codes.OK,
		},
		{
			desc:        "EmptyName",
			req:         &pb.GetPaymentRequest{Name: ""},
			wantPayment: nil,
			wantCode:    codes.InvalidArgument,
		},
		{
			desc:        "InvalidName",
			req:         &pb.GetPaymentRequest{Name: testresources.Alice.Name}, // name of a user
			wantPayment: nil,
			wantCode:    codes.InvalidArgument,
		},
		{
			desc:        "NotFound",
			req:         &pb.GetPaymentRequest{Name: testresources.Bar_Bob_Payment.Name},
			wantPayment: nil,
			wantCode:    codes.NotFound,
		},
	} {
		t.Run(test.desc, func(t *testing.T) {
			payment, err := c.GetPayment(ctx, test.req)
			if diff := cmp.Diff(payment, test.wantPayment, protocmp.Transform()); diff != "" {
				t.Errorf("c.GetPayment(%v, %v) payment = != test.wantPayment (-got +want)\n%s", ctx, test.req, diff)
			}
			if got := status.Code(err); got != test.wantCode {
				t.Errorf("status.Code(%v) = %v; want %v", err, got, test.wantCode)
			}
		})
	}
}

func TestService_ListPayments(t *testing.T) {
	ctx := context.Background()
	c := serveAndDial(ctx, t, seed(ctx, t))
	for _, test := range []struct {
		desc     string
		req      *pb.ListPaymentsRequest
		wantResp *pb.ListPaymentsResponse
		wantCode codes.Code
	}{
		{
			desc: "OK_Bar",
			req: &pb.ListPaymentsRequest{
				Parent:    testresources.Bar.Name,
				PageSize:  0,
				PageToken: "",
			},
			wantResp: &pb.ListPaymentsResponse{
				Payments: []*pb.Payment{
					testresources.Bar_Alice_Payment,
				},
			},
			wantCode: codes.OK,
		},
		{
			desc: "OK_Mall",
			req: &pb.ListPaymentsRequest{
				Parent:    testresources.Mall.Name,
				PageSize:  0,
				PageToken: "",
			},
			wantResp: &pb.ListPaymentsResponse{
				Payments: []*pb.Payment{
					testresources.Mall_Alice_Payment,
				},
			},
			wantCode: codes.OK,
		},
		{
			desc: "InvalidParent",
			req: &pb.ListPaymentsRequest{
				Parent:    testresources.Alice.Name, // name of a user
				PageSize:  0,
				PageToken: "",
			},
			wantResp: nil,
			wantCode: codes.InvalidArgument,
		},
		{
			desc: "NegativePageSize",
			req: &pb.ListPaymentsRequest{
				Parent:    testresources.Bar.Name,
				PageSize:  -1,
				PageToken: "",
			},
			wantResp: nil,
			wantCode: codes.InvalidArgument,
		},
		{
			desc: "PaginationUnimplemented_PositivePageSize",
			req: &pb.ListPaymentsRequest{
				Parent:    testresources.Bar.Name,
				PageSize:  1,
				PageToken: "",
			},
			wantResp: nil,
			wantCode: codes.Unimplemented,
		},
		{
			desc: "PaginationUnimplemented_NonEmptyPageToken",
			req: &pb.ListPaymentsRequest{
				Parent:    testresources.Bar.Name,
				PageSize:  0,
				PageToken: "token",
			},
			wantResp: nil,
			wantCode: codes.Unimplemented,
		},
	} {
		t.Run(test.desc, func(t *testing.T) {
			resp, err := c.ListPayments(ctx, test.req)
			if diff := cmp.Diff(
				resp, test.wantResp, protocmp.Transform(),
				protocmp.FilterField(new(pb.ListPaymentsResponse), "payments", protocmp.SortRepeated(paymentLess)),
			); diff != "" {
				t.Errorf("c.ListPayments(%v, %v) resp != test.wantResp (-got +want)\n%s", ctx, test.req, diff)
			}
			if got := status.Code(err); got != test.wantCode {
				t.Errorf("status.Code(%v) = %v; want %v", err, got, test.wantCode)
			}
		})
	}
}

func TestService_CreatePayment(t *testing.T) {
	ctx := context.Background()
	for _, test := range []struct {
		desc        string
		req         *pb.CreatePaymentRequest
		wantPayment *pb.Payment
		wantCode    codes.Code
	}{
		{
			desc: "OK",
			req: &pb.CreatePaymentRequest{
				Parent:  testresources.Bar.Name,
				Payment: testresources.Bar_Bob_Payment,
			},
			wantPayment: testresources.Bar_Bob_Payment,
			wantCode:    codes.OK,
		},
	} {
		t.Run(test.desc, func(t *testing.T) {
			c := serveAndDial(ctx, t, seed(ctx, t))
			payment, err := c.CreatePayment(ctx, test.req)
			if diff := cmp.Diff(
				payment, test.wantPayment, protocmp.Transform(),
				protocmp.IgnoreFields(new(pb.Payment), "name"),
			); diff != "" {
				t.Errorf("c.CreatePayment(%v, %v) payment != test.wantPayment (-got +want)\n%s", ctx, test.req, diff)
			}
			if got := status.Code(err); got != test.wantCode {
				t.Errorf("status.Code(%v) = %v; want %v", err, got, test.wantCode)
			}
		})
	}
}

func TestService_UpdatePayment(t *testing.T) {
	ctx := context.Background()
	// Test scenario(s) where the update is successful.
	t.Run("OK", func(t *testing.T) {
		oldPayment := payments.Clone(testresources.Bar_Alice_Payment)
		newPayment := payments.Clone(oldPayment)
		newPayment.Description = "Alice's new payment"
		newPayment.AmountCents = 50000
		for _, test := range []struct {
			desc string
			req  *pb.UpdatePaymentRequest
			want *pb.Payment
		}{
			{
				desc: "NoOp_NilUpdateMask",
				req: &pb.UpdatePaymentRequest{
					Payment:    oldPayment,
					UpdateMask: nil,
				},
				want: oldPayment,
			},
			{
				desc: "NoOp_AllPaths",
				req: &pb.UpdatePaymentRequest{
					Payment:    oldPayment,
					UpdateMask: &fieldmaskpb.FieldMask{Paths: []string{"description", "amount_cents"}},
				},
				want: oldPayment,
			},
			{
				desc: "NoOp_NoPaths",
				req: &pb.UpdatePaymentRequest{
					Payment:    newPayment,
					UpdateMask: &fieldmaskpb.FieldMask{Paths: nil},
				},
				want: oldPayment,
			},
			{
				desc: "FullUpdate_NilUpdateMask",
				req: &pb.UpdatePaymentRequest{
					Payment:    newPayment,
					UpdateMask: nil,
				},
				want: newPayment,
			},
			{
				desc: "FullUpdate_AllPaths",
				req: &pb.UpdatePaymentRequest{
					Payment:    newPayment,
					UpdateMask: &fieldmaskpb.FieldMask{Paths: []string{"description", "amount_cents"}},
				},
				want: newPayment,
			},
			{
				desc: "PartialUpdate_FullPayment_Description",
				req: &pb.UpdatePaymentRequest{
					Payment:    newPayment,
					UpdateMask: &fieldmaskpb.FieldMask{Paths: []string{"description"}},
				},
				want: func() *pb.Payment {
					payment := payments.Clone(oldPayment)
					payment.Description = newPayment.Description
					return payment
				}(),
			},
			{
				desc: "PartialUpdate_FullPayment_AmountCents",
				req: &pb.UpdatePaymentRequest{
					Payment:    newPayment,
					UpdateMask: &fieldmaskpb.FieldMask{Paths: []string{"amount_cents"}},
				},
				want: func() *pb.Payment {
					payment := payments.Clone(oldPayment)
					payment.AmountCents = newPayment.AmountCents
					return payment
				}(),
			},
			{
				desc: "PartialUpdate_PartialPayment_Description",
				req: &pb.UpdatePaymentRequest{
					Payment: &pb.Payment{
						Name:        oldPayment.Name,
						Description: "New payment",
					},
					UpdateMask: &fieldmaskpb.FieldMask{Paths: []string{"description"}},
				},
				want: func() *pb.Payment {
					payment := payments.Clone(oldPayment)
					payment.Description = "New payment"
					return payment
				}(),
			},
			{
				desc: "PartialUpdate_PartialPayment_AmountCents",
				req: &pb.UpdatePaymentRequest{
					Payment: &pb.Payment{
						Name:        oldPayment.Name,
						AmountCents: 50000,
					},
					UpdateMask: &fieldmaskpb.FieldMask{Paths: []string{"amount_cents"}},
				},
				want: func() *pb.Payment {
					payment := payments.Clone(oldPayment)
					payment.AmountCents = 50000
					return payment
				}(),
			},
		} {
			t.Run(test.desc, func(t *testing.T) {
				c := serveAndDial(ctx, t, seed(ctx, t))
				payment, err := c.UpdatePayment(ctx, test.req)
				if diff := cmp.Diff(payment, test.want, protocmp.Transform()); diff != "" {
					t.Errorf("c.UpdatePayment(%v, %v) payment != test.want (-got +want)\n%s", ctx, test.req, diff)
				}
				if got, want := status.Code(err), codes.OK; got != want {
					t.Errorf("status.Code(%v) = %v; want %v", err, got, want)
				}
			})
		}
	})
	// Test scenario(s) where the update fails.
	t.Run("Errors", func(t *testing.T) {
		for _, test := range []struct {
			desc string
			req  *pb.UpdatePaymentRequest
			want codes.Code
		}{
			{
				desc: "UpdateUser",
				req: &pb.UpdatePaymentRequest{
					Payment: func() *pb.Payment {
						payment := payments.Clone(testresources.Bar_Alice_Payment)
						payment.User = testresources.Bob.Name
						return payment
					}(),
					UpdateMask: nil,
				},
				want: codes.InvalidArgument,
			},
			{
				desc: "NegativeAmountCents",
				req: &pb.UpdatePaymentRequest{
					Payment: func() *pb.Payment {
						payment := payments.Clone(testresources.Bar_Alice_Payment)
						payment.AmountCents = -10000
						return payment
					}(),
					UpdateMask: nil,
				},
				want: codes.InvalidArgument,
			},
			{
				desc: "NotFound",
				req: &pb.UpdatePaymentRequest{
					Payment:    testresources.Bar_Bob_Payment,
					UpdateMask: nil,
				},
				want: codes.NotFound,
			},
			{
				desc: "InvalidUpdateMask",
				req: &pb.UpdatePaymentRequest{
					Payment:    testresources.Bar_Alice_Payment,
					UpdateMask: &fieldmaskpb.FieldMask{Paths: []string{"invalid"}},
				},
				want: codes.InvalidArgument,
			},
		} {
			t.Run(test.desc, func(t *testing.T) {
				c := serveAndDial(ctx, t, seed(ctx, t))
				_, err := c.UpdatePayment(ctx, test.req)
				if got := status.Code(err); got != test.want {
					t.Errorf("status.Code(%v) = %v; want %v", err, got, test.want)
				}
			})
		}
	})
}

func TestService_DeletePayment(t *testing.T) {
	ctx := context.Background()
	// Test scenario(s) where the delete is successful.
	t.Run("OK", func(t *testing.T) {
		c := serveAndDial(ctx, t, seed(ctx, t))
		{
			req := &pb.DeletePaymentRequest{Name: testresources.Bar_Alice_Payment.Name}
			_, err := c.DeletePayment(ctx, req)
			if got, want := status.Code(err), codes.OK; got != want {
				t.Errorf("status.Code(%v) = %v; want %v", err, got, want)
			}
		}
		{
			req := &pb.GetPaymentRequest{Name: testresources.Bar_Alice_Payment.Name}
			_, err := c.GetPayment(ctx, req)
			if got, want := status.Code(err), codes.NotFound; got != want {
				t.Errorf("status.Code(%v) = %v; want %v", err, got, want)
			}
		}
	})
	// Test scenario(s) where the delete fails.
	t.Run("Errors", func(t *testing.T) {
		c := serveAndDial(ctx, t, seed(ctx, t))
		for _, test := range []struct {
			desc string
			req  *pb.DeletePaymentRequest
			want codes.Code
		}{
			{
				desc: "EmptyName",
				req:  &pb.DeletePaymentRequest{Name: ""},
				want: codes.InvalidArgument,
			},
			{
				desc: "NotFound",
				req:  &pb.DeletePaymentRequest{Name: testresources.Bar_Bob_Payment.Name},
				want: codes.NotFound,
			},
		} {
			t.Run(test.desc, func(t *testing.T) {
				_, err := c.DeletePayment(ctx, test.req)
				if got := status.Code(err); got != test.want {
					t.Errorf("status.Code(%v) = %v; want %v", err, got, test.want)
				}
			})
		}
	})
}
