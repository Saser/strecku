package user

import (
	"go.uber.org/zap"
)

type Server struct {
	logger *zap.Logger
}

func NewServer(
	logger *zap.Logger,
) *Server {
	logger.Info("creating new User API server")
	server := &Server{
		logger: logger,
	}
	logger.Info("created new User API server")
	return server
}
