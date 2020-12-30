package repositories

import (
	"context"
	"errors"
	"fmt"
	"testing"

	pb "github.com/Saser/strecku/api/v1"
)

var (
	ErrUnauthenticated = errors.New("unauthenticated")
)

type EmailAddressNotFound struct {
	EmailAddress string
}

func (e *EmailAddressNotFound) Error() string {
	return fmt.Sprintf("email address not found: %q", e.EmailAddress)
}

func (e *EmailAddressNotFound) Is(target error) bool {
	other, ok := target.(*EmailAddressNotFound)
	return ok && e.EmailAddress == other.EmailAddress
}

type EmailAddressExists struct {
	EmailAddress string
}

func (e *EmailAddressExists) Error() string {
	return fmt.Sprintf("email address exists: %q", e.EmailAddress)
}

func (e *EmailAddressExists) Is(target error) bool {
	other, ok := target.(*EmailAddressExists)
	return ok && e.EmailAddress == other.EmailAddress
}

type Users interface {
	// Authenticate determines whether there exists a user with the given
	// name and password. The name and password will be validated using
	// package users. If no combination of the given name and password
	// exists, ErrUnauthenticated will be returned.
	Authenticate(ctx context.Context, name string, password string) error

	// Lookup returns the user corresponding to the given name, or returns a
	// non-nil error otherwise. The name will be validated using package
	// users. If no user is found, a NotFound error will be returned.
	Lookup(ctx context.Context, name string) (*pb.User, error)

	// ResolveEmail returns the resource name for the user with the given
	// email address, or an EmailAddressNotFound error otherwise.
	ResolveEmail(ctx context.Context, emailAddress string) (string, error)

	// List returns a list of all users.
	List(ctx context.Context) ([]*pb.User, error)

	// Create creates a new user resource based on the given user, and
	// associates it with the given password. The given user and the
	// password will be validated using package users. If a user already
	// exists with the email address, an EmailAddressExists error will be
	// returned. If a user already exists with the given name, an Exists
	// error will be returned.
	Create(ctx context.Context, user *pb.User, password string) error

	// Update updates an existing user to the version specified by the given
	// user. The given user will be validated using package user. The name
	// of the given user is used to identify which user to update. If no
	// user with that name exists, a NotFound error will be returned.
	Update(ctx context.Context, user *pb.User) error

	// Delete deletes the user corresponding to the given name. The name
	// will be validated using package users. If no user with that name
	// exists, a NotFound error will be returned.
	Delete(ctx context.Context, name string) error
}

func SeedUsers(ctx context.Context, t *testing.T, r Users, users []*pb.User, passwords []string) {
	t.Helper()
	t.Cleanup(func() {
		all, err := r.List(ctx)
		if err != nil {
			t.Error(err)
		}
		for _, user := range all {
			if err := r.Delete(ctx, user.Name); err != nil {
				t.Error(err)
			}
		}
	})
	if userCount, passwordCount := len(users), len(passwords); userCount != passwordCount {
		t.Fatalf("len(users), len(passwords) = %v, %v; want equal", userCount, passwordCount)
	}
	for i, user := range users {
		if err := r.Create(ctx, user, passwords[i]); err != nil {
			t.Errorf("r.Create(ctx, %v, %q) = %v; want nil", user, passwords[i], err)
		}
	}
	if t.Failed() {
		t.FailNow()
	}
}
