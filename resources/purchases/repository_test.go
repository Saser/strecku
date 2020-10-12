package purchases

import (
	"context"
	"fmt"
	"testing"

	pb "github.com/Saser/strecku/api/v1"
	"github.com/Saser/strecku/resources/testresources"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"google.golang.org/protobuf/testing/protocmp"
)

func purchaseLess(p1, p2 *pb.Purchase) bool {
	return p1.Name < p2.Name
}

func purchaseLineLess(l1, l2 *pb.Purchase_Line) bool {
	return l1.Description < l2.Description
}

func TestNotFoundError_Error(t *testing.T) {
	err := &NotFoundError{Name: testresources.Alice_Beer1.Name}
	want := fmt.Sprintf("purchase not found: %q", testresources.Alice_Beer1.Name)
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
			err:    &NotFoundError{Name: testresources.Alice_Beer1.Name},
			target: &NotFoundError{Name: testresources.Alice_Beer1.Name},
			want:   true,
		},
		{
			err:    &NotFoundError{Name: testresources.Alice_Beer1.Name},
			target: &NotFoundError{Name: testresources.Alice_Cocktail1.Name},
			want:   false,
		},
		{
			err:    &NotFoundError{Name: testresources.Alice_Beer1.Name},
			target: fmt.Errorf("purchase not found: %q", testresources.Alice_Beer1.Name),
			want:   false,
		},
	} {
		if got := test.err.Is(test.target); got != test.want {
			t.Errorf("test.err.Is(%v) = %v; want %v", test.target, got, test.want)
		}
	}
}

func TestExistsError_Error(t *testing.T) {
	err := &ExistsError{Name: testresources.Beer.Name}
	want := fmt.Sprintf("purchase exists: %q", testresources.Beer.Name)
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
			err:    &ExistsError{Name: testresources.Alice_Beer1.Name},
			target: &ExistsError{Name: testresources.Alice_Beer1.Name},
			want:   true,
		},
		{
			err:    &ExistsError{Name: testresources.Alice_Beer1.Name},
			target: &ExistsError{Name: testresources.Alice_Cocktail1.Name},
			want:   false,
		},
		{
			err:    &ExistsError{Name: testresources.Alice_Beer1.Name},
			target: fmt.Errorf("purchase exists: %q", testresources.Alice_Beer1.Name),
			want:   false,
		},
	} {
		if got := test.err.Is(test.target); got != test.want {
			t.Errorf("test.err.Is(%v) = %v; want %v", test.target, got, test.want)
		}
	}
}

func TestRepository_LookupPurchase(t *testing.T) {
	ctx := context.Background()
	r := SeedRepository(t, []*pb.Purchase{testresources.Alice_Beer1})
	for _, test := range []struct {
		desc         string
		name         string
		wantPurchase *pb.Purchase
		wantErr      error
	}{
		{
			desc:         "OK",
			name:         testresources.Alice_Beer1.Name,
			wantPurchase: testresources.Alice_Beer1,
			wantErr:      nil,
		},
		{
			desc:         "EmptyName",
			name:         "",
			wantPurchase: nil,
			wantErr:      ErrNameEmpty,
		},
		{
			desc:         "NotFound",
			name:         testresources.Alice_Cocktail1.Name,
			wantPurchase: nil,
			wantErr:      &NotFoundError{Name: testresources.Alice_Cocktail1.Name},
		},
	} {
		t.Run(test.desc, func(t *testing.T) {
			purchase, err := r.LookupPurchase(ctx, test.name)
			if diff := cmp.Diff(purchase, test.wantPurchase, protocmp.Transform()); diff != "" {
				t.Errorf("r.LookupPurchase(%v, %q) purchase != test.wantPurchase (-got +want)\n%s", ctx, test.name, diff)
			}
			if !cmp.Equal(err, test.wantErr, cmpopts.EquateErrors()) {
				t.Errorf("r.LookupPurchase(%v, %q) err = %v; want %v", ctx, test.name, err, test.wantErr)
			}
		})
	}
}

func TestRepository_ListPurchases(t *testing.T) {
	ctx := context.Background()
	allPurchases := []*pb.Purchase{
		testresources.Alice_Beer1,
		testresources.Alice_Cocktail1,
		testresources.Alice_Beer2_Cocktail2,
	}
	r := SeedRepository(t, allPurchases)
	purchases, err := r.ListPurchases(ctx)
	if diff := cmp.Diff(
		purchases, allPurchases, protocmp.Transform(),
		cmpopts.EquateEmpty(),
		cmpopts.SortSlices(purchaseLess),
		protocmp.FilterField(new(pb.Purchase), "lines", protocmp.SortRepeated(purchaseLineLess)),
	); diff != "" {
		t.Errorf("r.ListPurchases(%v) purchases != allPurchases (-got +want)\n%s", ctx, diff)
	}
	if err != nil {
		t.Errorf("r.ListPurchases(%v) err = %v; want nil", ctx, err)
	}
}

