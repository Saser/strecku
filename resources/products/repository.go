package products

import "fmt"

type ProductNotFoundError struct {
	Name string
}

func (e *ProductNotFoundError) Error() string {
	return fmt.Sprintf("product not found: %q", e.Name)
}

func (e *ProductNotFoundError) Is(target error) bool {
	other, ok := target.(*ProductNotFoundError)
	if !ok {
		return false
	}
	return other.Name == e.Name
}

type ProductExistsError struct {
	Name string
}

func (e *ProductExistsError) Error() string {
	return fmt.Sprintf("product exists: %q", e.Name)
}

func (e *ProductExistsError) Is(target error) bool {
	other, ok := target.(*ProductExistsError)
	if !ok {
		return false
	}
	return other.Name == e.Name
}
