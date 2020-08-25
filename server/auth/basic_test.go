package auth

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestBasic_GetRequestMetadata(t *testing.T) {
	ctx := context.Background()
	for _, test := range []struct {
		b        Basic
		wantMD   map[string]string
		wantCode codes.Code
	}{
		{
			b:      Basic{Username: "user@example.com", Password: "password"},
			wantMD: map[string]string{"authorization": "Basic dXNlckBleGFtcGxlLmNvbTpwYXNzd29yZA=="}, // pre-calculated value
		},
		{
			b:        Basic{Username: "user@example.com", Password: ""},
			wantCode: codes.InvalidArgument,
		},
		{
			b:        Basic{Username: "", Password: "password"},
			wantCode: codes.InvalidArgument,
		},
		{
			b:        Basic{Username: "", Password: ""},
			wantCode: codes.InvalidArgument,
		},
	} {
		gotMD, gotErr := test.b.GetRequestMetadata(ctx)
		if diff := cmp.Diff(gotMD, test.wantMD); diff != "" {
			t.Errorf("metadata mismatch (-got +want):\n%s", diff)
		}
		if got := status.Code(gotErr); got != test.wantCode {
			t.Errorf("status.Code(err) = %v; want %v", got, test.wantCode)
		}
	}
}

func TestBasic_RequireTransportSecurity(t *testing.T) {
	var b Basic
	if got, want := b.RequireTransportSecurity(), true; got != want {
		t.Errorf("b.RequireTransportSecurity() = %v; want %v", got, want)
	}
}
