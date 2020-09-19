package users

import (
	"fmt"

	"github.com/google/uuid"
)

const (
	CollectionID = "users"

	prefix = CollectionID + "/"
)

func GenerateName() string {
	return fmt.Sprintf("%s/%s", CollectionID, uuid.New().String())
}
