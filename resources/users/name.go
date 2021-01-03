package users

import (
	"github.com/Saser/strecku/resourcename"
	"github.com/google/uuid"
)

const (
	CollectionID = "users"
)

var NameFormat = resourcename.MustParseFormat(CollectionID + "/{user}")

func GenerateName() string {
	name, err := NameFormat.Format(resourcename.UUIDs{"user": uuid.New()})
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
	return uuids["user"], nil
}

func ValidateName(name string) error {
	_, err := ParseName(name)
	return err
}
