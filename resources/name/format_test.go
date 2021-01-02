package name

import "testing"

func TestParseFormat(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		for _, s := range []string{
			"settings",
			"users/{user}",
			"stores/{store}",
			"stores/{store}/products/{product}",
		} {
			if _, err := ParseFormat(s); err != nil {
				t.Errorf("ParseFormat(%q) err = %v; want nil", s, err)
			}
		}
	})

	t.Run("Errors", func(t *testing.T) {
		for _, s := range []string{
			// Empty string.
			"",
			// Empty segments.
			"/{user}",
			"users/",
			"/",
			"stores//products",
			// Invalid characters in collection IDs.
			"camelCase",
			"with-hyphen",
			"with_underscore",
			"with1number",
			// Invalid characters in variable names.
			"{camelCase}",
			"{with-hyphen}",
			"{with_underscore}",
			"{with1number}",
		} {
			if _, err := ParseFormat(s); err == nil {
				t.Errorf("ParseFormat(%q) err = nil; want non-nil", s)
			}
		}
	})
}

func TestFormat_String(t *testing.T) {
	for _, s := range []string{
		"users/{user}",
		"stores/{store}",
		"stores/{store}/products/{product}",
	} {
		f, err := ParseFormat(s)
		if err != nil {
			t.Errorf("ParseFormat(%q) err = %v; want nil", s, err)
			continue
		}
		if got, want := f.String(), s; got != want {
			t.Errorf("f.String() = %q; want %q", got, want)
		}
	}
}
