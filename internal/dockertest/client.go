package dockertest

import (
	"bufio"
	"context"
	"fmt"
	"net"
	"net/url"
	"os"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	dockerclient "github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"go.uber.org/zap"
)

type Client struct {
	dc     *dockerclient.Client
	logger *zap.Logger
	host   string
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
		dc:     dockerClient,
		logger: logger,
		host:   host,
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

func (c *Client) PullOfficialImage(ctx context.Context, imageName string, tag string) error {
	canonicalName := fmt.Sprintf("docker.io/library/%s:%s", imageName, tag)
	logger := c.logger.With(zap.String("canonicalName", canonicalName))
	logger.Info("pulling image")
	reader, err := c.dc.ImagePull(ctx, canonicalName, types.ImagePullOptions{})
	if err != nil {
		return fmt.Errorf("pull official image: %w", err)
	}
	scanner := bufio.NewScanner(reader)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		msg := scanner.Text()
		logger.Debug(msg)
	}
	logger.Info("pulled image")
	if err := reader.Close(); err != nil {
		return fmt.Errorf("pull official image: %w", err)
	}
	return nil
}

// func (c *Client) StartContainer(ctx context.Context, imageName string, tag string, env map[string]string) (string, error) {
func (c *Client) StartContainer(ctx context.Context, imageName string, tag string) (string, error) {
	imageLogger := c.logger.With(
		zap.String("imageName", imageName),
		zap.String("tag", tag),
	)
	imageLogger.Info(
		"pulling image of container to run",
		zap.String("imageName", imageName),
		zap.String("tag", tag),
	)
	if err := c.PullOfficialImage(ctx, imageName, tag); err != nil {
		return "", fmt.Errorf("start container: %w", err)
	}
	imageLogger.Info("image pulled, creating container")
	// envStrings := make([]string, 0, len(env))
	// for key, value := range env {
	//	envStrings = append(envStrings, fmt.Sprintf("%s=%s", key, value))
	//}
	containerConfig := &container.Config{
		Image: fmt.Sprintf("%v:%v", imageName, tag),
		// Env:   envStrings,
	}
	hostConfig := &container.HostConfig{
		PublishAllPorts: true,
	}
	createResponse, err := c.dc.ContainerCreate(ctx, containerConfig, hostConfig, nil, "")
	if err != nil {
		return "", fmt.Errorf("start container: %w", err)
	}
	id := createResponse.ID
	idLogger := c.logger.With(zap.String("id", id))
	idLogger.Info(
		"created container, now starting it",
	)
	if err := c.dc.ContainerStart(ctx, id, types.ContainerStartOptions{}); err != nil {
		return "", fmt.Errorf("start container: %w", err)
	}
	idLogger.Info("container started")
	return id, nil
}

func (c *Client) StopContainer(ctx context.Context, containerID string, timeout time.Duration) error {
	idLogger := c.logger.With(zap.String("id", containerID))
	idLogger.Info("stopping container")
	if err := c.dc.ContainerStop(ctx, containerID, &timeout); err != nil {
		return fmt.Errorf("stop container: %w", err)
	}
	idLogger.Info("container stopped, now removing it and its volumes")
	opts := types.ContainerRemoveOptions{
		RemoveVolumes: true,
	}
	if err := c.dc.ContainerRemove(ctx, containerID, opts); err != nil {
		return fmt.Errorf("stop container: %w", err)
	}
	idLogger.Info("container removed")
	return nil
}

func (c *Client) ContainerExists(ctx context.Context, containerID string) (bool, error) {
	idLogger := c.logger.With(zap.String("id", containerID))
	idLogger.Info("checking if container exists by inspecting it")
	_, err := c.dc.ContainerInspect(ctx, containerID)
	if err != nil {
		if dockerclient.IsErrContainerNotFound(err) {
			return false, nil
		}
		return false, fmt.Errorf("container exists: %w", err)
	}
	return true, nil
}

func (c *Client) GetTCPAddress(ctx context.Context, containerID string, port string) (*net.TCPAddr, error) {
	idLogger := c.logger.With(
		zap.String("id", containerID),
	)
	idLogger.Info("inspecting container")
	inspectResponse, err := c.dc.ContainerInspect(ctx, containerID)
	if err != nil {
		return nil, fmt.Errorf("get TCP address: %w", err)
	}
	idLogger.Info("inspected container, looking up binding for port", zap.String("port", port))
	portBindings := inspectResponse.NetworkSettings.Ports[nat.Port(port)]
	if len(portBindings) == 0 {
		return nil, fmt.Errorf("get TCP address: no bindings for port %v", port)
	}
	hostPort := portBindings[0].HostPort
	idLogger.Info(
		"found binding for port, now resolving TCP address",
		zap.String("hostPort", hostPort),
		zap.String("host", c.host),
	)
	addressString := fmt.Sprintf("%v:%v", c.host, hostPort)
	address, err := net.ResolveTCPAddr("tcp", addressString)
	if err != nil {
		return nil, fmt.Errorf("get TCP address: %w", err)
	}
	idLogger.Info("resolved TCP address", zap.Any("address", address))
	return address, nil
}
