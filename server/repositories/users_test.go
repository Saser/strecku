package repositories

import (
	"strings"
	"testing"

	streckuv1 "github.com/Saser/strecku/saser/strecku/v1"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"google.golang.org/protobuf/testing/protocmp"
)

var userLess = func(u1, u2 *streckuv1.User) bool {
	return u1.Name < u2.Name
}

func seedUsers(users []*streckuv1.User) *Users {
	mUsers := make(map[string]*streckuv1.User, len(users))
	mNames := make(map[string]string, len(users))
	for _, user := range users {
		mUsers[user.Name] = user
		mNames[user.EmailAddress] = user.Name
	}
	return newUsers(mUsers, mNames)
}

func TestUserNotFoundError_Error(t *testing.T) {
	for _, test := range []struct {
		name         string
		emailAddress string
		want         string
	}{
		{name: "users/foobar", want: `user not found: "users/foobar"`},
		{name: "some name", want: `user not found: "some name"`},
		{emailAddress: "user@example.com", want: `user email not found: "user@example.com"`},
		{emailAddress: "some email", want: `user email not found: "some email"`},
	} {
		err := &UserNotFoundError{
			Name:         test.name,
			EmailAddress: test.emailAddress,
		}
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
			r := seedUsers(test.users)
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

func TestUsers_LookupUserByEmail(t *testing.T) {
	for _, test := range []struct {
		testName     string
		users        []*streckuv1.User
		emailAddress string
		wantUser     *streckuv1.User
		wantErr      error
	}{
		{
			testName:     "EmptyDatabaseEmptyName",
			users:        nil,
			emailAddress: "",
			wantUser:     nil,
			wantErr:      &UserNotFoundError{EmailAddress: ""},
		},
		{
			testName:     "EmptyDatabaseNonEmptyName",
			users:        nil,
			emailAddress: "user@example.com",
			wantUser:     nil,
			wantErr:      &UserNotFoundError{EmailAddress: "user@example.com"},
		},
		{
			testName: "OneUserOK",
			users: []*streckuv1.User{
				{Name: "users/foobar", EmailAddress: "user@example.com", DisplayName: "User"},
			},
			emailAddress: "user@example.com",
			wantUser:     &streckuv1.User{Name: "users/foobar", EmailAddress: "user@example.com", DisplayName: "User"},
			wantErr:      nil,
		},
		{
			testName: "MultipleUsersOK",
			users: []*streckuv1.User{
				{Name: "users/foobar", EmailAddress: "foobar@example.com", DisplayName: "Foo Bar"},
				{Name: "users/barbaz", EmailAddress: "barbaz@example.com", DisplayName: "Barba Z."},
				{Name: "users/quux", EmailAddress: "quux@example.com", DisplayName: "Qu Ux"},
			},
			emailAddress: "barbaz@example.com",
			wantUser:     &streckuv1.User{Name: "users/barbaz", EmailAddress: "barbaz@example.com", DisplayName: "Barba Z."},
			wantErr:      nil,
		},
		{
			testName: "OneUserNotFound",
			users: []*streckuv1.User{
				{Name: "users/foobar", EmailAddress: "user@example.com", DisplayName: "User"},
			},
			emailAddress: "notfoobar@example.com",
			wantUser:     nil,
			wantErr:      &UserNotFoundError{EmailAddress: "notfoobar@example.com"},
		},
		{
			testName: "MultipleUsersNotFound",
			users: []*streckuv1.User{
				{Name: "users/foobar", EmailAddress: "foobar@example.com", DisplayName: "Foo Bar"},
				{Name: "users/barbaz", EmailAddress: "barbaz@example.com", DisplayName: "Barba Z."},
				{Name: "users/quux", EmailAddress: "quux@example.com", DisplayName: "Qu Ux"},
			},
			emailAddress: "notfoobar@example.com",
			wantUser:     nil,
			wantErr:      &UserNotFoundError{EmailAddress: "notfoobar@example.com"},
		},
	} {
		t.Run(test.testName, func(t *testing.T) {
			r := seedUsers(test.users)
			user, err := r.LookupUserByEmail(test.emailAddress)
			if diff := cmp.Diff(user, test.wantUser, protocmp.Transform()); diff != "" {
				t.Errorf("user != test.wantUser (-got +want)\n%s", diff)
			}
			if got, want := err, test.wantErr; !cmp.Equal(got, want) {
				t.Errorf("r.LookupUser(%q) err = %v; want %v", test.emailAddress, got, want)
			}
		})
	}
}

func TestUsers_ListUsers(t *testing.T) {
	for _, test := range []struct {
		name  string
		users []*streckuv1.User
	}{
		{name: "Empty", users: nil},
		{
			name: "OneUser",
			users: []*streckuv1.User{
				{Name: "users/foobar", EmailAddress: "user@example.com", DisplayName: "Foo Bar"},
			},
		},
		{
			name: "ThreeUsers",
			users: []*streckuv1.User{
				{Name: "users/foobar", EmailAddress: "foobar@example.com", DisplayName: "Foo Bar"},
				{Name: "users/barbaz", EmailAddress: "barbaz@example.com", DisplayName: "Barba Z."},
				{Name: "users/quux", EmailAddress: "quux@example.com", DisplayName: "Qu Ux"},
			},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			r := seedUsers(test.users)
			users, err := r.ListUsers()
			if err != nil {
				t.Errorf("r.ListUsers() err = %v; want nil", err)
			}
			if diff := cmp.Diff(
				users, test.users, protocmp.Transform(),
				cmpopts.EquateEmpty(),
				cmpopts.SortSlices(userLess),
			); diff != "" {
				t.Errorf("users != test.users (-got +want)\n%s", diff)
			}
		})
	}
}

func TestUsers_FilterUsers(t *testing.T) {
	for _, test := range []struct {
		name      string
		users     []*streckuv1.User
		predicate func(*streckuv1.User) bool
		want      []*streckuv1.User
	}{
		{
			name:      "Empty",
			users:     nil,
			predicate: func(*streckuv1.User) bool { return true },
			want:      nil,
		},
		{
			name: "OneUserNoneMatching",
			users: []*streckuv1.User{
				{Name: "users/foobar", EmailAddress: "user@example.com", DisplayName: "Foo Bar"},
			},
			predicate: func(user *streckuv1.User) bool { return false },
			want:      nil,
		},
		{
			name: "MultipleUsersNoneMatching",
			users: []*streckuv1.User{
				{Name: "users/foobar", EmailAddress: "foobar@example.com", DisplayName: "Foo Bar"},
				{Name: "users/barbaz", EmailAddress: "barbaz@example.com", DisplayName: "Barba Z."},
				{Name: "users/quux", EmailAddress: "quux@example.com", DisplayName: "Qu Ux"},
			},
			predicate: func(user *streckuv1.User) bool { return false },
			want:      nil,
		},
		{
			name: "OneUserOneMatching",
			users: []*streckuv1.User{
				{Name: "users/foobar", EmailAddress: "user@example.com", DisplayName: "Foo Bar"},
			},
			predicate: func(user *streckuv1.User) bool { return strings.HasPrefix(user.DisplayName, "Foo") },
			want: []*streckuv1.User{
				{Name: "users/foobar", EmailAddress: "user@example.com", DisplayName: "Foo Bar"},
			},
		},
		{
			name: "MultipleUsersOneMatching",
			users: []*streckuv1.User{
				{Name: "users/foobar", EmailAddress: "foobar@example.com", DisplayName: "Foo Bar"},
				{Name: "users/barbaz", EmailAddress: "barbaz@example.com", DisplayName: "Barba Z."},
				{Name: "users/quux", EmailAddress: "quux@example.com", DisplayName: "Qu Ux"},
			},
			predicate: func(user *streckuv1.User) bool { return strings.HasPrefix(user.DisplayName, "Foo") },
			want: []*streckuv1.User{
				{Name: "users/foobar", EmailAddress: "foobar@example.com", DisplayName: "Foo Bar"},
			},
		},
		{
			name: "MultipleUsersMultipleMatching",
			users: []*streckuv1.User{
				{Name: "users/foobar", EmailAddress: "foobar@example.com", DisplayName: "Foo Bar"},
				{Name: "users/barbaz", EmailAddress: "barbaz@example.com", DisplayName: "Barba Z."},
				{Name: "users/quux", EmailAddress: "quux@example.com", DisplayName: "Qu Ux"},
			},
			predicate: func(user *streckuv1.User) bool { return strings.Contains(user.DisplayName, "Bar") },
			want: []*streckuv1.User{
				{Name: "users/foobar", EmailAddress: "foobar@example.com", DisplayName: "Foo Bar"},
				{Name: "users/barbaz", EmailAddress: "barbaz@example.com", DisplayName: "Barba Z."},
			},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			r := seedUsers(test.users)
			got, err := r.FilterUsers(test.predicate)
			if err != nil {
				t.Errorf("r.FilterUsers(test.predicate) err = %v; want nil", err)
			}
			if diff := cmp.Diff(
				got, test.want, protocmp.Transform(),
				cmpopts.EquateEmpty(),
				cmpopts.SortSlices(userLess),
			); diff != "" {
				t.Errorf("got != test.want (-got +want)\n%s", diff)
			}
		})
	}
}
