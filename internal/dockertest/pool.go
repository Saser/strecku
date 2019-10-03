package dockertest

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"go.uber.org/zap"
)

type Pool struct {
	logger *zap.Logger
	cli    *Client
}

func NewPool(
	logger *zap.Logger,
	cli *Client,
) *Pool {
	logger.Info("creating dockertest pool")
	pool := &Pool{
		logger: logger,
		cli:    cli,
	}
	logger.Info("created dockertest pool")
	return pool
}

func (p *Pool) PullOfficialImage(ctx context.Context, imageName string, tag string) error {
	canonicalName := fmt.Sprintf("docker.io/library/%s:%s", imageName, tag)
	logger := p.logger.With(zap.String("canonicalName", canonicalName))
	logger.Info("pulling image")
	reader, err := p.cli.ImagePull(ctx, canonicalName, types.ImagePullOptions{})
	if err != nil {
		return fmt.Errorf("pull official image: %w", err)
	}
	_, err = io.Copy(ioutil.Discard, reader)
	if err != nil {
		return fmt.Errorf("pull official image: %w", err)
	}
	logger.Info("pulled image")
	return nil
}

func (p *Pool) StartContainer(ctx context.Context, imageName string, tag string) (string, error) {
	imageLogger := p.logger.With(
		zap.String("imageName", imageName),
		zap.String("tag", tag),
	)
	imageLogger.Info(
		"pulling image of container to run",
		zap.String("imageName", imageName),
		zap.String("tag", tag),
	)
	if err := p.PullOfficialImage(ctx, imageName, tag); err != nil {
		return "", fmt.Errorf("start container: %w", err)
	}
	imageLogger.Info("image pulled, creating container")
	containerConfig := &container.Config{
		Image: fmt.Sprintf("%v:%v", imageName, tag),
	}
	hostConfig := &container.HostConfig{
		PublishAllPorts: true,
	}
	createResponse, err := p.cli.ContainerCreate(ctx, containerConfig, hostConfig, nil, "")
	if err != nil {
		return "", fmt.Errorf("start container: %w", err)
	}
	id := createResponse.ID
	idLogger := p.logger.With(zap.String("id", id))
	idLogger.Info(
		"created container, now starting it",
	)
	if err := p.cli.ContainerStart(ctx, id, types.ContainerStartOptions{}); err != nil {
		return "", fmt.Errorf("start container: %w", err)
	}
	idLogger.Info("container started")
	return id, nil
}

func (p *Pool) StopContainer(ctx context.Context, containerID string, timeout time.Duration) error {
	idLogger := p.logger.With(zap.String("id", containerID))
	idLogger.Info("stopping container")
	if err := p.cli.ContainerStop(ctx, containerID, &timeout); err != nil {
		return fmt.Errorf("stop container: %w", err)
	}
	idLogger.Info("container stopped, now removing it")
	if err := p.cli.ContainerRemove(ctx, containerID, types.ContainerRemoveOptions{}); err != nil {
		return fmt.Errorf("stop container: %w", err)
	}
	idLogger.Info("container removed")
	return nil
}

func (p *Pool) ContainerExists(ctx context.Context, containerID string) (bool, error) {
	idLogger := p.logger.With(zap.String("id", containerID))
	idLogger.Info("checking if container exists by inspecting it")
	_, err := p.cli.ContainerInspect(ctx, containerID)
	if err != nil {
		if client.IsErrContainerNotFound(err) {
			return false, nil
		}
		return false, fmt.Errorf("container exists: %w", err)
	}
	return true, nil
}

func (p *Pool) WithContainer(
	ctx context.Context,
	imageName string,
	tag string,
	stopTimeout time.Duration,
	f func(string) error,
) (err error) {
	p.logger.Info(
		"running function with throwaway container",
		zap.String("imageName", imageName),
		zap.String("tag", tag),
	)
	id, startErr := p.StartContainer(ctx, imageName, tag)
	if startErr != nil {
		return fmt.Errorf("with container: %w", startErr)
	}
	defer func() {
		if stopErr := p.StopContainer(ctx, id, stopTimeout); stopErr != nil {
			err = stopErr
		}
	}()
	if err := f(id); err != nil {
		return fmt.Errorf("with container: %w", err)
	}
	return nil
}

func (p *Pool) GetPortBinding(ctx context.Context, containerID string, port string) (string, error) {
	idLogger := p.logger.With(
		zap.String("id", containerID),
	)
	idLogger.Info("inspecting container")
	inspectResponse, err := p.cli.ContainerInspect(ctx, containerID)
	if err != nil {
		return "", fmt.Errorf("get bound port: %w", err)
	}
	idLogger.Info("inspected container")
	portMappings := inspectResponse.NetworkSettings.Ports[nat.Port(port)]
	if len(portMappings) == 0 {
		return "", fmt.Errorf("get bound port: no mappings for port %v", port)
	}
	hostPort := portMappings[0].HostPort
	portBinding := fmt.Sprintf("%v:%v", p.cli.Host, hostPort)
	idLogger.Info(
		"found port mapping",
		zap.String("port", port),
		zap.String("portBinding", portBinding),
	)
	return portBinding, nil
}
