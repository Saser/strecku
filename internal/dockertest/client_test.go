package dockertest

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/Saser/strecku/internal/provide"
	"github.com/cenkalti/backoff/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestNewClient(t *testing.T) {
	logger := provide.ZapTestLogger(t)
	ctx := context.Background()
	defaultHost := os.Getenv("DOCKER_HOST")
	for _, tt := range []struct {
		host  string
		valid bool
	}{
		{host: defaultHost, valid: true},
		{host: "invalid", valid: false},
	} {
		tt := tt
		t.Run(fmt.Sprintf("host=%v,valid=%v", tt.host, tt.valid), func(t *testing.T) {
			defer func() {
				require.NoError(t, os.Setenv("DOCKER_HOST", defaultHost))
			}()
			require.NoError(t, os.Setenv("DOCKER_HOST", tt.host))
			client, cleanup, err := NewClient(logger)
			if tt.valid {
				require.NoError(t, err)
				defer cleanup()
				_, err := client.dc.Ping(ctx)
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}
		})
	}
}

func TestClient_PullOfficialImage(t *testing.T) {
	ctx := context.Background()
	logger := provide.ZapTestLogger(t)
	client, cleanup := client(ctx, t, logger)
	defer cleanup()
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
			err := client.PullOfficialImage(ctx, tt.image, tt.tag)
			if tt.valid {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
		})
	}
}

func TestClient_StartContainer_StopContainer_ContainerExists(t *testing.T) {
	ctx := context.Background()
	logger := provide.ZapTestLogger(t)
	client, cleanup := client(ctx, t, logger)
	defer cleanup()
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
			id, err := client.StartContainer(ctx, tt.image, tt.tag)
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
			operation := func() error { return client.StopContainer(stopCtx, id, stopTimeout) }
			policy := backoff.WithContext(backoff.NewExponentialBackOff(), stopCtx)
			err = backoff.Retry(operation, policy)
			require.NoError(t, err)
			exists, err := client.ContainerExists(ctx, id)
			require.NoError(t, err)
			require.False(t, exists)
		})
	}
}

func TestClient_GetTCPAddress(t *testing.T) {
	ctx := context.Background()
	logger := provide.ZapTestLogger(t)
	client, cleanup := client(ctx, t, logger)
	defer cleanup()
	id, err := client.StartContainer(ctx, "postgres", "11.5-alpine")
	require.NoError(t, err)
	stopTimeout := 10 * time.Second
	defer func() {
		err := client.StopContainer(ctx, id, stopTimeout)
		require.NoError(t, err)
	}()
	_, err = client.GetTCPAddress(ctx, id, "5432/tcp")
	require.NoError(t, err)
}

func client(ctx context.Context, t *testing.T, logger *zap.Logger) (*Client, func()) {
	c, cleanup, err := NewClient(logger)
	require.NoError(t, err)
	_, err = c.dc.Ping(ctx)
	require.NoError(t, err)
	return c, cleanup
}
