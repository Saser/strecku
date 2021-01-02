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

func TestFormat_Parse(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		for _, test := range []struct {
			f    *Format
			name string
		}{
			{
				f:    MustParseFormat("users/{user}"),
				name: "users/78da9161-aef1-49ed-bc92-0f136c95308f",
			},
			{
				f:    MustParseFormat("stores/{store}/products/{product}"),
				name: "stores/78da9161-aef1-49ed-bc92-0f136c95308f/products/1bba3dfe-5770-4d65-ae3f-1bff9e45b668",
			},
		} {
			if _, err := test.f.Parse(test.name); err != nil {
				t.Errorf("f.Parse(%q) err = %v; want nil", test.name, err)
			}
		}
	})

	t.Run("Errors", func(t *testing.T) {
		for _, test := range []struct {
			f    *Format
			name string
		}{
			{
				f:    MustParseFormat("users/{user}"),
				name: "users",
			},
			{
				f:    MustParseFormat("users/{user}"),
				name: "users/",
			},
			{
				f:    MustParseFormat("users/{user}"),
				name: "users/not-a-uuid",
			},
			{
				f:    MustParseFormat("users/{user}"),
				name: "users/78da9161-aef1-49ed-bc92-0f136c95308f/purchases/1bba3dfe-5770-4d65-ae3f-1bff9e45b668",
			},
			{
				f:    MustParseFormat("users/{user}"),
				name: "stores/78da9161-aef1-49ed-bc92-0f136c95308f",
			},
			{
				f:    MustParseFormat("stores/{store}/products/{product}"),
				name: "stores/78da9161-aef1-49ed-bc92-0f136c95308f/products",
			},
			{
				f:    MustParseFormat("stores/{store}/products/{product}"),
				name: "stores/78da9161-aef1-49ed-bc92-0f136c95308f/products/",
			},
			{
				f:    MustParseFormat("stores/{store}/products/{product}"),
				name: "stores/78da9161-aef1-49ed-bc92-0f136c95308f/products/not-a-uuid",
			},
			{
				f:    MustParseFormat("stores/{store}/products/{product}"),
				name: "stores/not-a-uuid/products/1bba3dfe-5770-4d65-ae3f-1bff9e45b668",
			},
			{
				f:    MustParseFormat("stores/{store}/products/{product}"),
				name: "stores/78da9161-aef1-49ed-bc92-0f136c95308f",
			},
			{
				f:    MustParseFormat("stores/{store}/products/{product}"),
				name: "products/78da9161-aef1-49ed-bc92-0f136c95308f",
			},
		} {
			if _, err := test.f.Parse(test.name); err == nil {
				t.Errorf("f.Parse(%q) err = nil; want non-nil", test.name)
			}
		}
	})
}
