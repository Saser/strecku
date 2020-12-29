package testdatabase

import (
	"context"
	"testing"
)

func TestNew(t *testing.T) {
	ctx := context.Background()
	pool := New(ctx, t, "../../database")
	row := pool.QueryRow(ctx, "SELECT 1;")
	var one int
	if err := row.Scan(&one); err != nil {
		t.Fatal(err)
	}
}
