package users

import (
	"strings"
	"testing"

	"github.com/google/uuid"
)

func TestGenerateName(t *testing.T) {
	got := GenerateName()
	prefix := CollectionID + "/"
	if !strings.HasPrefix(got, prefix) {
		t.Errorf("GenerateName() = %q; want prefix %q", got, prefix)
	}
	id := strings.TrimPrefix(got, prefix)
	if _, err := uuid.Parse(id); err != nil {
		t.Errorf("uuid.Parse(%q) err = %v; want nil", id, err)
	}
}

func TestValidateName(t *testing.T) {
	id := "6f2d193c-1460-491d-8157-7dd9535526c6"
	for _, test := range []struct {
		name string
		want error
	}{
		{name: "users/" + id, want: nil},
		{name: "", want: ErrNameInvalidFormat},
		{name: "invalidprefix/" + id, want: ErrNameInvalidFormat},
		{name: id, want: ErrNameInvalidFormat},
		{name: "users/not a UUID", want: ErrNameInvalidFormat},
	} {
		if got := ValidateName(test.name); got != test.want {
			t.Errorf("ValidateName(%q) = %v; want %v", test.name, got, test.want)
		}
	}
}
