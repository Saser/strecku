package stores

import (
	"testing"

	streckuv1 "github.com/Saser/strecku/saser/strecku/v1"
)

func TestValidate(t *testing.T) {
	for _, test := range []struct {
		store *streckuv1.Store
		want  error
	}{
		{store: &streckuv1.Store{Name: "stores/6f2d193c-1460-491d-8157-7dd9535526c6", DisplayName: "Store"}, want: nil},
		{store: &streckuv1.Store{Name: "", DisplayName: "Store"}, want: ErrNameEmpty},
		{store: &streckuv1.Store{Name: "invalidprefix/6f2d193c-1460-491d-8157-7dd9535526c6", DisplayName: "Store"}, want: ErrNameInvalidFormat},
		{store: &streckuv1.Store{Name: "6f2d193c-1460-491d-8157-7dd9535526c6", DisplayName: "Store"}, want: ErrNameInvalidFormat},
		{store: &streckuv1.Store{Name: "stores/6f2d193c-1460-491d-8157-7dd9535526c6", DisplayName: ""}, want: ErrDisplayNameEmpty},
	} {
		if got := Validate(test.store); got != test.want {
			t.Errorf("Validate(%v) = %v; want %v", test.store, got, test.want)
		}
	}
}

func TestValidateName(t *testing.T) {
	for _, test := range []struct {
		name string
		want error
	}{
		{name: "stores/6f2d193c-1460-491d-8157-7dd9535526c6", want: nil},
		{name: "", want: ErrNameEmpty},
		{name: "invalidprefix/6f2d193c-1460-491d-8157-7dd9535526c6", want: ErrNameInvalidFormat},
		{name: "6f2d193c-1460-491d-8157-7dd9535526c6", want: ErrNameInvalidFormat},
		{name: "stores/not a UUID", want: ErrNameInvalidFormat},
	} {
		if got := ValidateName(test.name); got != test.want {
			t.Errorf("ValidateName(%q) = %v; want %v", test.name, got, test.want)
		}
	}
}
