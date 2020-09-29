package memberships

import (
	"fmt"
)

type MembershipNotFoundError struct {
	Name        string
	User, Store string
}

func (e *MembershipNotFoundError) Error() string {
	var query string
	switch {
	case e.Name != "":
		query = fmt.Sprintf("%q", e.Name)
	case e.User != "" && e.Store != "":
		query = fmt.Sprintf("between %q and %q", e.User, e.Store)
	}
	return fmt.Sprintf("membership not found: %s", query)
}

func (e *MembershipNotFoundError) Is(target error) bool {
	other, ok := target.(*MembershipNotFoundError)
	if !ok {
		return false
	}
	return e.Name == other.Name && e.User == other.User && e.Store == other.Store
}

type MembershipExistsError struct {
	Name        string
	User, Store string
}

func (e *MembershipExistsError) Error() string {
	var query string
	switch {
	case e.Name != "":
		query = fmt.Sprintf("%q", e.Name)
	case e.User != "" && e.Store != "":
		query = fmt.Sprintf("between %q and %q", e.User, e.Store)
	}
	return fmt.Sprintf("membership exists: %s", query)
}

func (e *MembershipExistsError) Is(target error) bool {
	other, ok := target.(*MembershipExistsError)
	if !ok {
		return false
	}
	return e.Name == other.Name && e.User == other.User && e.Store == other.Store
}
