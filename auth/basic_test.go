package auth

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	exampleUsername = "user@example.com"
	examplePassword = "password"
	// exampleAuthorization is a pre-calculated value for exampleUsername and
	// examplePassword.
	exampleAuthorization = "Basic dXNlckBleGFtcGxlLmNvbTpwYXNzd29yZA=="
)

func TestParseBasic(t *testing.T) {
	for _, test := range []struct {
		s       string
		want    Basic
		wantErr bool
	}{
		{s: exampleAuthorization, want: Basic{Username: exampleUsername, Password: examplePassword}},
		{s: "Basic dXNlckBleGFtcGxlLmNvbTpwYXNzd29yZDp3aXRoOmNvbG9ucw==", want: Basic{Username: exampleUsername, Password: "password:with:colons"}},
		{s: "Bearer dXNlckBleGFtcGxlLmNvbTpwYXNzd29yZA==", wantErr: true},
		{s: "Basic ", wantErr: true},
		{s: "Basic", wantErr: true},
		{s: "", wantErr: true},
		{s: "Basic invalidbase64", wantErr: true},
		// "user@example.com:"
		{s: "Basic dXNlckBleGFtcGxlLmNvbTo=", wantErr: true},
		// ":password"
		{s: "Basic OnBhc3N3b3Jk=", wantErr: true},
		// ":"
		{s: "Basic Og==", wantErr: true},
	} {
		b, err := ParseBasic(test.s)
		if test.wantErr && err == nil {
			t.Errorf("ParseBasic(%q) did not return an error", test.s)
			continue
		}
		if diff := cmp.Diff(b, test.want); diff != "" {
			t.Errorf("-got +want:\n%s", diff)
		}
	}
}

func TestBasic_GetRequestMetadata(t *testing.T) {
	ctx := context.Background()
	for _, test := range []struct {
		b        Basic
		wantMD   map[string]string
		wantCode codes.Code
	}{
		{
			b:      Basic{Username: exampleUsername, Password: examplePassword},
			wantMD: map[string]string{"authorization": exampleAuthorization},
		},
		{
			b:        Basic{Username: exampleUsername, Password: ""},
			wantCode: codes.Unauthenticated,
		},
		{
			b:        Basic{Username: "", Password: examplePassword},
			wantCode: codes.Unauthenticated,
		},
		{
			b:        Basic{Username: "", Password: ""},
			wantCode: codes.Unauthenticated,
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
	if got := b.RequireTransportSecurity(); got != true {
		t.Errorf("b.RequireTransportSecurity() = %v; want true", got)
	}
}
