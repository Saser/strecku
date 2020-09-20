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
}

func userLess(u1, u2 *streckuv1.User) bool {
	return u1.Name < u2.Name
}

func seed(t *testing.T, users []*streckuv1.User) *Repository {
	t.Helper()
	mUsers := make(map[string]*streckuv1.User, len(users))
	mNames := make(map[string]string, len(users))
	for _, user := range users {
		if got := Validate(user); got != nil {
			t.Errorf("Validate(%v) = %v; want %v", user, got, nil)
		}
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

func TestUserExistsError_Error(t *testing.T) {
	err := &UserExistsError{EmailAddress: "user@example.com"}
	if got, want := err.Error(), `duplicate user email address: "user@example.com"`; got != want {
		t.Errorf("err.Error() = %q; want %q", got, want)
	}
}

func TestUsers_LookupUser(t *testing.T) {
	ctx := context.Background()
	for _, test := range []struct {
		desc     string
		users    []*streckuv1.User
		name     string
		wantUser *streckuv1.User
		wantErr  error
	}{
		{
			desc:     "EmptyDatabaseEmptyName",
			users:    nil,
			name:     "",
			wantUser: nil,
			wantErr:  &UserNotFoundError{Name: ""},
		},
		{
			desc:     "EmptyDatabaseNonEmptyName",
			users:    nil,
			name:     userNames["foobar"],
			wantUser: nil,
			wantErr:  &UserNotFoundError{Name: userNames["foobar"]},
		},
		{
			desc: "OneUserOK",
			users: []*streckuv1.User{
				{Name: userNames["foobar"], EmailAddress: "user@example.com", DisplayName: "User"},
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
			name:     userNames["barbaz"],
			wantUser: &streckuv1.User{Name: userNames["barbaz"], EmailAddress: "barbaz@example.com", DisplayName: "Barba Z."},
			wantErr:  nil,
		},
		{
			desc: "OneUserNotFound",
			users: []*streckuv1.User{
				{Name: userNames["foobar"], EmailAddress: "user@example.com", DisplayName: "User"},
			},
			name:     "users/notfoobar",
			wantUser: nil,
			wantErr:  &UserNotFoundError{Name: "users/notfoobar"},
		},
		{
			desc: "MultipleUsersNotFound",
			users: []*streckuv1.User{
				{Name: userNames["foobar"], EmailAddress: "foobar@example.com", DisplayName: "Foo Bar"},
				{Name: userNames["barbaz"], EmailAddress: "barbaz@example.com", DisplayName: "Barba Z."},
				{Name: userNames["quux"], EmailAddress: "quux@example.com", DisplayName: "Qu Ux"},
			},
			name:     "users/notfoobar",
			wantUser: nil,
			wantErr:  &UserNotFoundError{Name: "users/notfoobar"},
		},
	} {
		t.Run(test.desc, func(t *testing.T) {
			r := seed(t, test.users)
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
		emailAddress string
		wantUser     *streckuv1.User
		wantErr      error
	}{
		{
			desc:         "EmptyDatabaseEmptyName",
			users:        nil,
			emailAddress: "",
			wantUser:     nil,
			wantErr:      &UserNotFoundError{EmailAddress: ""},
		},
		{
			desc:         "EmptyDatabaseNonEmptyName",
			users:        nil,
			emailAddress: "user@example.com",
			wantUser:     nil,
			wantErr:      &UserNotFoundError{EmailAddress: "user@example.com"},
		},
		{
			desc: "OneUserOK",
			users: []*streckuv1.User{
				{Name: userNames["foobar"], EmailAddress: "user@example.com", DisplayName: "User"},
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
			emailAddress: "barbaz@example.com",
			wantUser:     &streckuv1.User{Name: userNames["barbaz"], EmailAddress: "barbaz@example.com", DisplayName: "Barba Z."},
			wantErr:      nil,
		},
		{
			desc: "OneUserNotFound",
			users: []*streckuv1.User{
				{Name: userNames["foobar"], EmailAddress: "user@example.com", DisplayName: "User"},
			},
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
			emailAddress: "notfoobar@example.com",
			wantUser:     nil,
			wantErr:      &UserNotFoundError{EmailAddress: "notfoobar@example.com"},
		},
	} {
		t.Run(test.desc, func(t *testing.T) {
			r := seed(t, test.users)
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
		name  string
		users []*streckuv1.User
	}{
		{name: "Empty", users: nil},
		{
			name: "OneUser",
			users: []*streckuv1.User{
				{Name: userNames["foobar"], EmailAddress: "user@example.com", DisplayName: "Foo Bar"},
			},
		},
		{
			name: "ThreeUsers",
			users: []*streckuv1.User{
				{Name: userNames["foobar"], EmailAddress: "foobar@example.com", DisplayName: "Foo Bar"},
				{Name: userNames["barbaz"], EmailAddress: "barbaz@example.com", DisplayName: "Barba Z."},
				{Name: userNames["quux"], EmailAddress: "quux@example.com", DisplayName: "Qu Ux"},
			},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			r := seed(t, test.users)
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
			predicate: func(user *streckuv1.User) bool { return false },
			want:      nil,
		},
		{
			name: "OneUserOneMatching",
			users: []*streckuv1.User{
				{Name: userNames["foobar"], EmailAddress: "user@example.com", DisplayName: "Foo Bar"},
			},
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
			predicate: func(user *streckuv1.User) bool { return strings.Contains(user.DisplayName, "Bar") },
			want: []*streckuv1.User{
				{Name: userNames["foobar"], EmailAddress: "foobar@example.com", DisplayName: "Foo Bar"},
				{Name: userNames["barbaz"], EmailAddress: "barbaz@example.com", DisplayName: "Barba Z."},
			},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			r := seed(t, test.users)
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
		name  string
		users []*streckuv1.User
		user  *streckuv1.User
		want  error
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
			user: &streckuv1.User{Name: userNames["barbaz"], EmailAddress: "barbaz@example.com", DisplayName: "Foo Bar"},
			want: nil,
		},
		{
			name: "MultipleUsersOK",
			users: []*streckuv1.User{
				{Name: userNames["foobar"], EmailAddress: "foobar@example.com", DisplayName: "Foo Bar"},
				{Name: userNames["barbaz"], EmailAddress: "barbaz@example.com", DisplayName: "Barba Z."},
				{Name: userNames["quux"], EmailAddress: "quux@example.com", DisplayName: "Qu Ux"},
			},
			user: &streckuv1.User{Name: "users/cookie", EmailAddress: "cookie@example.com", DisplayName: "Cookie"},
			want: nil,
		},
		{
			name: "OneUserDuplicateEmail",
			users: []*streckuv1.User{
				{Name: userNames["foobar"], EmailAddress: "foobar@example.com", DisplayName: "Foo Bar"},
			},
			user: &streckuv1.User{Name: userNames["barbaz"], EmailAddress: "foobar@example.com", DisplayName: "Foo Bar"},
			want: &UserExistsError{EmailAddress: "foobar@example.com"},
		},
		{
			name: "MultipleUsersDuplicateEmail",
			users: []*streckuv1.User{
				{Name: userNames["foobar"], EmailAddress: "foobar@example.com", DisplayName: "Foo Bar"},
				{Name: userNames["barbaz"], EmailAddress: "barbaz@example.com", DisplayName: "Barba Z."},
				{Name: userNames["quux"], EmailAddress: "quux@example.com", DisplayName: "Qu Ux"},
			},
			user: &streckuv1.User{Name: "users/cookie", EmailAddress: "foobar@example.com", DisplayName: "Cookie"},
			want: &UserExistsError{EmailAddress: "foobar@example.com"},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			r := seed(t, test.users)
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
			updated:       &streckuv1.User{Name: "users/notfound", EmailAddress: "new-foobar@example.com", DisplayName: "Foo Bar"},
			wantUpdateErr: &UserNotFoundError{Name: "users/notfound"},
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
			updated:       &streckuv1.User{Name: "users/notfound", EmailAddress: "new-barbaz@example.com", DisplayName: "Barba Z."},
			wantUpdateErr: &UserNotFoundError{Name: "users/notfound"},
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
			updated:       &streckuv1.User{Name: userNames["barbaz"], EmailAddress: "foobar@example.com", DisplayName: "All New Barba Z."},
			wantUpdateErr: &UserExistsError{EmailAddress: "foobar@example.com"},
			lookupEmail:   "barbaz@example.com",
			wantUser:      &streckuv1.User{Name: userNames["barbaz"], EmailAddress: "barbaz@example.com", DisplayName: "Barba Z."},
			wantLookupErr: nil,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			r := seed(t, test.users)
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
		name          string
		want          error
		lookupName    string
		wantUser      *streckuv1.User
		wantLookupErr error
	}{
		{
			desc:          "Empty",
			users:         nil,
			name:          "users/notfound",
			want:          &UserNotFoundError{Name: "users/notfound"},
			lookupName:    "users/alsonotfound",
			wantUser:      nil,
			wantLookupErr: &UserNotFoundError{Name: "users/alsonotfound"},
		},
		{
			desc: "OneUserOK",
			users: []*streckuv1.User{
				{Name: userNames["foobar"], EmailAddress: "foobar@example.com", DisplayName: "Foo Bar"},
			},
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
			name:          "users/notfound",
			want:          &UserNotFoundError{Name: "users/notfound"},
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
			name:          "users/notfound",
			want:          &UserNotFoundError{Name: "users/notfound"},
			lookupName:    userNames["foobar"],
			wantUser:      &streckuv1.User{Name: userNames["foobar"], EmailAddress: "foobar@example.com", DisplayName: "Foo Bar"},
			wantLookupErr: nil,
		},
	} {
		t.Run(test.desc, func(t *testing.T) {
			r := seed(t, test.users)
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
