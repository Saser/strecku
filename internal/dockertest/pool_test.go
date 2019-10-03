package dockertest

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/Saser/strecku/internal/provide"
	"github.com/cenkalti/backoff/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestNewPool(t *testing.T) {
	logger := provide.ZapTestLogger(t)
	_, cleanup := pool(t, logger)
	defer cleanup()
}

func TestPool_PullOfficialImage(t *testing.T) {
	logger := provide.ZapTestLogger(t)
	pool, cleanup := pool(t, logger)
	defer cleanup()
	ctx := context.Background()
	for _, tt := range []struct {
		image string
		tag   string
		valid bool
	}{
		{image: "hello-world", tag: "linux", valid: true},
		{image: "invalid", tag: "invalid", valid: false},
	} {
		tt := tt
		t.Run(fmt.Sprintf("image=%v,tag=%v", tt.image, tt.tag), func(t *testing.T) {
			err := pool.PullOfficialImage(ctx, tt.image, tt.tag)
			if tt.valid {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
		})
	}
}

func TestPool_StartContainer_StopContainer_ContainerExists(t *testing.T) {
	logger := provide.ZapTestLogger(t)
	pool, cleanup := pool(t, logger)
	defer cleanup()
	ctx := context.Background()
	for _, tt := range []struct {
		image string
		tag   string
		valid bool
	}{
		{image: "hello-world", tag: "linux", valid: true},
		{image: "invalid", tag: "invalid", valid: false},
	} {
		tt := tt
		t.Run(fmt.Sprintf("image=%v,tag=%v,valid=%v", tt.image, tt.tag, tt.valid), func(t *testing.T) {
			id, err := pool.StartContainer(ctx, tt.image, tt.tag)
			if tt.valid {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
				return
			}
			totalTimeout := 1 * time.Minute
			stopTimeout := 10 * time.Second
			stopCtx, stopCancel := context.WithTimeout(ctx, totalTimeout)
			defer stopCancel()
			operation := func() error { return pool.StopContainer(stopCtx, id, stopTimeout) }
			policy := backoff.WithContext(backoff.NewExponentialBackOff(), stopCtx)
			err = backoff.Retry(operation, policy)
			require.NoError(t, err)
			exists, err := pool.ContainerExists(ctx, id)
			require.NoError(t, err)
			require.False(t, exists)
		})
	}
}

func pool(t *testing.T, logger *zap.Logger) (*Pool, func()) {
	cli, cleanup, err := NewClient(logger)
	require.NoError(t, err)
	pool := NewPool(logger, cli)
	return pool, cleanup
}
