package users

import (
	"strings"
	"testing"

	"github.com/google/uuid"
)

func TestGenerateName(t *testing.T) {
	got := GenerateName()
	if !strings.HasPrefix(got, prefix) {
		t.Errorf("GenerateName() = %q; want prefix %q", got, prefix)
	}
	id := strings.TrimPrefix(got, prefix)
	if _, err := uuid.Parse(id); err != nil {
		t.Errorf("uuid.Parse(%q) err = %v; want nil", id, err)
	}
}
