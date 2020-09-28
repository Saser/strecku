package memberships

import (
	"fmt"

	"github.com/google/uuid"
)

const (
	CollectionID = "memberships"

	prefix = CollectionID + "/"
)

func GenerateName() string {
	return fmt.Sprintf("%s/%s", CollectionID, uuid.New().String())
}