func TestRepository_FilterPurchases(t *testing.T) {
	ctx := context.Background()
	r := SeedRepository(t, []*pb.Purchase{
		testresources.Alice_Beer1,
		testresources.Alice_Cocktail1,
		testresources.Alice_Beer2_Cocktail2,
	})
	for _, test := range []struct {
		desc      string
		predicate func(*pb.Purchase) bool
		want      []*pb.Purchase
	}{
		{
			desc:      "NoneMatching",
			predicate: func(*pb.Purchase) bool { return false },
			want:      nil,
		},
		{
			desc:      "OneMatching",
			predicate: func(purchase *pb.Purchase) bool { return purchase.Name == testresources.Alice_Beer1.Name },
			want: []*pb.Purchase{
				testresources.Alice_Beer1,
			},
		},
		{
			desc:      "MultipleMatching",
			predicate: func(purchase *pb.Purchase) bool { return len(purchase.Lines) == 1 },
			want: []*pb.Purchase{
				testresources.Alice_Beer1,
				testresources.Alice_Cocktail1,
			},
		},
	} {
		t.Run(test.desc, func(t *testing.T) {
			filtered, err := r.FilterPurchases(ctx, test.predicate)
			if diff := cmp.Diff(
				filtered, test.want, protocmp.Transform(),
				cmpopts.EquateEmpty(),
				cmpopts.SortSlices(purchaseLess),
				protocmp.FilterField(new(pb.Purchase), "lines", protocmp.SortRepeated(purchaseLineLess)),
			); diff != "" {
				t.Errorf("r.FilterPurchases(%v, test.predicate) filtered != test.want (-got +want)\n%s", ctx, diff)
			}
			if err != nil {
				t.Errorf("r.FilterPurchases(%v, test.predicate) err = %v; want nil", ctx, err)
			}
		})
	}
}

func TestRepository_CreatePurchase(t *testing.T) {
	ctx := context.Background()
	for _, test := range []struct {
		desc     string
		purchase *pb.Purchase
		want     error
	}{
		{
			desc:     "OK",
			purchase: testresources.Alice_Cocktail1,
			want:     nil,
		},
		{
			desc: "DuplicateName",
			purchase: func() *pb.Purchase {
				purchase := Clone(testresources.Alice_Cocktail1)
				purchase.Name = testresources.Alice_Beer1.Name
				return purchase
			}(),
			want: &ExistsError{Name: testresources.Alice_Beer1.Name},
		},
		{
			desc: "InvalidPurchase",
			purchase: func() *pb.Purchase {
				// Create a purchase with no lines.
				// This type of invalidity was chosen arbitrarily.
				purchase := Clone(testresources.Alice_Cocktail1)
				purchase.Lines = nil
				return purchase
			}(),
			want: ErrLinesEmpty,
		},
	} {
		t.Run(test.desc, func(t *testing.T) {
			r := SeedRepository(t, []*pb.Purchase{testresources.Alice_Beer1})
			if got := r.CreatePurchase(ctx, test.purchase); !cmp.Equal(got, test.want, cmpopts.EquateErrors()) {
				t.Errorf("r.CreatePurchase(%v, %v) = %v; want %v", ctx, test.purchase, got, test.want)
			}
		})
	}
}

