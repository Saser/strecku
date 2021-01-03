package stores

import (
	"github.com/Saser/strecku/resourcename"
	"github.com/google/uuid"
)

const CollectionID = "stores"

var NameFormat = resourcename.MustParseFormat(CollectionID + "/{store}")

func GenerateName() string {
	name, err := NameFormat.Format(resourcename.UUIDs{"store": uuid.New()})
	if err != nil {
		panic(err)
	}
	return name
}

func ParseName(name string) (uuid.UUID, error) {
	uuids, err := NameFormat.Parse(name)
	if err != nil {
		return uuid.UUID{}, err
	}
	return uuids["store"], nil
}

func ValidateName(name string) error {
	_, err := ParseName(name)
	return err
}
