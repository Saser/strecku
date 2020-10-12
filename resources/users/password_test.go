package users

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestValidatePassword(t *testing.T) {
	for _, test := range []struct {
		password string
		want     error
	}{
		{
			password: "",
			want:     ErrPasswordEmpty,
		},
	} {
		if got := ValidatePassword(test.password); !cmp.Equal(got, test.want, cmpopts.EquateErrors()) {
			t.Errorf("ValidatePassword(%q) = %v; want %v", test.password, got, test.want)
		}
	}
}
