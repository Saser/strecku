package dockertest

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/Saser/strecku/internal/provide"
	"github.com/stretchr/testify/require"
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
				_, err := client.Ping(ctx)
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}
		})
	}
}
