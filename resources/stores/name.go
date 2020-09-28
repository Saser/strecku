package stores

import (
	"fmt"

	"github.com/google/uuid"
)

const (
	CollectionID = "stores"

	prefix = CollectionID + "/"
)

func GenerateName() string {
	return fmt.Sprintf("%s/%s", CollectionID, uuid.New().String())
}
