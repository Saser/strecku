package users

import (
	"context"
	"fmt"
	"strings"
	"testing"

	streckuv1 "github.com/Saser/strecku/saser/strecku/v1"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"google.golang.org/protobuf/testing/protocmp"
)

var userNames = map[string]string{
	"foobar": "users/6f2d193c-1460-491d-8157-7dd9535526c6",
	"barbaz": "users/d8bbf79e-8c59-4fae-aef9-634fcac00e07",
	"quux":   "users/9cd3ec05-e7af-418c-bd50-80a7c39a18cc",
	"cookie": "users/7f8f2c29-3860-49fa-923d-896a53f0ca26",
}

func userLess(u1, u2 *streckuv1.User) bool {
	return u1.Name < u2.Name
}

func seed(t *testing.T, users []*streckuv1.User, passwords []string) *Repository {
	t.Helper()
	userCount := len(users)
	if passwordCount := len(passwords); userCount != passwordCount {
		t.Fatalf("len(users), len(passwords) = %v, %v; want them to be equal", userCount, passwordCount)
	}
	mUsers := make(map[string]*streckuv1.User, userCount)
	mPasswords := make(map[string]string, userCount)
	mNames := make(map[string]string, userCount)
	for i, user := range users {
		if got := Validate(user); got != nil {
			t.Errorf("Validate(%v) = %v; want %v", user, got, nil)
		}
		mUsers[user.Name] = user
		mPasswords[user.Name] = passwords[i]
		mNames[user.EmailAddress] = user.Name
	}
	return newRepository(mUsers, mPasswords, mNames)
}