func TestRepository_UpdatePurchase(t *testing.T) {
	ctx := context.Background()
	// Test scenario(s) where the update is successful.
	t.Run("OK", func(t *testing.T) {
		for _, test := range []struct {
			desc   string
			modify func(purchase *pb.Purchase)
		}{
			{
				desc: "UpdateLine",
				modify: func(purchase *pb.Purchase) {
					line := purchase.Lines[0]
					line.Description = "Cider"
					line.Quantity = 2
					line.PriceCents = -1500
					line.Product = ""
				},
			},
			{
				desc: "AppendLine",
				modify: func(purchase *pb.Purchase) {
					purchase.Lines = append(purchase.Lines, &pb.Purchase_Line{
						Description: testresources.Cocktail.DisplayName,
						Quantity:    1,
						PriceCents:  testresources.Cocktail.FullPriceCents,
						Product:     testresources.Cocktail.Name,
					})
				},
			},
		} {
			t.Run(test.desc, func(t *testing.T) {
				r := SeedRepository(t, []*pb.Purchase{testresources.Alice_Beer1})
				oldPurchase := Clone(testresources.Alice_Beer1)
				newPurchase := Clone(oldPurchase)
				test.modify(newPurchase)
				if err := r.UpdatePurchase(ctx, newPurchase); err != nil {
					t.Errorf("r.UpdatePurchase(%v, %v) = %v; want nil", ctx, newPurchase, err)
				}
				purchase, err := r.LookupPurchase(ctx, newPurchase.Name)
				if diff := cmp.Diff(
					purchase, newPurchase, protocmp.Transform(),
					protocmp.FilterField(new(pb.Purchase), "lines", protocmp.SortRepeated(purchaseLineLess)),
				); diff != "" {
					t.Errorf("r.LookupPurchase(%v, %q) purchase != newPurchase (-got +want)\n%s", ctx, newPurchase.Name, diff)
				}
				if err != nil {
					t.Errorf("r.LookupPurchase(%v, %q) err = %v; want nil", ctx, newPurchase.Name, err)
				}
			})
		}
	})
	// Test scenario(s) where the update failed.
	t.Run("Errors", func(t *testing.T) {
		r := SeedRepository(t, []*pb.Purchase{testresources.Alice_Beer1})
		for _, test := range []struct {
			desc   string
			modify func(purchase *pb.Purchase)
			want   error
		}{
			{
				desc:   "UpdateUser",
				modify: func(purchase *pb.Purchase) { purchase.User = testresources.Bob.Name },
				want:   ErrUpdateUser,
			},
			{
				desc: "UpdateStore",
				modify: func(purchase *pb.Purchase) {
					purchase.Store = testresources.Mall.Name
					// The store is updated, so the product must be modified to be either unset or a
					// product that belongs to the new store (otherwise it is not a valid purchase).
					// Unsetting seems simplest.
					purchase.Lines[0].Product = ""
				},
				want: ErrUpdateStore,
			},
			{
				desc:   "NotFound",
				modify: func(purchase *pb.Purchase) { purchase.Name = testresources.Alice_Cocktail1.Name },
				want:   &NotFoundError{Name: testresources.Alice_Cocktail1.Name},
			},
		} {
			t.Run(test.desc, func(t *testing.T) {
				oldPurchase := Clone(testresources.Alice_Beer1)
				newPurchase := Clone(oldPurchase)
				test.modify(newPurchase)
				if got := r.UpdatePurchase(ctx, newPurchase); !cmp.Equal(got, test.want, cmpopts.EquateErrors()) {
					t.Errorf("r.UpdatePurchase(%v, %v) = %v; want %v", ctx, newPurchase, got, test.want)
				}
				purchase, err := r.LookupPurchase(ctx, oldPurchase.Name)
				if diff := cmp.Diff(
					purchase, oldPurchase, protocmp.Transform(),
					protocmp.FilterField(new(pb.Purchase), "lines", protocmp.SortRepeated(purchaseLineLess)),
				); diff != "" {
					t.Errorf("r.LookupPurchase(%v, %q) purchase != oldPurchase (-got +want)\n%s", ctx, oldPurchase.Name, diff)
				}
				if err != nil {
					t.Errorf("r.LookupPurchase(%v, %q) err = %v; want nil", ctx, oldPurchase.Name, err)
				}
			})
		}
	})
}

func TestRepository_DeletePurchase(t *testing.T) {
	ctx := context.Background()
	t.Run("OK", func(t *testing.T) {
		r := SeedRepository(t, []*pb.Purchase{
			testresources.Alice_Beer1,
			testresources.Alice_Cocktail1,
		})
		if err := r.DeletePurchase(ctx, testresources.Alice_Beer1.Name); err != nil {
			t.Errorf("r.DeletePurchase(%v, %q) = %v; want nil", ctx, testresources.Alice_Beer1.Name, err)
		}
		for _, test := range []struct {
			desc         string
			name         string
			wantPurchase *pb.Purchase
			wantErr      error
		}{
			{
				desc:         "LookupDeleted",
				name:         testresources.Alice_Beer1.Name,
				wantPurchase: nil,
				wantErr:      &NotFoundError{Name: testresources.Alice_Beer1.Name},
			},
			{
				desc:         "LookupExisting",
				name:         testresources.Alice_Cocktail1.Name,
				wantPurchase: testresources.Alice_Cocktail1,
				wantErr:      nil,
			},
		} {
			t.Run(test.desc, func(t *testing.T) {
				store, err := r.LookupPurchase(ctx, test.name)
				if diff := cmp.Diff(store, test.wantPurchase, protocmp.Transform()); diff != "" {
					t.Errorf("r.LookupPurchase(%v, %q) store != test.wantPurchase (-got +want)\n%s", ctx, test.name, diff)
				}
				if !cmp.Equal(err, test.wantErr, cmpopts.EquateErrors()) {
					t.Errorf("r.LookupPurchase(%v, %q) err = %v; want %v", ctx, test.name, err, test.wantErr)
				}
			})
		}
	})
	t.Run("Errors", func(t *testing.T) {
		r := SeedRepository(t, []*pb.Purchase{testresources.Alice_Beer1})
		for _, test := range []struct {
			desc string
			name string
			want error
		}{
			{
				desc: "EmptyName",
				name: "",
				want: ErrNameEmpty,
			},
			{
				desc: "NotFound",
				name: testresources.Alice_Cocktail1.Name,
				want: &NotFoundError{Name: testresources.Alice_Cocktail1.Name},
			},
		} {
			t.Run(test.desc, func(t *testing.T) {
				if got := r.DeletePurchase(ctx, test.name); !cmp.Equal(got, test.want, cmpopts.EquateErrors()) {
					t.Errorf("r.DeletePurchase(%v, %q) = %v; want %v", ctx, test.name, got, test.want)
				}
			})
		}
	})
}
