package name

import (
	"fmt"
	"regexp"
	"strings"
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
	for i, s := range segments {
		switch {
		case collectionRegexp.MatchString(s):
			matchers[i] = exactMatcher(s)
		case variableRegexp.MatchString(s):
			matchers[i] = uuidMatcher(strings.Trim(s, "{}"))
		default:
			return nil, fmt.Errorf("invalid format: %q is not a valid format segment", s)
		}
	}
	return &Format{
		matchers: matchers,
	}, nil
}

func (f *Format) String() string {
	segments := make([]string, len(f.matchers))
	for i, m := range f.matchers {
		segments[i] = m.String()
	}
	return strings.Join(segments, "/")
}
