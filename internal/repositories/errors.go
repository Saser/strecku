package repositories

import "fmt"

type NotFound struct {
	Name string
}

func (e *NotFound) Error() string {
	return fmt.Sprintf("not found: %q", e.Name)
}

func (e *NotFound) Is(target error) bool {
	other, ok := target.(*NotFound)
	return ok && e.Name == other.Name
}

type Exists struct {
	Name string
}

func (e *Exists) Error() string {
	return fmt.Sprintf("exists: %q", e.Name)
}

func (e *Exists) Is(target error) bool {
	other, ok := target.(*Exists)
	return ok && e.Name == other.Name
}
