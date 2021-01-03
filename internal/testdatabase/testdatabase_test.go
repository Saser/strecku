package testdatabase

import (
	"context"
	"testing"

	"github.com/cenkalti/backoff/v4"
	"golang.org/x/sync/errgroup"
)

func TestTestDatabase_Serve(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	tdb := New("../../database")
	g, ctx := errgroup.WithContext(ctx)
	g.Go(func() error {
		return tdb.Serve(ctx)
	})
	t.Cleanup(func() {
		cancel()
		if err := g.Wait(); err != nil {
			t.Errorf("g.Wait() = %v; want nil", err)
		}
	})

	op := func() error {
		_, err := tdb.ConnString()
		return err
	}
	if err := backoff.Retry(op, backoff.NewExponentialBackOff()); err != nil {
		t.Errorf("tdb.ConnString() err = %v; want nil", err)
	}
}
