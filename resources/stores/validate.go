package stores

import (
	"errors"
	"fmt"
	"strings"

	streckuv1 "github.com/Saser/strecku/saser/strecku/v1"
	"github.com/google/uuid"
)

var (
	ErrNameEmpty         = errors.New("name is empty")
	ErrNameInvalidFormat = fmt.Errorf("name must have format %q", prefix+"<uuid>")
	ErrDisplayNameEmpty  = errors.New("display name is empty")
)

func Validate(store *streckuv1.Store) error {
	if store.Name == "" {
		return ErrNameEmpty
	}
	if !strings.HasPrefix(store.Name, prefix) {
		return ErrNameInvalidFormat
	}
	if _, err := uuid.Parse(strings.TrimPrefix(store.Name, prefix)); err != nil {
		return ErrNameInvalidFormat
	}
	if store.DisplayName == "" {
		return ErrDisplayNameEmpty
	}
	return nil
}
