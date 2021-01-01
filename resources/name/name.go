package name

import (
	"fmt"
	"strings"
)

type Name string

func (n Name) Validate() error {
	s := n.segments()
	if n := len(s); n%2 != 0 {
		return fmt.Errorf("name: odd number of segments (%d)", n)
	}
	for i := 0; i < len(s); i += 2 {
		c, r := s[i], s[i+1]
		if err := CollectionID(c).Validate(); err != nil {
			return err
		}
		if err := ResourceID(r).Validate(); err != nil {
			return err
		}
	}
	return nil
}

func (n Name) segments() []string {
	return strings.Split(string(n), "/")
}
