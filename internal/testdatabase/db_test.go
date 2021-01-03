package testdatabase

import (
	"context"
	"testing"
)

func TestDB(t *testing.T) {
	Init()
	t.Cleanup(Cleanup)
	_ = DB(context.Background(), t, "../../database")
}
