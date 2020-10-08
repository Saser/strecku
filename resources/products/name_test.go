package products

import (
	"fmt"
	"strings"
	"testing"

	"github.com/Saser/strecku/resources/stores"
	"github.com/Saser/strecku/resources/users"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
)

func TestGenerateName(t *testing.T) {
	store := stores.GenerateName()
	got := GenerateName(store)
	wantPrefix := prefix(store)
	if !strings.HasPrefix(got, wantPrefix) {
		t.Errorf("GenerateName() = %q; want prefix %q", got, wantPrefix)
	}
	id := strings.TrimPrefix(got, wantPrefix)
	if _, err := uuid.Parse(id); err != nil {
		t.Errorf("uuid.Parse(%q) err = %v; want nil", id, err)
	}
}

func TestValidateName(t *testing.T) {
	storeID := "6729f7fa-dc5a-41ae-b00d-5cd67d5e1e15"
	store := fmt.Sprintf("%s/%s", stores.CollectionID, storeID)
	productID := "90e3eaaa-4d9c-423f-b468-bb7322fb5d4f"
	for _, test := range []struct {
		name string
		want error
	}{
		{name: fmt.Sprintf("%s/%s/%s", store, CollectionID, productID), want: nil},
		{name: "", want: ErrNameEmpty},
		{name: productID, want: ErrNameInvalidFormat},
		{name: fmt.Sprintf("%s/%s/%s/%s", users.CollectionID, storeID, CollectionID, productID), want: ErrNameInvalidFormat},
		{name: prefix(store) + "not a UUID", want: ErrNameInvalidFormat},
	} {
		if got := ValidateName(test.name); !cmp.Equal(got, test.want, cmpopts.EquateErrors()) {
			t.Errorf("ValidateName(%q) = %v; want %v", test.name, got, test.want)
		}
	}
}
