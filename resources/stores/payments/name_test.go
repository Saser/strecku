package payments

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
	prefix := store + "/" + CollectionID + "/"
	if !strings.HasPrefix(got, prefix) {
		t.Errorf("GenerateName() = %q; want prefix %q", got, prefix)
	}
	id := strings.TrimPrefix(got, prefix)
	if _, err := uuid.Parse(id); err != nil {
		t.Errorf("uuid.Parse(%q) err = %v; want nil", id, err)
	}
}

func TestValidateName(t *testing.T) {
	storeID := "6729f7fa-dc5a-41ae-b00d-5cd67d5e1e15"
	store := fmt.Sprintf("%s/%s", stores.CollectionID, storeID)
	paymentID := "90e3eaaa-4d9c-423f-b468-bb7322fb5d4f"
	for _, test := range []struct {
		name string
		want error
	}{
		{name: fmt.Sprintf("%s/%s/%s", store, CollectionID, paymentID), want: nil},
		{name: "", want: ErrNameInvalidFormat},
		{name: paymentID, want: ErrNameInvalidFormat},
		{name: fmt.Sprintf("%s/%s/%s/%s", users.CollectionID, storeID, CollectionID, paymentID), want: ErrNameInvalidFormat},
		{name: store + "/" + CollectionID + "/not a UUID", want: ErrNameInvalidFormat},
	} {
		if got := ValidateName(test.name); !cmp.Equal(got, test.want, cmpopts.EquateErrors()) {
			t.Errorf("ValidateName(%q) = %v; want %v", test.name, got, test.want)
		}
	}
}

func TestParent(t *testing.T) {
	storeID := "6729f7fa-dc5a-41ae-b00d-5cd67d5e1e15"
	store := fmt.Sprintf("%s/%s", stores.CollectionID, storeID)
	paymentID := "90e3eaaa-4d9c-423f-b468-bb7322fb5d4f"
	for _, test := range []struct {
		name       string
		wantParent string
		wantErr    error
	}{
		{
			name:       store + "/" + CollectionID + "/" + paymentID,
			wantParent: store,
			wantErr:    nil,
		},
		{
			name:       users.CollectionID + "/" + storeID + "/" + CollectionID + "/" + paymentID,
			wantParent: "",
			wantErr:    ErrNameInvalidFormat,
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
