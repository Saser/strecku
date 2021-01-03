package memberships

import (
	"github.com/Saser/strecku/resources/stores"
	"github.com/google/uuid"
)

const CollectionID = "memberships"

var NameFormat = stores.NameFormat.MustAppend("/" + CollectionID + "/{membership}")

func GenerateName(store string) string {
	uuids, err := stores.NameFormat.Parse(store)
	if err != nil {
		panic(err)
	}
	uuids["membership"] = uuid.New()
	name, err := NameFormat.Format(uuids)
	if err != nil {
		panic(err)
	}
	return name
}

func ParseName(name string) (store uuid.UUID, membership uuid.UUID, err error) {
	uuids, err := NameFormat.Parse(name)
	if err != nil {
		return uuid.UUID{}, uuid.UUID{}, err
	}
	return uuids["store"], uuids["membership"], nil
}

func ValidateName(name string) error {
	_, _, err := ParseName(name)
	return err
}

func Parent(name string) (string, error) {
	uuids, err := NameFormat.Parse(name)
	if err != nil {
		return "", err
	}
	return stores.NameFormat.Format(uuids)
}
