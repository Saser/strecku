package stores

import (
	"errors"
	"fmt"

	pb "github.com/Saser/strecku/api/v1"
)

var (
	ErrNameEmpty         = errors.New("name is empty")
	ErrNameInvalidFormat = fmt.Errorf("name must have format %q", prefix+"<uuid>")
	ErrDisplayNameEmpty  = errors.New("display name is empty")
)

func Validate(store *pb.Store) error {
	if err := ValidateName(store.Name); err != nil {
		return err
	}
	if store.DisplayName == "" {
		return ErrDisplayNameEmpty
	}
	return nil
}
