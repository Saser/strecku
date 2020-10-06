package users

import (
	"context"
	"fmt"
	"testing"

	pb "github.com/Saser/strecku/api/v1"
	"github.com/Saser/strecku/resources/users/testusers"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"google.golang.org/protobuf/testing/protocmp"
)

func userLess(u1, u2 *pb.User) bool {
	return u1.Name < u2.Name
}

func seedAlice(t *testing.T) *Repository {
	return SeedRepository(t, []*pb.User{testusers.Alice}, []string{testusers.AlicePassword})
}

func seedAliceBob(t *testing.T) *Repository {
	return SeedRepository(
		t,
		[]*pb.User{testusers.Alice, testusers.Bob},
		[]string{testusers.AlicePassword, testusers.BobPassword},
	)
}

func seedAliceBobCarol(t *testing.T) *Repository {
	return SeedRepository(
		t,
		[]*pb.User{testusers.Alice, testusers.Bob, testusers.Carol},
		[]string{testusers.AlicePassword, testusers.BobPassword, testusers.CarolPassword},
	)
}

func TestUserNotFoundError_Error(t *testing.T) {
	for _, test := range []struct {
		name         string
		emailAddress string
		want         string
	}{
		{name: testusers.Alice.Name, want: fmt.Sprintf("user not found: %q", testusers.Alice.Name)},
		{emailAddress: testusers.Alice.EmailAddress, want: fmt.Sprintf("user email not found: %q", testusers.Alice.EmailAddress)},
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
			err:    &UserNotFoundError{Name: testusers.Alice.Name},
			target: &UserNotFoundError{Name: testusers.Alice.Name},
			want:   true,
		},
		{
			err:    &UserNotFoundError{Name: testusers.Alice.Name},
			target: &UserNotFoundError{Name: testusers.Bob.Name},
			want:   false,
		},
		{
			err:    &UserNotFoundError{Name: testusers.Alice.Name},
			target: &UserNotFoundError{EmailAddress: testusers.Alice.EmailAddress},
			want:   false,
		},
		{
			err:    &UserNotFoundError{EmailAddress: testusers.Alice.EmailAddress},
			target: &UserNotFoundError{Name: testusers.Alice.Name},
			want:   false,
		},
		{
			err:    &UserNotFoundError{EmailAddress: testusers.Alice.EmailAddress},
			target: &UserNotFoundError{EmailAddress: testusers.Alice.EmailAddress},
			want:   true,
		},
		{
			err:    &UserNotFoundError{EmailAddress: testusers.Alice.EmailAddress},
			target: &UserNotFoundError{EmailAddress: testusers.Bob.EmailAddress},
			want:   false,
		},
		{
			err:    &UserNotFoundError{Name: testusers.Alice.Name},
			target: fmt.Errorf("user not found: %q", testusers.Alice.Name),
			want:   false,
		},
		{
			err:    &UserNotFoundError{EmailAddress: testusers.Alice.EmailAddress},
			target: fmt.Errorf("user email not found: %q", testusers.Alice.EmailAddress),
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
		{name: testusers.Alice.Name, want: fmt.Sprintf("user exists: %q", testusers.Alice.Name)},
		{emailAddress: testusers.Alice.EmailAddress, want: fmt.Sprintf("user email exists: %q", testusers.Alice.EmailAddress)},
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
			err:    &UserExistsError{Name: testusers.Alice.Name},
			target: &UserExistsError{Name: testusers.Alice.Name},
			want:   true,
		},
		{
			err:    &UserExistsError{Name: testusers.Alice.Name},
			target: &UserExistsError{Name: testusers.Bob.Name},
			want:   false,
		},
		{
			err:    &UserExistsError{Name: testusers.Alice.Name},
			target: &UserExistsError{EmailAddress: testusers.Alice.EmailAddress},
			want:   false,
		},
		{
			err:    &UserExistsError{EmailAddress: testusers.Alice.EmailAddress},
			target: &UserExistsError{Name: testusers.Alice.Name},
			want:   false,
		},
		{
			err:    &UserExistsError{EmailAddress: testusers.Alice.EmailAddress},
			target: &UserExistsError{EmailAddress: testusers.Alice.EmailAddress},
			want:   true,
		},
		{
			err:    &UserExistsError{EmailAddress: testusers.Alice.EmailAddress},
			target: &UserExistsError{EmailAddress: testusers.Bob.EmailAddress},
			want:   false,
		},
		{
			err:    &UserExistsError{Name: testusers.Alice.Name},
			target: fmt.Errorf("user exists: %q", testusers.Alice.Name),
			want:   false,
		},
		{
			err:    &UserExistsError{EmailAddress: testusers.Alice.EmailAddress},
			target: fmt.Errorf("user email exists: %q", testusers.Alice.EmailAddress),
			want:   false,
		},
	} {
		if got := test.err.Is(test.target); got != test.want {
			t.Errorf("test.err.Is(%v) = %v; want %v", test.target, got, test.want)
		}
	}
}

