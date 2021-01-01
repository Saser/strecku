package name

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestCollectionID_Validate(t *testing.T) {
	t.Run("Valid", func(t *testing.T) {
		for _, c := range []CollectionID{
			"memberships",
			"payments",
			"products",
			"purchases",
			"stores",
			"users",
		} {
			if err := c.Validate(); err != nil {
				t.Errorf("Collection(%q).Validate() = %v; want nil", c, err)
			}
		}
	})

	t.Run("Invalid", func(t *testing.T) {
		for _, c := range []CollectionID{
			"camelCase",
			"with-hyphens",
			"with1number",
		} {
			if got, want := c.Validate(), InvalidCollectionID(c); !cmp.Equal(got, want, cmpopts.EquateErrors()) {
				t.Errorf("Collection(%q).Validate() = %v; want %v", c, got, want)
			}
		}
	})
}
