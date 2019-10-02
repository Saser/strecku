package dockertest

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"

	"github.com/docker/docker/api/types"
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
