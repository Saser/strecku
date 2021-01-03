package testdatabase

import (
	"context"
	"testing"
)

func TestDB(t *testing.T) {
	_ = DB(context.Background(), t, "../../database")
}
