package auth

import (
	"context"
	"encoding/base64"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"
)

// Basic implements the credentials.PerRPCCredentials interface to provide
// authentication similar to HTTP Basic.
type Basic struct {
	Username string
	Password string
}

// Compile-time assertion that Basic implements the
// credentials.PerRPCCredentials interface.
var _ credentials.PerRPCCredentials = (*Basic)(nil)

// GetRequestMetadata returns metadata to attach to each request. The metadata
// contains one key-value pair: the key is "authorization" and the value is
// "Basic <enc>", where <enc> is the base64 encoding of the string formed by
// b.Username + ":" + b.Password.
func (b Basic) GetRequestMetadata(context.Context, ...string) (map[string]string, error) {
	if b.Username == "" || b.Password == "" {
		return nil, status.Error(codes.Unauthenticated, "Username and password are required.")
	}
	auth := fmt.Sprintf("%s:%s", b.Username, b.Password)
	enc := base64.StdEncoding.EncodeToString([]byte(auth))
	return map[string]string{
		"authorization": fmt.Sprintf("Basic %s", enc),
	}, nil
}

// RequireTransportSecurity returns true.
func (b Basic) RequireTransportSecurity() bool { return true }
