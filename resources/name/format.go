package name

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/google/uuid"
)

var (
	collectionRegexp = regexp.MustCompile("^[a-z]+$")
	variableRegexp   = regexp.MustCompile(`^\{[a-z]+\}$`)
	uuidRegexp       = func() *regexp.Regexp {
		const hex = "[0-9a-f]"
		counts := []int{8, 4, 4, 4, 12}
		parts := make([]string, len(counts))
		for i, n := range counts {
			parts[i] = hex + fmt.Sprintf("{%d}", n)
		}
		return regexp.MustCompile("^" + strings.Join(parts, "-") + "$")
	}()
)

type matcher interface {
	VarName() string
	Match(string) bool
	String() string
}

type exactMatcher string

func (m exactMatcher) VarName() string {
	return ""
}

func (m exactMatcher) Match(s string) bool {
	return s == string(m)
}

func (m exactMatcher) String() string {
	return string(m)
}

type uuidMatcher string

func (m uuidMatcher) VarName() string {
	return string(m)
}

func (m uuidMatcher) Match(s string) bool {
	return uuidRegexp.MatchString(s)
}

func (m uuidMatcher) String() string {
	return "{" + m.VarName() + "}"
}

type Format struct {
	matchers []matcher
}

func ParseFormat(s string) (*Format, error) {
	segments := strings.Split(s, "/")
	matchers := make([]matcher, len(segments))
	seen := make(map[string]bool)
	for i, s := range segments {
		switch {
		case collectionRegexp.MatchString(s):
			matchers[i] = exactMatcher(s)
		case variableRegexp.MatchString(s):
			varName := strings.Trim(s, "{}")
			if seen[varName] {
				return nil, fmt.Errorf("invalid format: variable %q occurs multiple times", varName)
			}
			seen[varName] = true
			matchers[i] = uuidMatcher(varName)
		default:
			return nil, fmt.Errorf("invalid format: %q is not a valid format segment", s)
		}
	}
	return &Format{
		matchers: matchers,
	}, nil
}

func MustParseFormat(s string) *Format {
	f, err := ParseFormat(s)
	if err != nil {
		panic(err)
	}
	return f
}

func (f *Format) Append(s string) (*Format, error) {
	return ParseFormat(f.String() + s)
}

func (f *Format) MustAppend(s string) *Format {
	return MustParseFormat(f.String() + s)
}

func (f *Format) String() string {
	segments := make([]string, len(f.matchers))
	for i, m := range f.matchers {
		segments[i] = m.String()
	}
	return strings.Join(segments, "/")
}

type UUIDs map[string]uuid.UUID

func (f *Format) Parse(name string) (UUIDs, error) {
	segments := strings.Split(name, "/")
	if got, want := len(segments), len(f.matchers); got != want {
		return nil, fmt.Errorf("invalid name: got %d segments, want %d", got, want)
	}
	uuids := make(map[string]uuid.UUID)
	for i, s := range segments {
		m := f.matchers[i]
		if !m.Match(s) {
			return nil, fmt.Errorf("invalid name: got %q, want name in format %q", name, f.String())
		}
		if varName := m.VarName(); varName != "" {
			uuids[varName] = uuid.MustParse(s)
		}
	}
	return uuids, nil
}

func (f *Format) Format(uuids UUIDs) (string, error) {
	segments := make([]string, len(f.matchers))
	for i, m := range f.matchers {
		var s string
		if varName := m.VarName(); varName != "" {
			u, ok := uuids[varName]
			if !ok {
				return "", fmt.Errorf("invalid UUIDs: does not contain variable %q", varName)
			}
			s = u.String()
		} else {
			s = m.String()
		}
		segments[i] = s
	}
	return strings.Join(segments, "/"), nil
}