func TestWrongPasswordError_Error(t *testing.T) {
	err := &WrongPasswordError{Name: testusers.Alice.Name}
	if got, want := err.Error(), fmt.Sprintf("wrong password for user %q", testusers.Alice.Name); got != want {
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
			err:    &WrongPasswordError{Name: testusers.Alice.Name},
			target: &WrongPasswordError{Name: testusers.Alice.Name},
			want:   true,
		},
		{
			err:    &WrongPasswordError{Name: testusers.Alice.Name},
			target: &WrongPasswordError{Name: testusers.Bob.Name},
			want:   false,
		},
		{
			err:    &WrongPasswordError{Name: testusers.Alice.Name},
			target: fmt.Errorf("wrong password for user %q", testusers.Alice.Name),
			want:   false,
		},
	} {
		if got := test.err.Is(test.target); got != test.want {
			t.Errorf("test.err.Is(%v) = %v; want %v", test.target, got, test.want)
		}
	}
}

func TestRepository_Authenticate(t *testing.T) {
	ctx := context.Background()
	r := seedAlice(t)
	for _, test := range []struct {
		desc     string
		name     string
		password string
		want     error
	}{
		{
			desc:     "OK",
			name:     testusers.Alice.Name,
			password: testusers.AlicePassword,
			want:     nil,
		},
		{
			desc:     "NotFound",
			name:     testusers.Bob.Name,
			password: testusers.BobPassword,
			want:     &UserNotFoundError{Name: testusers.Bob.Name},
		},
		{
			desc:     "WrongPassword",
			name:     testusers.Alice.Name,
			password: "wrong password",
			want:     &WrongPasswordError{Name: testusers.Alice.Name},
		},
	} {
		t.Run(test.desc, func(t *testing.T) {
			if got := r.Authenticate(ctx, test.name, test.password); !cmp.Equal(got, test.want) {
				t.Errorf("r.Authenticate(%v, %q, %q) = %v; want %v", ctx, test.name, test.password, got, test.want)
			}
		})
	}
}

func TestRepository_LookupUser(t *testing.T) {
	ctx := context.Background()
	r := seedAlice(t)
	for _, test := range []struct {
		desc     string
		name     string
		wantUser *pb.User
		wantErr  error
	}{
		{
			desc:     "OK",
			name:     testusers.Alice.Name,
			wantUser: testusers.Alice,
			wantErr:  nil,
		},
		{
			desc:     "EmptyName",
			name:     "",
			wantUser: nil,
			wantErr:  ErrNameEmpty,
		},
		{
			desc:     "NotFound",
			name:     testusers.Bob.Name,
			wantUser: nil,
			wantErr:  &UserNotFoundError{Name: testusers.Bob.Name},
		},
	} {
		t.Run(test.desc, func(t *testing.T) {
			user, err := r.LookupUser(ctx, test.name)
			if diff := cmp.Diff(user, test.wantUser, protocmp.Transform()); diff != "" {
				t.Errorf("r.LookupUser(%v, %q) user != test.wantUser (-got +want)\n%s", ctx, test.name, diff)
			}
			if got, want := err, test.wantErr; !cmp.Equal(got, want, cmpopts.EquateErrors()) {
				t.Errorf("r.LookupUser(%v, %q) err = %v; want %v", ctx, test.name, got, want)
			}
		})
	}
}