func TestUserNotFoundError_Error(t *testing.T) {
	for _, test := range []struct {
		name         string
		emailAddress string
		want         string
	}{
		{name: userNames["foobar"], want: fmt.Sprintf("user not found: %q", userNames["foobar"])},
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

func TestUserNotFoundError_Is(t *testing.T) {
	for _, test := range []struct {
		err    *UserNotFoundError
		target error
		want   bool
	}{
		{
			err:    &UserNotFoundError{Name: userNames["foobar"]},
			target: &UserNotFoundError{Name: userNames["foobar"]},
			want:   true,
		},
		{
			err:    &UserNotFoundError{Name: userNames["foobar"]},
			target: &UserNotFoundError{Name: userNames["barbaz"]},
			want:   false,
		},
		{
			err:    &UserNotFoundError{Name: userNames["foobar"]},
			target: &UserNotFoundError{EmailAddress: "foobar@example.com"},
			want:   false,
		},
		{
			err:    &UserNotFoundError{EmailAddress: "foobar@example.com"},
			target: &UserNotFoundError{Name: userNames["foobar"]},
			want:   false,
		},
		{
			err:    &UserNotFoundError{EmailAddress: "foobar@example.com"},
			target: &UserNotFoundError{EmailAddress: "foobar@example.com"},
			want:   true,
		},
		{
			err:    &UserNotFoundError{EmailAddress: "foobar@example.com"},
			target: &UserNotFoundError{EmailAddress: "barbaz@example.com"},
			want:   false,
		},
		{
			err:    &UserNotFoundError{Name: userNames["foobar"]},
			target: fmt.Errorf("user not found: %q", userNames["foobar"]),
			want:   false,
		},
		{
			err:    &UserNotFoundError{EmailAddress: "foobar@example.com"},
			target: fmt.Errorf("user email not found: %q", "foobar@example.com"),
			want:   false,
		},
	} {
		if got := test.err.Is(test.target); got != test.want {
			t.Errorf("test.err.Is(%v) = %v; want %v", test.target, got, test.want)
		}
	}
}

func TestUserExistsError_Error(t *testing.T) {
	for _, test := range []struct {
		name         string
		emailAddress string
		want         string
	}{
		{name: userNames["foobar"], want: fmt.Sprintf("user exists: %q", userNames["foobar"])},
		{name: "some name", want: `user exists: "some name"`},
		{emailAddress: "user@example.com", want: `user email exists: "user@example.com"`},
		{emailAddress: "some email", want: `user email exists: "some email"`},
	} {
		err := &UserExistsError{
			Name:         test.name,
			EmailAddress: test.emailAddress,
		}
		if got := err.Error(); got != test.want {
			t.Errorf("err.Error() = %q; want %q", got, test.want)
		}
	}
}

func TestUserExistsError_Is(t *testing.T) {
	for _, test := range []struct {
		err    *UserExistsError
		target error
		want   bool
	}{
		{
			err:    &UserExistsError{Name: userNames["foobar"]},
			target: &UserExistsError{Name: userNames["foobar"]},
			want:   true,
		},
		{
			err:    &UserExistsError{Name: userNames["foobar"]},
			target: &UserExistsError{Name: userNames["barbaz"]},
			want:   false,
		},
		{
			err:    &UserExistsError{Name: userNames["foobar"]},
			target: &UserExistsError{EmailAddress: "foobar@example.com"},
			want:   false,
		},
		{
			err:    &UserExistsError{EmailAddress: "foobar@example.com"},
			target: &UserExistsError{Name: userNames["foobar"]},
			want:   false,
		},
		{
			err:    &UserExistsError{EmailAddress: "foobar@example.com"},
			target: &UserExistsError{EmailAddress: "foobar@example.com"},
			want:   true,
		},
		{
			err:    &UserExistsError{EmailAddress: "foobar@example.com"},
			target: &UserExistsError{EmailAddress: "barbaz@example.com"},
			want:   false,
		},
		{
			err:    &UserExistsError{Name: userNames["foobar"]},
			target: fmt.Errorf("user exists: %q", userNames["foobar"]),
			want:   false,
		},
		{
			err:    &UserExistsError{EmailAddress: "foobar@example.com"},
			target: fmt.Errorf("user email exists: %q", "foobar@example.com"),
			want:   false,
		},
	} {
		if got := test.err.Is(test.target); got != test.want {
			t.Errorf("test.err.Is(%v) = %v; want %v", test.target, got, test.want)
		}
	}
}

func TestWrongPasswordError_Error(t *testing.T) {
	err := &WrongPasswordError{Name: userNames["foobar"]}
	if got, want := err.Error(), fmt.Sprintf("wrong password for user %q", userNames["foobar"]); got != want {
		t.Errorf("err.Error() = %q; want %q", got, want)
	}
}

func TestWrongPasswordError_Is(t *testing.T) {
	for _, test := range []struct {
		err    *WrongPasswordError
		target error
		want   bool
	}{
		{
			err:    &WrongPasswordError{Name: userNames["foobar"]},
			target: &WrongPasswordError{Name: userNames["foobar"]},
			want:   true,
		},
		{
			err:    &WrongPasswordError{Name: userNames["foobar"]},
			target: &WrongPasswordError{Name: userNames["barbaz"]},
			want:   false,
		},
		{
			err:    &WrongPasswordError{Name: userNames["foobar"]},
			target: fmt.Errorf("wrong password for user %q", userNames["foobar"]),
			want:   false,
		},
	} {
		if got := test.err.Is(test.target); got != test.want {
			t.Errorf("test.err.Is(%v) = %v; want %v", test.target, got, test.want)
		}
	}
}

func TestUsers_Authenticate(t *testing.T) {
	ctx := context.Background()
	for _, test := range []struct {
		desc      string
		users     []*streckuv1.User
		passwords []string
		name      string
		password  string
		want      error
	}{
		{
			desc:      "Empty",
			users:     nil,
			passwords: nil,
			name:      userNames["foobar"],
			password:  "foobar",
			want:      &UserNotFoundError{Name: userNames["foobar"]},
		},
		{
			desc: "OneUserOK",
			users: []*streckuv1.User{
				{Name: userNames["foobar"], EmailAddress: "foobar@example.com", DisplayName: "Foo Bar"},
			},
			passwords: []string{
				"foobar",
			},
			name:     userNames["foobar"],
			password: "foobar",
			want:     nil,
		},
		{
			desc: "MultipleUsersOK",
			users: []*streckuv1.User{
				{Name: userNames["foobar"], EmailAddress: "foobar@example.com", DisplayName: "Foo Bar"},
				{Name: userNames["barbaz"], EmailAddress: "barbaz@example.com", DisplayName: "Barba Z."},
				{Name: userNames["quux"], EmailAddress: "quux@example.com", DisplayName: "Qu Ux"},
			},
			passwords: []string{
				"foobar",
				"barbaz",
				"quux",
			},
			name:     userNames["barbaz"],
			password: "barbaz",
			want:     nil,
		},
		{
			desc: "OneUserNotFound",
			users: []*streckuv1.User{
				{Name: userNames["foobar"], EmailAddress: "foobar@example.com", DisplayName: "Foo Bar"},
			},
			passwords: []string{
				"foobar",
			},
			name:     userNames["barbaz"],
			password: "barbaz",
			want:     &UserNotFoundError{Name: userNames["barbaz"]},
		},
		{
			desc: "MultipleUsersNotFound",
			users: []*streckuv1.User{
				{Name: userNames["foobar"], EmailAddress: "foobar@example.com", DisplayName: "Foo Bar"},
				{Name: userNames["barbaz"], EmailAddress: "barbaz@example.com", DisplayName: "Barba Z."},
				{Name: userNames["quux"], EmailAddress: "quux@example.com", DisplayName: "Qu Ux"},
			},
			passwords: []string{
				"foobar",
				"barbaz",
				"quux",
			},
			name:     userNames["cookie"],
			password: "cookie",
			want:     &UserNotFoundError{Name: userNames["cookie"]},
		},
		{
			desc: "OneUserWrongPassword",
			users: []*streckuv1.User{
				{Name: userNames["foobar"], EmailAddress: "foobar@example.com", DisplayName: "Foo Bar"},
			},
			passwords: []string{
				"foobar",
			},
			name:     userNames["foobar"],
			password: "wrong password",
			want:     &WrongPasswordError{Name: userNames["foobar"]},
		},
		{
			desc: "MultipleUsersWrongPassword",
			users: []*streckuv1.User{
				{Name: userNames["foobar"], EmailAddress: "foobar@example.com", DisplayName: "Foo Bar"},
				{Name: userNames["barbaz"], EmailAddress: "barbaz@example.com", DisplayName: "Barba Z."},
				{Name: userNames["quux"], EmailAddress: "quux@example.com", DisplayName: "Qu Ux"},
			},
			passwords: []string{
				"foobar",
				"barbaz",
				"quux",
			},
			name:     userNames["foobar"],
			password: "wrong password",
			want:     &WrongPasswordError{Name: userNames["foobar"]},
		},
	} {
		t.Run(test.desc, func(t *testing.T) {
			r := seed(t, test.users, test.passwords)
			if got := r.Authenticate(ctx, test.name, test.password); !cmp.Equal(got, test.want) {
				t.Errorf("r.Authenticate(%v, %q, %q) = %v; want %v", ctx, test.name, test.password, got, test.want)
			}
		})
	}
}

func TestUsers_LookupUser(t *testing.T) {
	ctx := context.Background()
	for _, test := range []struct {
		desc      string
		users     []*streckuv1.User
		passwords []string
		name      string
		wantUser  *streckuv1.User
		wantErr   error
	}{
		{
			desc:      "EmptyDatabaseEmptyName",
			users:     nil,
			passwords: nil,
			name:      "",
			wantUser:  nil,
			wantErr:   &UserNotFoundError{Name: ""},
		},
		{
			desc:      "EmptyDatabaseNonEmptyName",
			users:     nil,
			passwords: nil,
			name:      userNames["foobar"],
			wantUser:  nil,
			wantErr:   &UserNotFoundError{Name: userNames["foobar"]},
		},
		{
			desc: "OneUserOK",
			users: []*streckuv1.User{
				{Name: userNames["foobar"], EmailAddress: "user@example.com", DisplayName: "User"},
			},
			passwords: []string{
				"foobar",
			},
			name:     userNames["foobar"],
			wantUser: &streckuv1.User{Name: userNames["foobar"], EmailAddress: "user@example.com", DisplayName: "User"},
			wantErr:  nil,
		},
		{
			desc: "MultipleUsersOK",
			users: []*streckuv1.User{
				{Name: userNames["foobar"], EmailAddress: "foobar@example.com", DisplayName: "Foo Bar"},
				{Name: userNames["barbaz"], EmailAddress: "barbaz@example.com", DisplayName: "Barba Z."},
				{Name: userNames["quux"], EmailAddress: "quux@example.com", DisplayName: "Qu Ux"},
			},
			passwords: []string{
				"foobar",
				"barbaz",
				"quux",
			},
			name:     userNames["barbaz"],
			wantUser: &streckuv1.User{Name: userNames["barbaz"], EmailAddress: "barbaz@example.com", DisplayName: "Barba Z."},
			wantErr:  nil,
		},
		{
			desc: "OneUserNotFound",
			users: []*streckuv1.User{
				{Name: userNames["foobar"], EmailAddress: "user@example.com", DisplayName: "User"},
			},
			passwords: []string{"foobar"},
			name:      userNames["barbaz"],
			wantUser:  nil,
			wantErr:   &UserNotFoundError{Name: userNames["barbaz"]},
		},
		{
			desc: "MultipleUsersNotFound",
			users: []*streckuv1.User{
				{Name: userNames["foobar"], EmailAddress: "foobar@example.com", DisplayName: "Foo Bar"},
				{Name: userNames["barbaz"], EmailAddress: "barbaz@example.com", DisplayName: "Barba Z."},
				{Name: userNames["quux"], EmailAddress: "quux@example.com", DisplayName: "Qu Ux"},
			},
			passwords: []string{
				"foobar",
				"barbaz",
				"quux",
			},
			name:     userNames["cookie"],
			wantUser: nil,
			wantErr:  &UserNotFoundError{Name: userNames["cookie"]},
		},
	} {
		t.Run(test.desc, func(t *testing.T) {
			r := seed(t, test.users, test.passwords)
			user, err := r.LookupUser(ctx, test.name)
			if diff := cmp.Diff(user, test.wantUser, protocmp.Transform()); diff != "" {
				t.Errorf("r.LookupUser(%v, %q) user != test.wantUser (-got +want)\n%s", ctx, test.name, diff)
			}
			if got, want := err, test.wantErr; !cmp.Equal(got, want) {
				t.Errorf("r.LookupUser(%v, %q) err = %v; want %v", ctx, test.name, got, want)
			}
		})
	}
}

func TestUsers_LookupUserByEmail(t *testing.T) {
	ctx := context.Background()
	for _, test := range []struct {
		desc         string
		users        []*streckuv1.User
		passwords    []string
		emailAddress string
		wantUser     *streckuv1.User
		wantErr      error
	}{
		{
			desc:         "EmptyDatabaseEmptyName",
			users:        nil,
			passwords:    nil,
			emailAddress: "",
			wantUser:     nil,
			wantErr:      &UserNotFoundError{EmailAddress: ""},
		},
		{
			desc:         "EmptyDatabaseNonEmptyName",
			users:        nil,
			passwords:    nil,
			emailAddress: "user@example.com",
			wantUser:     nil,
			wantErr:      &UserNotFoundError{EmailAddress: "user@example.com"},
		},
		{
			desc: "OneUserOK",
			users: []*streckuv1.User{
				{Name: userNames["foobar"], EmailAddress: "user@example.com", DisplayName: "User"},
			},
			passwords: []string{
				"foobar",
			},
			emailAddress: "user@example.com",
			wantUser:     &streckuv1.User{Name: userNames["foobar"], EmailAddress: "user@example.com", DisplayName: "User"},
			wantErr:      nil,
		},
		{
			desc: "MultipleUsersOK",
			users: []*streckuv1.User{
				{Name: userNames["foobar"], EmailAddress: "foobar@example.com", DisplayName: "Foo Bar"},
				{Name: userNames["barbaz"], EmailAddress: "barbaz@example.com", DisplayName: "Barba Z."},
				{Name: userNames["quux"], EmailAddress: "quux@example.com", DisplayName: "Qu Ux"},
			},
			passwords: []string{
				"foobar",
				"barbaz",
				"quux",
			},
			emailAddress: "barbaz@example.com",
			wantUser:     &streckuv1.User{Name: userNames["barbaz"], EmailAddress: "barbaz@example.com", DisplayName: "Barba Z."},
			wantErr:      nil,
		},
		{
			desc: "OneUserNotFound",
			users: []*streckuv1.User{
				{Name: userNames["foobar"], EmailAddress: "user@example.com", DisplayName: "User"},
			},
			passwords:    []string{"foobar"},
			emailAddress: "notfoobar@example.com",
			wantUser:     nil,
			wantErr:      &UserNotFoundError{EmailAddress: "notfoobar@example.com"},
		},
		{
			desc: "MultipleUsersNotFound",
			users: []*streckuv1.User{
				{Name: userNames["foobar"], EmailAddress: "foobar@example.com", DisplayName: "Foo Bar"},
				{Name: userNames["barbaz"], EmailAddress: "barbaz@example.com", DisplayName: "Barba Z."},
				{Name: userNames["quux"], EmailAddress: "quux@example.com", DisplayName: "Qu Ux"},
			},
			passwords: []string{
				"foobar",
				"barbaz",
				"quux",
			},
			emailAddress: "notfoobar@example.com",
			wantUser:     nil,
			wantErr:      &UserNotFoundError{EmailAddress: "notfoobar@example.com"},
		},
	} {
		t.Run(test.desc, func(t *testing.T) {
			r := seed(t, test.users, test.passwords)
			user, err := r.LookupUserByEmail(ctx, test.emailAddress)
			if diff := cmp.Diff(user, test.wantUser, protocmp.Transform()); diff != "" {
				t.Errorf("r.LookupUserByEmail(%v, %q) user != test.wantUser (-got +want)\n%s", ctx, test.emailAddress, diff)
			}
			if got, want := err, test.wantErr; !cmp.Equal(got, want) {
				t.Errorf("r.LookupUserByEmail(%v, %q) err = %v; want %v", ctx, test.emailAddress, got, want)
			}
		})
	}
}

func TestUsers_ListUsers(t *testing.T) {
	ctx := context.Background()
	for _, test := range []struct {
		name      string
		users     []*streckuv1.User
		passwords []string
	}{
		{name: "Empty", users: nil},
		{
			name: "OneUser",
			users: []*streckuv1.User{
				{Name: userNames["foobar"], EmailAddress: "user@example.com", DisplayName: "Foo Bar"},
			},
			passwords: []string{"foobar"},
		},
		{
			name: "ThreeUsers",
			users: []*streckuv1.User{
				{Name: userNames["foobar"], EmailAddress: "foobar@example.com", DisplayName: "Foo Bar"},
				{Name: userNames["barbaz"], EmailAddress: "barbaz@example.com", DisplayName: "Barba Z."},
				{Name: userNames["quux"], EmailAddress: "quux@example.com", DisplayName: "Qu Ux"},
			},
			passwords: []string{
				"foobar",
				"barbaz",
				"quux",
			},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			r := seed(t, test.users, test.passwords)
			users, err := r.ListUsers(ctx)
			if diff := cmp.Diff(
				users, test.users, protocmp.Transform(),
				cmpopts.EquateEmpty(),
				cmpopts.SortSlices(userLess),
			); diff != "" {
				t.Errorf("r.ListUsers(%v) users != test.users (-got +want)\n%s", ctx, diff)
			}
			if got, want := err, error(nil); !cmp.Equal(got, want) {
				t.Errorf("r.ListUsers(%v) err = %v; want %v", ctx, got, want)
			}
		})
	}
}

func TestUsers_FilterUsers(t *testing.T) {
	ctx := context.Background()
	for _, test := range []struct {
		name      string
		users     []*streckuv1.User
		passwords []string
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
				{Name: userNames["foobar"], EmailAddress: "user@example.com", DisplayName: "Foo Bar"},
			},
			passwords: []string{"foobar"},
			predicate: func(user *streckuv1.User) bool { return false },
			want:      nil,
		},
		{
			name: "MultipleUsersNoneMatching",
			users: []*streckuv1.User{
				{Name: userNames["foobar"], EmailAddress: "foobar@example.com", DisplayName: "Foo Bar"},
				{Name: userNames["barbaz"], EmailAddress: "barbaz@example.com", DisplayName: "Barba Z."},
				{Name: userNames["quux"], EmailAddress: "quux@example.com", DisplayName: "Qu Ux"},
			},
			passwords: []string{
				"foobar",
				"barbaz",
				"quux",
			},
			predicate: func(user *streckuv1.User) bool { return false },
			want:      nil,
		},
		{
			name: "OneUserOneMatching",
			users: []*streckuv1.User{
				{Name: userNames["foobar"], EmailAddress: "user@example.com", DisplayName: "Foo Bar"},
			},
			passwords: []string{"foobar"},
			predicate: func(user *streckuv1.User) bool { return strings.HasPrefix(user.DisplayName, "Foo") },
			want: []*streckuv1.User{
				{Name: userNames["foobar"], EmailAddress: "user@example.com", DisplayName: "Foo Bar"},
			},
		},
		{
			name: "MultipleUsersOneMatching",
			users: []*streckuv1.User{
				{Name: userNames["foobar"], EmailAddress: "foobar@example.com", DisplayName: "Foo Bar"},
				{Name: userNames["barbaz"], EmailAddress: "barbaz@example.com", DisplayName: "Barba Z."},
				{Name: userNames["quux"], EmailAddress: "quux@example.com", DisplayName: "Qu Ux"},
			},
			passwords: []string{
				"foobar",
				"barbaz",
				"quux",
			},
			predicate: func(user *streckuv1.User) bool { return strings.HasPrefix(user.DisplayName, "Foo") },
			want: []*streckuv1.User{
				{Name: userNames["foobar"], EmailAddress: "foobar@example.com", DisplayName: "Foo Bar"},
			},
		},
		{
			name: "MultipleUsersMultipleMatching",
			users: []*streckuv1.User{
				{Name: userNames["foobar"], EmailAddress: "foobar@example.com", DisplayName: "Foo Bar"},
				{Name: userNames["barbaz"], EmailAddress: "barbaz@example.com", DisplayName: "Barba Z."},
				{Name: userNames["quux"], EmailAddress: "quux@example.com", DisplayName: "Qu Ux"},
			},
			passwords: []string{
				"foobar",
				"barbaz",
				"quux",
			},
			predicate: func(user *streckuv1.User) bool { return strings.Contains(user.DisplayName, "Bar") },
			want: []*streckuv1.User{
				{Name: userNames["foobar"], EmailAddress: "foobar@example.com", DisplayName: "Foo Bar"},
				{Name: userNames["barbaz"], EmailAddress: "barbaz@example.com", DisplayName: "Barba Z."},
			},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			r := seed(t, test.users, test.passwords)
			users, err := r.FilterUsers(ctx, test.predicate)
			if diff := cmp.Diff(
				users, test.want, protocmp.Transform(),
				cmpopts.EquateEmpty(),
				cmpopts.SortSlices(userLess),
			); diff != "" {
				t.Errorf("r.FilterUsers(%v, test.predicate) users != test.want (-got +want)\n%s", ctx, diff)
			}
			if got, want := err, error(nil); !cmp.Equal(got, want) {
				t.Errorf("r.FilterUsers(%v, test.predicate) err = %v; want %v", ctx, got, want)
			}
		})
	}
}

