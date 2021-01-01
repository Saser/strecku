package name

import (
	"fmt"
	"regexp"
)

var (
	collectionRE = regexp.MustCompile("^[a-z]+$")
)

type InvalidCollectionID string

func (e InvalidCollectionID) Error() string {
	return fmt.Sprintf("name: invalid collection ID: %s", string(e))
}

type CollectionID string

func (c CollectionID) Validate() error {
	if !collectionRE.MatchString(string(c)) {
		return InvalidCollectionID(c)
	}
	return nil
}
