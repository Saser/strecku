package dockertest

import (
	"fmt"
	"net/url"
	"os"

	dockerclient "github.com/docker/docker/client"
	"go.uber.org/zap"
)

type Client struct {
	*dockerclient.Client
	Host string
}

func NewClient(
	logger *zap.Logger,
) (*Client, func(), error) {
	logger.Info(
		"creating docker client from environment",
		zap.String("DOCKER_HOST", os.Getenv("DOCKER_HOST")),
		zap.String("DOCKER_API_VERSION", os.Getenv("DOCKER_API_VERSION")),
		zap.String("DOCKER_CERT_PATH", os.Getenv("DOCKER_CERT_PATH")),
		zap.String("DOCKER_TLS_VERIFY", os.Getenv("DOCKER_TLS_VERIFY")),
	)
	dockerClient, err := dockerclient.NewEnvClient()
	if err != nil {
		return nil, nil, fmt.Errorf("create docker client: %w", err)
	}
	logger.Info("created docker client from environment")
	host := os.Getenv("DOCKER_HOST")
	if host != "" {
		hostURL, err := url.Parse(host)
		if err != nil {
			return nil, nil, fmt.Errorf("create docker client: %w", err)
		}
		host = hostURL.Host
	}
	logger.Info("creating dockertest client wrapper", zap.String("host", host))
	client := &Client{
		Client: dockerClient,
		Host:   host,
	}
	logger.Info("created dockertest client wrapper")
	cleanup := func() {
		logger.Info("closing docker client")
		if err := dockerClient.Close(); err != nil {
			logger.Error("closing docker client failed", zap.Error(err))
			return
		}
		logger.Info("closed docker client")
	}
	return client, cleanup, nil
}
