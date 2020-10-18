package stores

import (
	"fmt"

	"github.com/Saser/strecku/resources/names"
	"github.com/google/uuid"
)

const CollectionID = "stores"

var (
	Regexp = names.MustCompile(CollectionID, names.UUID)

	ErrNameInvalidFormat = fmt.Errorf("name must have format %q", CollectionID+"/<uuid>")
)

func GenerateName() string {
	return fmt.Sprintf("%s/%s", CollectionID, uuid.New().String())
}

func ValidateName(name string) error {
	if !Regexp.MatchString(name) {
		return ErrNameInvalidFormat
	}
	return nil
}
