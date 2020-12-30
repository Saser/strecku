package testdatabase

import (
	"strings"
	"testing"
)

func TestConnString(t *testing.T) {
	Init()
	t.Cleanup(Cleanup)
	c := ConnString()
	scheme := "postgres://"
	if !strings.HasPrefix(c, "postgres://") {
		t.Errorf("ConnString() = %q; want prefix %q", c, scheme)
	}
}
