package users

import (
	"strings"
	"testing"
)

func TestGenerateName(t *testing.T) {
	if got, want := GenerateName(), CollectionID+"/"; !strings.HasPrefix(got, want) {
		t.Errorf("GenerateName() = %q; want prefix %q", got, want)
	}
}
