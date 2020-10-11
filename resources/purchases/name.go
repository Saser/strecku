package purchases

import (
	"errors"
	"fmt"

	"github.com/Saser/strecku/resources/names"
	"github.com/google/uuid"
)

const CollectionID = "purchases"

var (
	Regexp = names.MustCompile(CollectionID, names.UUID)

	ErrNameEmpty         = errors.New("name is empty")
	ErrNameInvalidFormat = fmt.Errorf("name must have format %q", CollectionID+"/<uuid>")
)

func GenerateName() string {
	return fmt.Sprintf("%s/%s", CollectionID, uuid.New().String())
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