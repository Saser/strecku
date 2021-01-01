package name

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestResourceID_Validate(t *testing.T) {
	t.Run("Valid", func(t *testing.T) {
		r := ResourceID("0464cc9f-506f-415b-9a3f-fb7305026781")
		if err := r.Validate(); err != nil {
			t.Errorf("ResourceID(%q).Validate() = %v; want nil", r, err)
		}
	})

	t.Run("Invalid", func(t *testing.T) {
		for _, r := range []ResourceID{
			"0464CC9F-506F-415B-9A3F-FB7305026781",
			"0464cc9f506f415b9a3ffb7305026781",
			"urn:uuid:0464cc9f-506f-415b-9a3f-fb7305026781",
			"{0464cc9f-506f-415b-9a3f-fb7305026781}",
		} {
			if got, want := r.Validate(), InvalidResourceID(r); !cmp.Equal(got, want, cmpopts.EquateErrors()) {
				t.Errorf("ResourceID(%q).Validate() = %v; want %v", r, got, want)
			}
		}
	})
}
