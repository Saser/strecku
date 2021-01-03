package stores

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
)

func TestGenerateName(t *testing.T) {
	got := GenerateName()
	if err := ValidateName(got); err != nil {
		t.Errorf("ValidateName(GenerateName() = %q) = %v; want nil", got, err)
	}
}

func TestParseName_ValidateName(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		want := uuid.MustParse("6f2d193c-1460-491d-8157-7dd9535526c6")
		name := "stores/" + want.String()
		got, err := ParseName(name)
		if err != nil {
			t.Errorf("ParseName(%q) err = %v; want nil", name, err)
		}
		if !cmp.Equal(got, want) {
			t.Errorf("ParseName(%q) uuid = %v; want %v", name, got, want)
		}
		if err := ValidateName(name); err != nil {
			t.Errorf("ValidateName(%q) = %v; want nil", name, err)
		}
	})

	t.Run("Errors", func(t *testing.T) {
		id := "6f2d193c-1460-491d-8157-7dd9535526c6"
		for _, s := range []string{
			"",
			"invalidprefix/" + id,
			id,
			"stores/not a UUID",
		} {
			_, err := ParseName(s)
			if err == nil {
				t.Errorf("ParseName(%q) err = nil; want non-nil", s)
			}
			if err := ValidateName(s); err == nil {
				t.Errorf("ValidateName(%q) = nil; want non-nil", s)
			}
		}
	})
}
