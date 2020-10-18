package auth

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestBearer_GetRequestMetadata(t *testing.T) {
	ctx := context.Background()
	for _, test := range []struct {
		b        Bearer
		wantMD   map[string]string
		wantCode codes.Code
	}{
		{
			b:      Bearer{Token: "foobar"},
			wantMD: map[string]string{"authorization": "Bearer foobar"},
		},
		{
			b:        Bearer{Token: ""},
			wantCode: codes.Unauthenticated,
		},
	} {
		gotMD, gotErr := test.b.GetRequestMetadata(ctx)
		if diff := cmp.Diff(gotMD, test.wantMD); diff != "" {
			t.Errorf("metadata mismatch (-got +want):\n%s", diff)
		}
		if got := status.Code(gotErr); got != test.wantCode {
			t.Errorf("status.Code(%v) = %v; want %v", gotErr, got, test.wantCode)
		}
	}
}

func TestBearer_GetTransportSecurity(t *testing.T) {
	var b Bearer
	if got := b.RequireTransportSecurity(); got != true {
		t.Errorf("b.RequireTransportSecurity() = %v; want true", got)
	}
}
