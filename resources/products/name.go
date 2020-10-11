package products

import (
	"errors"
	"fmt"

	"github.com/Saser/strecku/resources/names"
	"github.com/Saser/strecku/resources/stores"
	"github.com/google/uuid"
)

const CollectionID = "products"

var (
	Regexp = names.MustCompile(
		fmt.Sprintf("(?P<store>%s)", stores.Regexp.String()),
		CollectionID, names.UUID,
	)

	ErrNameEmpty         = errors.New("name is empty")
	ErrNameInvalidFormat = fmt.Errorf("name must have format %q", stores.CollectionID+"/<uuid>/"+CollectionID+"/<uuid>")
)

func GenerateName(store string) string {
	return fmt.Sprintf("%s/%s/%s", store, CollectionID, uuid.New().String())
}

func ValidateName(name string) error {
	if name == "" {
		return ErrNameEmpty
	}
	if !Regexp.MatchString(name) {
		return ErrNameInvalidFormat
	}
	return nil
}

func Parent(name string) (string, error) {
	if err := ValidateName(name); err != nil {
		return "", err
	}
	matches := Regexp.FindStringSubmatch(name)
	return matches[Regexp.SubexpIndex("store")], nil
}
