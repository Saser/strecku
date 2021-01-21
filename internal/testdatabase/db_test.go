package testdatabase

import (
	"context"
	"testing"
)

func TestDB(t *testing.T) {
	if testing.Short() {
		t.Skipf("skipping: -short is set")
	}
	_ = DB(context.Background(), t, "../../database")
}
