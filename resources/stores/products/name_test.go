package products

import (
	"fmt"
	"testing"

	"github.com/Saser/strecku/resourcename"
	"github.com/Saser/strecku/resources/stores"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
)

func TestGenerateName(t *testing.T) {
	store := stores.GenerateName()
	got := GenerateName(store)
	if err := ValidateName(got); err != nil {
		t.Errorf("ValidateName(GenerateName() = %q) = %v; want nil", got, err)
	}
}

func TestParseName_ValidateName(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		wantStore := uuid.MustParse("6729f7fa-dc5a-41ae-b00d-5cd67d5e1e15")
		wantProduct := uuid.MustParse("90e3eaaa-4d9c-423f-b468-bb7322fb5d4f")
		name := "stores/" + wantStore.String() + "/products/" + wantProduct.String()
		gotStore, gotProduct, err := ParseName(name)
		if err != nil {
			t.Errorf("ParseName(%q) err = %v; want nil", name, err)
		}
		if !cmp.Equal(gotStore, wantStore) {
			t.Errorf("ParseName(%q) store = %v; want %v", name, gotStore, wantStore)
		}
		if !cmp.Equal(gotProduct, wantProduct) {
			t.Errorf("ParseName(%q) product = %v; want %v", name, gotProduct, wantProduct)
		}
		if err := ValidateName(name); err != nil {
			t.Errorf("ValidateName(%q) = %v; want nil", name, err)
		}
	})

	t.Run("Errors", func(t *testing.T) {
		storeID := "6729f7fa-dc5a-41ae-b00d-5cd67d5e1e15"
		productID := "90e3eaaa-4d9c-423f-b468-bb7322fb5d4f"
		for _, name := range []string{
			"",
			productID,
			"users/" + storeID + "/products/" + productID,
			"stores/" + storeID + "/products/not-a-UUID",
			"stores/not-a-UUID/products/" + productID,
		} {
			_, _, err := ParseName(name)
			if err == nil {
				t.Errorf("ParseName(%q) err = nil; want non-nil", name)
			}
			if err := ValidateName(name); err == nil {
				t.Errorf("ValidateName(%q) = nil; want non-nil", name)
			}
		}
	})
}

func TestParent(t *testing.T) {
	storeID := "6729f7fa-dc5a-41ae-b00d-5cd67d5e1e15"
	store := fmt.Sprintf("%s/%s", stores.CollectionID, storeID)
	productID := "90e3eaaa-4d9c-423f-b468-bb7322fb5d4f"
	for _, test := range []struct {
		name       string
		wantParent string
		wantErr    error
	}{
		{
			name:       store + "/products/" + productID,
			wantParent: store,
			wantErr:    nil,
		},
		{
			name:       "users/" + storeID + "/products/" + productID,
			wantParent: "",
			wantErr:    resourcename.ErrInvalidName,
		},
	} {
		parent, err := Parent(test.name)
		if parent != test.wantParent {
			t.Errorf("Parent(%q) parent = %q; want %q", test.name, parent, test.wantParent)
		}
		if !cmp.Equal(err, test.wantErr, cmpopts.EquateErrors()) {
			t.Errorf("Parent(%q) err = %v; want %q", test.name, err, test.wantErr)
		}
	}
}