func TestUsers_CreateUser(t *testing.T) {
	ctx := context.Background()
	for _, test := range []struct {
		name      string
		users     []*streckuv1.User
		passwords []string
		user      *streckuv1.User
		want      error
	}{
		{
			name:  "Empty",
			users: nil,
			user:  &streckuv1.User{Name: userNames["foobar"], EmailAddress: "foobar@example.com", DisplayName: "Foo Bar"},
			want:  nil,
		},
		{
			name: "OneUserOK",
			users: []*streckuv1.User{
				{Name: userNames["foobar"], EmailAddress: "foobar@example.com", DisplayName: "Foo Bar"},
			},
			passwords: []string{"foobar"},
			user:      &streckuv1.User{Name: userNames["barbaz"], EmailAddress: "barbaz@example.com", DisplayName: "Foo Bar"},
			want:      nil,
		},
		{
			name: "MultipleUsersOK",
			users: []*streckuv1.User{
				{Name: userNames["foobar"], EmailAddress: "foobar@example.com", DisplayName: "Foo Bar"},
				{Name: userNames["barbaz"], EmailAddress: "barbaz@example.com", DisplayName: "Barba Z."},
				{Name: userNames["quux"], EmailAddress: "quux@example.com", DisplayName: "Qu Ux"},
			},
			passwords: []string{
				"foobar",
				"barbaz",
				"quux",
			},
			user: &streckuv1.User{Name: userNames["cookie"], EmailAddress: "cookie@example.com", DisplayName: "Cookie"},
			want: nil,
		},
		{
			name: "OneUserDuplicateEmail",
			users: []*streckuv1.User{
				{Name: userNames["foobar"], EmailAddress: "foobar@example.com", DisplayName: "Foo Bar"},
			},
			passwords: []string{"foobar"},
			user:      &streckuv1.User{Name: userNames["barbaz"], EmailAddress: "foobar@example.com", DisplayName: "Barba Z."},
			want:      &UserExistsError{EmailAddress: "foobar@example.com"},
		},
		{
			name: "OneUserDuplicateName",
			users: []*streckuv1.User{
				{Name: userNames["foobar"], EmailAddress: "foobar@example.com", DisplayName: "Foo Bar"},
			},
			passwords: []string{"foobar"},
			user:      &streckuv1.User{Name: userNames["foobar"], EmailAddress: "new-foobar@example.com", DisplayName: "New Foo Bar"},
			want:      &UserExistsError{Name: userNames["foobar"]},
		},
		{
			name: "MultipleUsersDuplicateEmail",
			users: []*streckuv1.User{
				{Name: userNames["foobar"], EmailAddress: "foobar@example.com", DisplayName: "Foo Bar"},
				{Name: userNames["barbaz"], EmailAddress: "barbaz@example.com", DisplayName: "Barba Z."},
				{Name: userNames["quux"], EmailAddress: "quux@example.com", DisplayName: "Qu Ux"},
			},
			passwords: []string{
				"foobar",
				"barbaz",
				"quux",
			},
			user: &streckuv1.User{Name: userNames["cookie"], EmailAddress: "foobar@example.com", DisplayName: "Cookie"},
			want: &UserExistsError{EmailAddress: "foobar@example.com"},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			r := seed(t, test.users, test.passwords)
			if got := r.CreateUser(ctx, test.user); !cmp.Equal(got, test.want) {
				t.Errorf("r.CreateUser(%v, %v) = %v; want %v", ctx, test.user, got, test.want)
			}
		})
	}
}

