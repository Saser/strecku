package names

import (
	"regexp"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestUUID(t *testing.T) {
	re, err := regexp.Compile(UUID)
	if err != nil {
		t.Fatalf("regexp.Compile(%q) err = %v; want nil", UUID, err)
	}
	for _, test := range []struct {
		want bool
		s    string
	}{
		// Strings that should be matched.
		{want: true, s: "ef299eb2-2f9e-41fd-8f41-24a578d6e1f2"}, // standard format for UUIDs
		{want: true, s: "ef299eb22f9e41fd8f4124a578d6e1f2"},     // same, without hyphens
		{want: true, s: "EF299EB2-2F9E-41FD-8F41-24A578D6E1F2"}, // standard format for UUIDs, only uppercase
		{want: true, s: "EF299EB22F9E41FD8F4124A578D6E1F2"},     // same, without hyphens
		{want: true, s: "eF299Eb2-2f9e-41fD-8f41-24a578d6e1f2"}, // standard format for UUIDs, mixed case
		{want: true, s: "eF299Eb22f9e41fD8f4124a578d6e1f2"},     // same, without hyphens
		// Strings that should not be matched.
		{want: false, s: "not a UUID"},
		{want: false, s: "kl299kh2-2l9k-41lj-8l41-24g578j6k1l2"}, // looks like a UUID, but non-hex characters
		{want: false, s: "ef299eb2-2f9e-41fd-8f41-24a578d6e1f"},  // a UUID missing the last character
		{want: false, s: "ef299eb22f9e41fd8f4124a578d6e1f"},      // same, without hyphens
	} {
		if got := re.MatchString(test.s); got != test.want {
			t.Errorf("re.MatchString(%q) = %v; want %v", test.s, got, test.want)
		}
	}
}

func TestCompile(t *testing.T) {
	for _, test := range []struct {
		parts []string
		want  error
	}{
		{
			parts: []string{UUID},
			want:  nil,
		},
		{
			parts: []string{"foos", UUID},
			want:  nil,
		},
		{
			parts: []string{"foos", UUID, "bars", UUID},
			want:  nil,
		},
		{
			parts: []string{"foos", UUID, "bars", "bazs"},
			want:  nil,
		},
	} {
		if _, got := Compile(test.parts...); !cmp.Equal(got, test.want, cmpopts.EquateErrors()) {
			t.Errorf("Compile(%q) = %v; want %v", test.parts, got, test.want)
		}
	}
}
