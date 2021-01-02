package name

import "testing"

var (
	good = []string{
		"settings",
		"users/{user}",
		"stores/{store}",
		"stores/{store}/products/{product}",
	}
	bad = []string{
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
	}
)

func TestParseFormat(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		for _, s := range good {
			if _, err := ParseFormat(s); err != nil {
				t.Errorf("ParseFormat(%q) err = %v; want nil", s, err)
			}
		}
	})

	t.Run("Errors", func(t *testing.T) {
		for _, s := range bad {
			if _, err := ParseFormat(s); err == nil {
				t.Errorf("ParseFormat(%q) err = nil; want non-nil", s)
			}
		}
	})
}

func TestMustParseFormat(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		for _, s := range good {
			func() {
				defer func() {
					if err := recover(); err != nil {
						t.Errorf("MustParseFormat(%q) err = %v; want nil", s, err)
					}
				}()
				_ = MustParseFormat(s)
			}()
		}
	})

	t.Run("Errors", func(t *testing.T) {
		for _, s := range bad {
			func() {
				defer func() {
					if err := recover(); err == nil {
						t.Errorf("MustParseFormat(%q) err = nil; want non-nil", s)
					}
				}()
				_ = MustParseFormat(s)
			}()
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