func TestUsers_UpdateUser(t *testing.T) {
	ctx := context.Background()
	for _, test := range []struct {
		name          string
		users         []*streckuv1.User
		passwords     []string
		updated       *streckuv1.User
		wantUpdateErr error
		lookupEmail   string
		wantUser      *streckuv1.User
		wantLookupErr error
	}{
		{
			name: "OneUserEmailAddressLookupNew",
			users: []*streckuv1.User{
				{Name: userNames["foobar"], EmailAddress: "foobar@example.com", DisplayName: "Foo Bar"},
			},
			passwords:     []string{"foobar"},
			updated:       &streckuv1.User{Name: userNames["foobar"], EmailAddress: "new-foobar@example.com", DisplayName: "Foo Bar"},
			wantUpdateErr: nil,
			lookupEmail:   "new-foobar@example.com",
			wantUser:      &streckuv1.User{Name: userNames["foobar"], EmailAddress: "new-foobar@example.com", DisplayName: "Foo Bar"},
			wantLookupErr: nil,
		},
		{
			name: "OneUserEmailAddressLookupOld",
			users: []*streckuv1.User{
				{Name: userNames["foobar"], EmailAddress: "foobar@example.com", DisplayName: "Foo Bar"},
			},
			passwords:     []string{"foobar"},
			updated:       &streckuv1.User{Name: userNames["foobar"], EmailAddress: "new-foobar@example.com", DisplayName: "Foo Bar"},
			wantUpdateErr: nil,
			lookupEmail:   "foobar@example.com",
			wantUser:      nil,
			wantLookupErr: &UserNotFoundError{EmailAddress: "foobar@example.com"},
		},
		{
			name: "OneUserDisplayName",
			users: []*streckuv1.User{
				{Name: userNames["foobar"], EmailAddress: "foobar@example.com", DisplayName: "Foo Bar"},
			},
			passwords:     []string{"foobar"},
			updated:       &streckuv1.User{Name: userNames["foobar"], EmailAddress: "foobar@example.com", DisplayName: "New Foo Bar"},
			wantUpdateErr: nil,
			lookupEmail:   "foobar@example.com",
			wantUser:      &streckuv1.User{Name: userNames["foobar"], EmailAddress: "foobar@example.com", DisplayName: "New Foo Bar"},
			wantLookupErr: nil,
		},
		{
			name: "OneUserMultipleFields",
			users: []*streckuv1.User{
				{Name: userNames["foobar"], EmailAddress: "foobar@example.com", DisplayName: "Foo Bar"},
			},
			passwords:     []string{"foobar"},
			updated:       &streckuv1.User{Name: userNames["foobar"], EmailAddress: "new-foobar@example.com", DisplayName: "New Foo Bar"},
			wantUpdateErr: nil,
			lookupEmail:   "new-foobar@example.com",
			wantUser:      &streckuv1.User{Name: userNames["foobar"], EmailAddress: "new-foobar@example.com", DisplayName: "New Foo Bar"},
			wantLookupErr: nil,
		},
		{
			name: "OneUserNotFound",
			users: []*streckuv1.User{
				{Name: userNames["foobar"], EmailAddress: "foobar@example.com", DisplayName: "Foo Bar"},
			},
			passwords:     []string{"foobar"},
			updated:       &streckuv1.User{Name: userNames["barbaz"], EmailAddress: "new-foobar@example.com", DisplayName: "Foo Bar"},
			wantUpdateErr: &UserNotFoundError{Name: userNames["barbaz"]},
			lookupEmail:   "new-foobar@example.com",
			wantUser:      nil,
			wantLookupErr: &UserNotFoundError{EmailAddress: "new-foobar@example.com"},
		},
		{
			name: "MultipleUsersEmailAddressLookupNew",
			users: []*streckuv1.User{
				{Name: userNames["foobar"], EmailAddress: "foobar@example.com", DisplayName: "Foo Bar"},
				{Name: userNames["barbaz"], EmailAddress: "barbaz@example.com", DisplayName: "Barba Z."},
				{Name: userNames["quux"], EmailAddress: "quux@example.com", DisplayName: "Qu Ux"},
			},
			passwords: []string{
				"foobar",
				"barbaz",
				"quux",
			},
			updated:       &streckuv1.User{Name: userNames["barbaz"], EmailAddress: "new-barbaz@example.com", DisplayName: "Barba Z."},
			wantUpdateErr: nil,
			lookupEmail:   "new-barbaz@example.com",
			wantUser:      &streckuv1.User{Name: userNames["barbaz"], EmailAddress: "new-barbaz@example.com", DisplayName: "Barba Z."},
			wantLookupErr: nil,
		},
		{
			name: "MultipleUsersEmailAddressLookupOld",
			users: []*streckuv1.User{
				{Name: userNames["foobar"], EmailAddress: "foobar@example.com", DisplayName: "Foo Bar"},
				{Name: userNames["barbaz"], EmailAddress: "barbaz@example.com", DisplayName: "Barba Z."},
				{Name: userNames["quux"], EmailAddress: "quux@example.com", DisplayName: "Qu Ux"},
			},
			passwords: []string{
				"foobar",
				"barbaz",
				"quux",
			},
			updated:       &streckuv1.User{Name: userNames["barbaz"], EmailAddress: "new-barbaz@example.com", DisplayName: "Barba Z."},
			wantUpdateErr: nil,
			lookupEmail:   "barbaz@example.com",
			wantUser:      nil,
			wantLookupErr: &UserNotFoundError{EmailAddress: "barbaz@example.com"},
		},
		{
			name: "MultipleUsersDisplayName",
			users: []*streckuv1.User{
				{Name: userNames["foobar"], EmailAddress: "foobar@example.com", DisplayName: "Foo Bar"},
				{Name: userNames["barbaz"], EmailAddress: "barbaz@example.com", DisplayName: "Barba Z."},
				{Name: userNames["quux"], EmailAddress: "quux@example.com", DisplayName: "Qu Ux"},
			},
			passwords: []string{
				"foobar",
				"barbaz",
				"quux",
			},
			updated:       &streckuv1.User{Name: userNames["barbaz"], EmailAddress: "barbaz@example.com", DisplayName: "All New Barba Z."},
			wantUpdateErr: nil,
			lookupEmail:   "barbaz@example.com",
			wantUser:      &streckuv1.User{Name: userNames["barbaz"], EmailAddress: "barbaz@example.com", DisplayName: "All New Barba Z."},
			wantLookupErr: nil,
		},
		{
			name: "MultipleUsersMultipleFields",
			users: []*streckuv1.User{
				{Name: userNames["foobar"], EmailAddress: "foobar@example.com", DisplayName: "Foo Bar"},
				{Name: userNames["barbaz"], EmailAddress: "barbaz@example.com", DisplayName: "Barba Z."},
				{Name: userNames["quux"], EmailAddress: "quux@example.com", DisplayName: "Qu Ux"},
			},
			passwords: []string{
				"foobar",
				"barbaz",
				"quux",
			},
			updated:       &streckuv1.User{Name: userNames["barbaz"], EmailAddress: "new-barbaz@example.com", DisplayName: "All New Barba Z."},
			wantUpdateErr: nil,
			lookupEmail:   "new-barbaz@example.com",
			wantUser:      &streckuv1.User{Name: userNames["barbaz"], EmailAddress: "new-barbaz@example.com", DisplayName: "All New Barba Z."},
			wantLookupErr: nil,
		},
		{
			name: "MultipleUsersNotFound",
			users: []*streckuv1.User{
				{Name: userNames["foobar"], EmailAddress: "foobar@example.com", DisplayName: "Foo Bar"},
				{Name: userNames["barbaz"], EmailAddress: "barbaz@example.com", DisplayName: "Barba Z."},
				{Name: userNames["quux"], EmailAddress: "quux@example.com", DisplayName: "Qu Ux"},
			},
			passwords: []string{
				"foobar",
				"barbaz",
				"quux",
			},
			updated:       &streckuv1.User{Name: userNames["cookie"], EmailAddress: "new-barbaz@example.com", DisplayName: "Barba Z."},
			wantUpdateErr: &UserNotFoundError{Name: userNames["cookie"]},
			lookupEmail:   "barbaz@example.com",
			wantUser:      &streckuv1.User{Name: userNames["barbaz"], EmailAddress: "barbaz@example.com", DisplayName: "Barba Z."},
			wantLookupErr: nil,
		},
		{
			name: "MultipleUsersDuplicateEmailAddressLookupNew",
			users: []*streckuv1.User{
				{Name: userNames["foobar"], EmailAddress: "foobar@example.com", DisplayName: "Foo Bar"},
				{Name: userNames["barbaz"], EmailAddress: "barbaz@example.com", DisplayName: "Barba Z."},
				{Name: userNames["quux"], EmailAddress: "quux@example.com", DisplayName: "Qu Ux"},
			},
			passwords: []string{
				"foobar",
				"barbaz",
				"quux",
			},
			updated:       &streckuv1.User{Name: userNames["barbaz"], EmailAddress: "foobar@example.com", DisplayName: "Barba Z."},
			wantUpdateErr: &UserExistsError{EmailAddress: "foobar@example.com"},
			lookupEmail:   "foobar@example.com",
			wantUser:      &streckuv1.User{Name: userNames["foobar"], EmailAddress: "foobar@example.com", DisplayName: "Foo Bar"},
			wantLookupErr: nil,
		},
		{
			name: "MultipleUsersDuplicateEmailAddressLookupOld",
			users: []*streckuv1.User{
				{Name: userNames["foobar"], EmailAddress: "foobar@example.com", DisplayName: "Foo Bar"},
				{Name: userNames["barbaz"], EmailAddress: "barbaz@example.com", DisplayName: "Barba Z."},
				{Name: userNames["quux"], EmailAddress: "quux@example.com", DisplayName: "Qu Ux"},
			},
			passwords: []string{
				"foobar",
				"barbaz",
				"quux",
			},
			updated:       &streckuv1.User{Name: userNames["barbaz"], EmailAddress: "foobar@example.com", DisplayName: "All New Barba Z."},
			wantUpdateErr: &UserExistsError{EmailAddress: "foobar@example.com"},
			lookupEmail:   "barbaz@example.com",
			wantUser:      &streckuv1.User{Name: userNames["barbaz"], EmailAddress: "barbaz@example.com", DisplayName: "Barba Z."},
			wantLookupErr: nil,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			r := seed(t, test.users, test.passwords)
			if got := r.UpdateUser(ctx, test.updated); !cmp.Equal(got, test.wantUpdateErr) {
				t.Errorf("r.UpdateUser(%v, %v) = %v; want %v", ctx, test.updated, got, test.wantUpdateErr)
			}
			user, err := r.LookupUserByEmail(ctx, test.lookupEmail)
			if diff := cmp.Diff(user, test.wantUser, protocmp.Transform()); diff != "" {
				t.Errorf("r.LookupUserByEmail(%v, %v) user != test.wantUser (-got +want)\n%s", ctx, test.lookupEmail, diff)
			}
			if got, want := err, test.wantLookupErr; !cmp.Equal(got, want) {
				t.Errorf("r.LookupUserByEmail(%v, %v) err = %v; want %v", ctx, test.lookupEmail, got, want)
			}
		})
	}
}

