package users

import (
	"fmt"

	"github.com/google/uuid"
)

const CollectionID = "users"

func GenerateName() string {
	return fmt.Sprintf("%s/%s", CollectionID, uuid.New().String())
}
