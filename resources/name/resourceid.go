package name

import (
	"fmt"
	"regexp"
	"strings"
)

var (
	resourceRE *regexp.Regexp
)

func init() {
	const hex = "[0-9a-f]"
	counts := []int{8, 4, 4, 4, 12}
	parts := make([]string, len(counts))
	for i, n := range counts {
		parts[i] = hex + fmt.Sprintf("{%d}", n)
	}
	resourceRE = regexp.MustCompile("^" + strings.Join(parts, "-") + "$")
}

type InvalidResourceID string

func (e InvalidResourceID) Error() string {
	return fmt.Sprintf("name: invalid resource ID: %s", string(e))
}

type ResourceID string

func (r ResourceID) Validate() error {
	if !resourceRE.MatchString(string(r)) {
		return InvalidResourceID(r)
	}
	return nil
}
