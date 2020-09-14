package repositories

import (
	"testing"

	streckuv1 "github.com/Saser/strecku/saser/strecku/v1"
	"github.com/google/go-cmp/cmp"
	"google.golang.org/protobuf/testing/protocmp"
)

func toMap(users []*streckuv1.User) map[string]*streckuv1.User {
	m := make(map[string]*streckuv1.User)
	for _, user := range users {
		m[user.Name] = user
	}
	return m
}

func TestUserNotFoundError_Error(t *testing.T) {
	for _, test := range []struct {
		name string
		want string
	}{
		{name: "", want: `user not found: ""`},
		{name: "users/foobar", want: `user not found: "users/foobar"`},
		{name: "some name", want: `user not found: "some name"`},
	} {
		err := &UserNotFoundError{Name: test.name}
		if got := err.Error(); got != test.want {
			t.Errorf("err.Error() = %q; want %q", got, test.want)
		}
	}
}

func TestUsers_LookupUser(t *testing.T) {
	for _, test := range []struct {
		testName string
		users    []*streckuv1.User
		name     string
		wantUser *streckuv1.User
		wantErr  error
	}{
		{
			testName: "EmptyDatabaseEmptyName",
			users:    nil,
			name:     "",
			wantUser: nil,
			wantErr:  &UserNotFoundError{Name: ""},
		},
		{
			testName: "EmptyDatabaseNonEmptyName",
			users:    nil,
			name:     "users/foobar",
			wantUser: nil,
			wantErr:  &UserNotFoundError{Name: "users/foobar"},
		},
		{
			testName: "OneUserOK",
			users: []*streckuv1.User{
				{Name: "users/foobar", EmailAddress: "user@example.com", DisplayName: "User"},
			},
			name:     "users/foobar",
			wantUser: &streckuv1.User{Name: "users/foobar", EmailAddress: "user@example.com", DisplayName: "User"},
			wantErr:  nil,
		},
		{
			testName: "MultipleUsersOK",
			users: []*streckuv1.User{
				{Name: "users/foobar", EmailAddress: "foobar@example.com", DisplayName: "Foo Bar"},
				{Name: "users/barbaz", EmailAddress: "barbaz@example.com", DisplayName: "Barba Z."},
				{Name: "users/quux", EmailAddress: "quux@example.com", DisplayName: "Qu Ux"},
			},
			name:     "users/barbaz",
			wantUser: &streckuv1.User{Name: "users/barbaz", EmailAddress: "barbaz@example.com", DisplayName: "Barba Z."},
			wantErr:  nil,
		},
		{
			testName: "OneUserNotFound",
			users: []*streckuv1.User{
				{Name: "users/foobar", EmailAddress: "user@example.com", DisplayName: "User"},
			},
			name:     "users/notfoobar",
			wantUser: nil,
			wantErr:  &UserNotFoundError{Name: "users/notfoobar"},
		},
		{
			testName: "MultipleUsersNotFound",
			users: []*streckuv1.User{
				{Name: "users/foobar", EmailAddress: "foobar@example.com", DisplayName: "Foo Bar"},
				{Name: "users/barbaz", EmailAddress: "barbaz@example.com", DisplayName: "Barba Z."},
				{Name: "users/quux", EmailAddress: "quux@example.com", DisplayName: "Qu Ux"},
			},
			name:     "users/notfoobar",
			wantUser: nil,
			wantErr:  &UserNotFoundError{Name: "users/notfoobar"},
		},
	} {
		t.Run(test.testName, func(t *testing.T) {
			r := newUsers(toMap(test.users))
			user, err := r.LookupUser(test.name)
			if diff := cmp.Diff(user, test.wantUser, protocmp.Transform()); diff != "" {
				t.Errorf("user != test.wantUser (-got +want)\n%s", diff)
			}
			if got, want := err, test.wantErr; !cmp.Equal(got, want) {
				t.Errorf("r.LookupUser(%q) err = %v; want %v", test.name, got, want)
			}
		})
	}
}
