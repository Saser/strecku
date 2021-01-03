package resourcename

import (
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
)

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
		// Duplicate variables.
		"stores/{duplicate}/products/{duplicate}",
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
			_, err := ParseFormat(s)
			if err == nil {
				t.Errorf("ParseFormat(%q) err = nil; want non-nil", s)
			}
			var want *InvalidFormat
			if !errors.As(err, &want) {
				t.Errorf("errors.As(%v, %T) = false; want true", err, &want)
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

func TestFormat_Append(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		for _, s := range good {
			f := MustParseFormat(s)
			for _, s2 := range []string{
				"/bar",
				"/foos/{foo}",
				"/bars/{bar}/foos/{foo}",
			} {
				if _, err := f.Append(s2); err != nil {
					t.Errorf("f.Append(%q) err = %v; want nil", s2, err)
				}
			}
		}
	})

	t.Run("Errors", func(t *testing.T) {
		for _, s := range good {
			f := MustParseFormat(s)
			for _, s2 := range bad {
				s3 := "/" + s2
				if _, err := f.Append(s3); err == nil {
					t.Errorf("f.Append(%q) err = nil; want non-nil", s3)
				}
			}
		}
	})
}

func TestFormat_MustAppend(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		for _, s := range good {
			f := MustParseFormat(s)
			for _, s2 := range []string{
				"/bar",
				"/foos/{foo}",
				"/bars/{bar}/foos/{foo}",
			} {
				func() {
					defer func() {
						if err := recover(); err != nil {
							t.Errorf("f.MustAppend(%q) err = %v; want nil", s2, err)
						}
					}()
					_ = f.MustAppend(s2)
				}()
			}
		}
	})

	t.Run("Errors", func(t *testing.T) {
		for _, s := range good {
			f := MustParseFormat(s)
			for _, s2 := range bad {
				func() {
					s3 := "/" + s2
					defer func() {
						if err := recover(); err == nil {
							t.Errorf("f.MustAppend(%q) err = nil; want non-nil", s3)
						}
					}()
					_ = f.MustAppend(s3)
				}()
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

func TestFormat_Parse(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		for _, test := range []struct {
			f    *Format
			name string
			want UUIDs
		}{
			{
				f:    MustParseFormat("users/{user}"),
				name: "users/78da9161-aef1-49ed-bc92-0f136c95308f",
				want: UUIDs{
					"user": uuid.MustParse("78da9161-aef1-49ed-bc92-0f136c95308f"),
				},
			},
			{
				f:    MustParseFormat("stores/{store}/products/{product}"),
				name: "stores/78da9161-aef1-49ed-bc92-0f136c95308f/products/1bba3dfe-5770-4d65-ae3f-1bff9e45b668",
				want: UUIDs{
					"store":   uuid.MustParse("78da9161-aef1-49ed-bc92-0f136c95308f"),
					"product": uuid.MustParse("1bba3dfe-5770-4d65-ae3f-1bff9e45b668"),
				},
			},
		} {
			got, err := test.f.Parse(test.name)
			if err != nil {
				t.Errorf("f.Parse(%q) err = %v; want nil", test.name, err)
			}
			less := func(s1, s2 string) bool { return s1 < s2 }
			if diff := cmp.Diff(test.want, got, cmpopts.SortMaps(less)); diff != "" {
				t.Errorf("f.Parse(%q) UUIDs != test.want (-want +got)\n%s", test.name, diff)
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
			_, err := test.f.Parse(test.name)
			if err == nil {
				t.Errorf("f.Parse(%q) err = nil; want non-nil", test.name)
			}
			var want *InvalidName
			if !errors.As(err, &want) {
				t.Errorf("errors.As(%v, %T) = false; want true", err, &want)
			}
		}
	})
}

func TestFormat_Format(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		for _, test := range []struct {
			f     *Format
			uuids UUIDs
			want  string
		}{
			{
				f: MustParseFormat("users/{user}"),
				uuids: UUIDs{
					"user": uuid.MustParse("78da9161-aef1-49ed-bc92-0f136c95308f"),
				},
				want: "users/78da9161-aef1-49ed-bc92-0f136c95308f",
			},
			{
				f: MustParseFormat("stores/{store}/settings"),
				uuids: UUIDs{
					"store": uuid.MustParse("78da9161-aef1-49ed-bc92-0f136c95308f"),
				},
				want: "stores/78da9161-aef1-49ed-bc92-0f136c95308f/settings",
			},
			{
				f: MustParseFormat("stores/{store}/products/{product}"),
				uuids: UUIDs{
					"store":   uuid.MustParse("78da9161-aef1-49ed-bc92-0f136c95308f"),
					"product": uuid.MustParse("1bba3dfe-5770-4d65-ae3f-1bff9e45b668"),
				},
				want: "stores/78da9161-aef1-49ed-bc92-0f136c95308f/products/1bba3dfe-5770-4d65-ae3f-1bff9e45b668",
			},
		} {
			got, err := test.f.Format(test.uuids)
			if err != nil {
				t.Errorf("f.Format(%v) err = %v; want nil", test.uuids, err)
			}
			if got != test.want {
				t.Errorf("f.Format(%v) name = %q; want %q", test.uuids, got, test.want)
			}
		}
	})

	t.Run("Errors", func(t *testing.T) {
		for _, test := range []struct {
			f     *Format
			uuids UUIDs
		}{
			{
				f: MustParseFormat("users/{user}"),
				uuids: UUIDs{
					"store": uuid.MustParse("78da9161-aef1-49ed-bc92-0f136c95308f"),
				},
			},
			{
				f: MustParseFormat("stores/{store}/products/{product}"),
				uuids: UUIDs{
					"store": uuid.MustParse("78da9161-aef1-49ed-bc92-0f136c95308f"),
				},
			},
			{
				f: MustParseFormat("stores/{store}/products/{product}"),
				uuids: UUIDs{
					"product": uuid.MustParse("78da9161-aef1-49ed-bc92-0f136c95308f"),
				},
			},
		} {
			if _, err := test.f.Format(test.uuids); err == nil {
				t.Errorf("f.Format(%v) err = nil; want non-nil", test.uuids)
			}
		}
	})
}
