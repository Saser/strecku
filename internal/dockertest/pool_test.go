package dockertest

import (
	"context"
	"fmt"
	"testing"

	"github.com/Saser/strecku/internal/provide"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewPool(t *testing.T) {
	logger := provide.ZapTestLogger(t)
	cli, cleanup, err := NewClient(logger)
	require.NoError(t, err)
	defer cleanup()
	_ = NewPool(logger, cli)
}

func TestPool_PullOfficialImage(t *testing.T) {
	logger := provide.ZapTestLogger(t)
	cli, cleanup, err := NewClient(logger)
	require.NoError(t, err)
	defer cleanup()
	pool := NewPool(logger, cli)
	ctx := context.Background()
	for _, tt := range []struct {
		image string
		tag   string
		valid bool
	}{
		{image: "postgres", tag: "11.5-alpine", valid: true},
		{image: "invalid", tag: "invalid", valid: false},
	} {
		tt := tt
		t.Run(fmt.Sprintf("image=%v,tag=%v", tt.image, tt.tag), func(t *testing.T) {
			assertFunction := assert.NoError
			if !tt.valid {
				assertFunction = assert.Error
			}
			assertFunction(t, pool.PullOfficialImage(ctx, tt.image, tt.tag))
		})
	}
}
