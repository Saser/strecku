package stores

import (
	"errors"

	pb "github.com/Saser/strecku/api/v1"
)

var ErrDisplayNameEmpty = errors.New("display name is empty")

func Validate(store *pb.Store) error {
	if err := ValidateName(store.Name); err != nil {
		return err
	}
	if store.DisplayName == "" {
		return ErrDisplayNameEmpty
	}
	return nil
}
