package stores

import (
	"fmt"

	"github.com/Saser/strecku/resources/names"
	"github.com/google/uuid"
)

const (
	CollectionID = "stores"

	prefix = CollectionID + "/"
)

var Regexp = names.MustCompile(CollectionID, names.UUID)

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