func TestUsers_DeleteUser(t *testing.T) {
	ctx := context.Background()
	for _, test := range []struct {
		desc          string
		users         []*streckuv1.User
		passwords     []string
		name          string
		want          error
		lookupName    string
		wantUser      *streckuv1.User
		wantLookupErr error
	}{
		{
			desc:          "Empty",
			users:         nil,
			name:          userNames["foobar"],
			want:          &UserNotFoundError{Name: userNames["foobar"]},
			lookupName:    userNames["barbaz"],
			wantUser:      nil,
			wantLookupErr: &UserNotFoundError{Name: userNames["barbaz"]},
		},
		{
			desc: "OneUserOK",
			users: []*streckuv1.User{
				{Name: userNames["foobar"], EmailAddress: "foobar@example.com", DisplayName: "Foo Bar"},
			},
			passwords:     []string{"foobar"},
			name:          userNames["foobar"],
			want:          nil,
			lookupName:    userNames["foobar"],
			wantUser:      nil,
			wantLookupErr: &UserNotFoundError{Name: userNames["foobar"]},
		},
		{
			desc: "MultipleUsersLookupDeleted",
			users: []*streckuv1.User{
				{Name: userNames["foobar"], EmailAddress: "foobar@example.com", DisplayName: "Foo Bar"},
				{Name: userNames["barbaz"], EmailAddress: "barbaz@example.com", DisplayName: "Barba Z."},
				{Name: userNames["quux"], EmailAddress: "quux@example.com", DisplayName: "Qu Ux"},
			},
			passwords: []string{
				"foobar",
				"barbaz",
				"quux",
			},
			name:          userNames["barbaz"],
			want:          nil,
			lookupName:    userNames["barbaz"],
			wantUser:      nil,
			wantLookupErr: &UserNotFoundError{Name: userNames["barbaz"]},
		},
		{
			desc: "MultipleUsersLookupExisting",
			users: []*streckuv1.User{
				{Name: userNames["foobar"], EmailAddress: "foobar@example.com", DisplayName: "Foo Bar"},
				{Name: userNames["barbaz"], EmailAddress: "barbaz@example.com", DisplayName: "Barba Z."},
				{Name: userNames["quux"], EmailAddress: "quux@example.com", DisplayName: "Qu Ux"},
			},
			passwords: []string{
				"foobar",
				"barbaz",
				"quux",
			},
			name:          userNames["barbaz"],
			want:          nil,
			lookupName:    userNames["foobar"],
			wantUser:      &streckuv1.User{Name: userNames["foobar"], EmailAddress: "foobar@example.com", DisplayName: "Foo Bar"},
			wantLookupErr: nil,
		},
		{
			desc: "OneUserNotFound",
			users: []*streckuv1.User{
				{Name: userNames["foobar"], EmailAddress: "foobar@example.com", DisplayName: "Foo Bar"},
			},
			passwords:     []string{"foobar"},
			name:          userNames["barbaz"],
			want:          &UserNotFoundError{Name: userNames["barbaz"]},
			lookupName:    userNames["foobar"],
			wantUser:      &streckuv1.User{Name: userNames["foobar"], EmailAddress: "foobar@example.com", DisplayName: "Foo Bar"},
			wantLookupErr: nil,
		},
		{
			desc: "MultipleUsersNotFound",
			users: []*streckuv1.User{
				{Name: userNames["foobar"], EmailAddress: "foobar@example.com", DisplayName: "Foo Bar"},
				{Name: userNames["barbaz"], EmailAddress: "barbaz@example.com", DisplayName: "Barba Z."},
				{Name: userNames["quux"], EmailAddress: "quux@example.com", DisplayName: "Qu Ux"},
			},
			passwords: []string{
				"foobar",
				"barbaz",
				"quux",
			},
			name:          userNames["cookie"],
			want:          &UserNotFoundError{Name: userNames["cookie"]},
			lookupName:    userNames["foobar"],
			wantUser:      &streckuv1.User{Name: userNames["foobar"], EmailAddress: "foobar@example.com", DisplayName: "Foo Bar"},
			wantLookupErr: nil,
		},
	} {
		t.Run(test.desc, func(t *testing.T) {
			r := seed(t, test.users, test.passwords)
			err := r.DeleteUser(ctx, test.name)
			if got, want := err, test.want; !cmp.Equal(got, want) {
				t.Errorf("r.DeleteUser(%v, %q) = %v; want %v", ctx, test.name, got, want)
			}
			user, err := r.LookupUser(ctx, test.lookupName)
			if diff := cmp.Diff(user, test.wantUser, protocmp.Transform()); diff != "" {
				t.Errorf("r.LookupUser(%v, %q) user != test.wantUser (-got +want)\n%s", ctx, test.lookupName, diff)
			}
			if got, want := err, test.wantLookupErr; !cmp.Equal(got, want) {
				t.Errorf("r.LookupUser(%v, %q) err = %v; want %v", ctx, test.lookupName, got, want)
			}
		})
	}
}