func TestRepository_LookupUserByEmail(t *testing.T) {
	ctx := context.Background()
	r := seedAlice(t)
	for _, test := range []struct {
		desc         string
		emailAddress string
		wantUser     *pb.User
		wantErr      error
	}{
		{
			desc:         "OK",
			emailAddress: testusers.Alice.EmailAddress,
			wantUser:     testusers.Alice,
			wantErr:      nil,
		},
		{
			desc:         "EmptyEmailAddress",
			emailAddress: "",
			wantUser:     nil,
			wantErr:      &UserNotFoundError{EmailAddress: ""},
		},
		{
			desc:         "NotFound",
			emailAddress: testusers.Bob.EmailAddress,
			wantUser:     nil,
			wantErr:      &UserNotFoundError{EmailAddress: testusers.Bob.EmailAddress},
		},
	} {
		t.Run(test.desc, func(t *testing.T) {
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

func TestRepository_ListUsers(t *testing.T) {
	ctx := context.Background()
	allUsers := []*pb.User{
		testusers.Alice,
		testusers.Bob,
		testusers.Carol,
	}
	r := seedAliceBobCarol(t)
	users, err := r.ListUsers(ctx)
	if diff := cmp.Diff(
		users, allUsers, protocmp.Transform(),
		cmpopts.SortSlices(userLess),
	); diff != "" {
		t.Errorf("r.ListUsers(%v) users != allUsers (-got +want)\n%s", ctx, diff)
	}
	if err != nil {
		t.Errorf("r.ListUsers(%v) err = %v; want nil", ctx, err)
	}
}

func TestRepository_FilterUsers(t *testing.T) {
	ctx := context.Background()
	r := seedAliceBobCarol(t)
	for _, test := range []struct {
		name      string
		predicate func(*pb.User) bool
		want      []*pb.User
	}{
		{
			name:      "NoneMatching",
			predicate: func(user *pb.User) bool { return false },
			want:      nil,
		},
		{
			name:      "OneMatching",
			predicate: func(user *pb.User) bool { return user.DisplayName == "Alice" },
			want: []*pb.User{
				testusers.Alice,
			},
		},
		{
			name:      "SeveralMatching",
			predicate: func(user *pb.User) bool { return user.DisplayName == "Alice" || user.DisplayName == "Bob" },
			want: []*pb.User{
				testusers.Alice,
				testusers.Bob,
			},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
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

func TestRepository_CreateUser(t *testing.T) {
	ctx := context.Background()
	// The Repository will be seeded with testusers.Alice.
	// However, since the repository is possibly mutated, it needs to be seeded
	// for each test case.
	for _, test := range []struct {
		name     string
		user     *pb.User
		password string
		want     error
	}{
		{
			name:     "OK",
			user:     testusers.Bob,
			password: testusers.BobPassword,
			want:     nil,
		},
		{
			name:     "EmptyPassword",
			user:     testusers.Bob,
			password: "",
			want:     ErrEmptyPassword,
		},
		{
			name:     "DuplicateEmail",
			user:     &pb.User{Name: testusers.Bob.Name, EmailAddress: testusers.Alice.EmailAddress, DisplayName: testusers.Bob.DisplayName},
			password: testusers.BobPassword,
			want:     &UserExistsError{EmailAddress: testusers.Alice.EmailAddress},
		},
		{
			name:     "DuplicateName",
			user:     &pb.User{Name: testusers.Alice.Name, EmailAddress: testusers.Bob.EmailAddress, DisplayName: testusers.Bob.DisplayName},
			password: testusers.BobPassword,
			want:     &UserExistsError{Name: testusers.Alice.Name},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			r := seedAlice(t)
			if got := r.CreateUser(ctx, test.user, test.password); !cmp.Equal(got, test.want, cmpopts.EquateErrors()) {
				t.Errorf("r.CreateUser(%v, %v) = %v; want %v", ctx, test.user, got, test.want)
			}
		})
	}
}

func TestRepository_UpdateUser(t *testing.T) {
	ctx := context.Background()
	type testCase struct {
		desc        string
		name, email string
		wantUser    *pb.User
		wantErr     error
	}
	testF := func(t *testing.T, r *Repository, lookups []testCase) {
		for _, test := range lookups {
			t.Run(test.desc, func(t *testing.T) {
				if (test.name == "" && test.email == "") || (test.name != "" && test.email != "") {
					t.Fatalf("test.name, test.email = %q, %q; want exactly one to be non-empty", test.name, test.email)
				}
				var (
					user   *pb.User
					err    error
					lookup string
				)
				switch {
				case test.name != "":
					user, err = r.LookupUser(ctx, test.name)
					lookup = fmt.Sprintf("r.LookupUser(%v, %q)", ctx, test.name)
				case test.email != "":
					user, err = r.LookupUserByEmail(ctx, test.email)
					lookup = fmt.Sprintf("r.LookupUserByEmail(%v, %q)", ctx, test.email)
				}
				if diff := cmp.Diff(user, test.wantUser, protocmp.Transform()); diff != "" {
					t.Errorf("%s user != test.wantUser (-got +want)\n%s", lookup, diff)
				}
				if !cmp.Equal(err, test.wantErr, cmpopts.EquateErrors()) {
					t.Errorf("%s err = %v; want %v", lookup, err, test.wantErr)
				}
			})
		}
	}

	// Test scenarios where the update is successful.
	t.Run("OK", func(t *testing.T) {
		r := seedAlice(t)
		oldAlice := Clone(testusers.Alice)
		newAlice := Clone(oldAlice)
		newAlice.EmailAddress = "new-alice@example.com"
		newAlice.DisplayName = "New Alice"
		if err := r.UpdateUser(ctx, newAlice); err != nil {
			t.Errorf("r.UpdateUser(%v, %v) = %v; want nil", ctx, newAlice, err)
		}
		testCases := []testCase{
			{
				desc:     "ByName",
				name:     testusers.Alice.Name,
				email:    "",
				wantUser: newAlice,
				wantErr:  nil,
			},
			{
				desc:     "ByOldEmail",
				name:     "",
				email:    oldAlice.EmailAddress,
				wantUser: nil,
				wantErr:  &UserNotFoundError{EmailAddress: oldAlice.EmailAddress},
			},
			{
				desc:     "ByNewEmail",
				name:     "",
				email:    newAlice.EmailAddress,
				wantUser: newAlice,
				wantErr:  nil,
			},
		}
		testF(t, r, testCases)
	})

	// Test scenarios where the update fails due to trying to update to an existing email address.
	t.Run("ExistingEmailAddress", func(t *testing.T) {
		r := seedAliceBob(t)
		oldAlice := Clone(testusers.Alice)
		newAlice := Clone(oldAlice)
		newAlice.EmailAddress = testusers.Bob.EmailAddress
		want := &UserExistsError{EmailAddress: testusers.Bob.EmailAddress}
		if got := r.UpdateUser(ctx, newAlice); !cmp.Equal(got, want, cmpopts.EquateErrors()) {
			t.Errorf("r.UpdateUser(%v, %v) = %v; want %v", ctx, newAlice, got, want)
		}
		testCases := []testCase{
			{
				desc:     "ByName_Alice",
				name:     oldAlice.Name,
				email:    "",
				wantUser: oldAlice,
				wantErr:  nil,
			},
			{
				desc:     "ByName_Bob",
				name:     testusers.Bob.Name,
				email:    "",
				wantUser: testusers.Bob,
				wantErr:  nil,
			},
			{
				desc:     "ByEmail_OldAlice",
				name:     "",
				email:    oldAlice.EmailAddress,
				wantUser: oldAlice,
				wantErr:  nil,
			},
			{
				desc:     "ByEmail_NewAlice",
				name:     "",
				email:    newAlice.EmailAddress,
				wantUser: testusers.Bob,
				wantErr:  nil,
			},
			{
				desc:     "ByEmail_Bob",
				name:     "",
				email:    testusers.Bob.EmailAddress,
				wantUser: testusers.Bob,
				wantErr:  nil,
			},
		}
		testF(t, r, testCases)
	})

	// Test scenarios where the update fails for other reasons.
	t.Run("Errors", func(t *testing.T) {
		r := seedAlice(t)
		for _, test := range []struct {
			desc   string
			modify func(alice *pb.User)
			want   error
		}{
			{
				desc:   "EmptyName",
				modify: func(alice *pb.User) { alice.Name = "" },
				want:   ErrNameEmpty,
			},
			{
				desc:   "EmptyEmailAddress",
				modify: func(alice *pb.User) { alice.EmailAddress = "" },
				want:   ErrEmailAddressEmpty,
			},
			{
				desc:   "EmptyDisplayName",
				modify: func(alice *pb.User) { alice.DisplayName = "" },
				want:   ErrDisplayNameEmpty,
			},
			{
				desc:   "NotFound",
				modify: func(alice *pb.User) { alice.Name = testusers.Bob.Name },
				want:   &UserNotFoundError{Name: testusers.Bob.Name},
			},
		} {
			t.Run(test.desc, func(t *testing.T) {
				updated := Clone(testusers.Alice)
				test.modify(updated)
				if got := r.UpdateUser(ctx, updated); !cmp.Equal(got, test.want, cmpopts.EquateErrors()) {
					t.Errorf("r.UpdateUser(%v, %v) = %v; want %v", ctx, updated, got, test.want)
				}
			})
		}
	})
}

func TestRepository_DeleteUser(t *testing.T) {
	ctx := context.Background()

	// Test scenarios where the delete succeeded.
	t.Run("OK", func(t *testing.T) {
		r := seedAliceBob(t)
		if err := r.DeleteUser(ctx, testusers.Alice.Name); err != nil {
			t.Fatalf("r.DeleteUser(%v, %q) = %v; want nil", ctx, testusers.Alice.Name, err)
		}
		// Test the effects of the delete by looking up users.
		for _, test := range []struct {
			desc     string
			name     string
			wantUser *pb.User
			wantErr  error
		}{
			{
				desc:     "LookupDeleted",
				name:     testusers.Alice.Name,
				wantUser: nil,
				wantErr:  &UserNotFoundError{Name: testusers.Alice.Name},
			},
			{
				desc:     "LookupExisting",
				name:     testusers.Bob.Name,
				wantUser: testusers.Bob,
				wantErr:  nil,
			},
		} {
			t.Run(test.desc, func(t *testing.T) {
				user, err := r.LookupUser(ctx, test.name)
				if diff := cmp.Diff(user, test.wantUser, protocmp.Transform()); diff != "" {
					t.Errorf("r.LookupUser(%v, %q) user != test.wantUser (-got +want)\n%s", ctx, test.name, diff)
				}
				if got, want := err, test.wantErr; !cmp.Equal(got, want) {
					t.Errorf("r.LookupUser(%v, %q) err = %v; want %v", ctx, test.name, got, want)
				}
			})
		}
	})

	// Test scenarios where the delete should not succeed.
	t.Run("Errors", func(t *testing.T) {
		// The repository will be seeded with only Alice.
		for _, test := range []struct {
			desc string
			name string
			want error
		}{
			{
				desc: "EmptyName",
				name: "",
				want: ErrNameEmpty,
			},
			{
				desc: "NotFound",
				name: testusers.Bob.Name,
				want: &UserNotFoundError{Name: testusers.Bob.Name},
			},
		} {
			t.Run(test.desc, func(t *testing.T) {
				r := seedAlice(t)
				if got := r.DeleteUser(ctx, test.name); !cmp.Equal(got, test.want, cmpopts.EquateErrors()) {
					t.Errorf("r.DeleteUser(%v, %q) = %v; want %v", ctx, test.name, got, test.want)
				}
			})
		}
	})
}
