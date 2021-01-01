package name

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestName_Validate(t *testing.T) {
	t.Run("Valid", func(t *testing.T) {
		for _, n := range []Name{
			"users/e71ee48e-1469-4a52-a338-8f78fe1c6cf7",
			"stores/0464cc9f-506f-415b-9a3f-fb7305026781/users/e71ee48e-1469-4a52-a338-8f78fe1c6cf7",
		} {
			if err := n.Validate(); err != nil {
				t.Errorf("Name(%q).Validate() = %v; want nil", n, err)
			}
		}
	})

	t.Run("Invalid", func(t *testing.T) {
		for _, test := range []struct {
			n    Name
			want error
		}{
			{
				n:    Name("camelCase/e71ee48e-1469-4a52-a338-8f78fe1c6cf7"),
				want: InvalidCollectionID("camelCase"),
			},
			{
				n:    Name("stores/0464cc9f-506f-415b-9a3f-fb7305026781/camelCase/e71ee48e-1469-4a52-a338-8f78fe1c6cf7"),
				want: InvalidCollectionID("camelCase"),
			},
			{
				n:    Name("users/e71ee48e14694a52a3388f78fe1c6cf7"),
				want: InvalidResourceID("e71ee48e14694a52a3388f78fe1c6cf7"),
			},
			{
				n:    Name("stores/0464cc9f-506f-415b-9a3f-fb7305026781/users/e71ee48e14694a52a3388f78fe1c6cf7"),
				want: InvalidResourceID("e71ee48e14694a52a3388f78fe1c6cf7"),
			},
		} {
			if got := test.n.Validate(); !cmp.Equal(got, test.want, cmpopts.EquateErrors()) {
				t.Errorf("Name(%q).Validate() = %v; want %v", test.n, got, test.want)
			}
		}
	})
}

func TestName_ResourceIDs(t *testing.T) {
	for _, test := range []struct {
		n    Name
		want map[CollectionID]ResourceID
	}{
		{
			n: Name("users/e71ee48e-1469-4a52-a338-8f78fe1c6cf7"),
			want: map[CollectionID]ResourceID{
				"users": "e71ee48e-1469-4a52-a338-8f78fe1c6cf7",
			},
		},
		{
			n: Name("stores/0464cc9f-506f-415b-9a3f-fb7305026781/users/e71ee48e-1469-4a52-a338-8f78fe1c6cf7"),
			want: map[CollectionID]ResourceID{
				"stores": "0464cc9f-506f-415b-9a3f-fb7305026781",
				"users":  "e71ee48e-1469-4a52-a338-8f78fe1c6cf7",
			},
		},
	} {
		got := test.n.ResourceIDs()
		less := func(c1, c2 CollectionID) bool {
			return c1 < c2
		}
		if diff := cmp.Diff(test.want, got, cmpopts.SortMaps(less)); diff != "" {
			t.Errorf("Name(%q).ResourceIDs() diff (-want +got)\n%s", test.n, diff)
		}
	}
}
