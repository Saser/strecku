package auth

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"
)

// Bearer implements the credentials.PerRPCCredentials interface to provide
// authentication similar to Bearer token authorization as commonly used for
// REST interfaces.
type Bearer struct {
	Token string
}

// Compile-time assertion that Bearer implements the
// credentials.PerRPCCredentials interface.
var _ credentials.PerRPCCredentials = (*Bearer)(nil)

// GetRequestMetadata returns metadata to attach to each request. The metadata
// contains one key-value pair: the key is "authorization" and the value is
// "Bearer <token>" where <token> is b.Token.
//
// If b.Token is empty, GetRequestMetadata returns an error with code
// codes.Unauthenticated.
func (b Bearer) GetRequestMetadata(context.Context, ...string) (map[string]string, error) {
	if b.Token == "" {
		return nil, status.Error(codes.Unauthenticated, "Token is required.")
	}
	return map[string]string{
		"authorization": "Bearer " + b.Token,
	}, nil
}

// RequireTransportSecurity returns true.
func (b Bearer) RequireTransportSecurity() bool { return true }
