package users

import (
	"context"
	"fmt"
	"testing"

	pb "github.com/Saser/strecku/api/v1"
	"github.com/Saser/strecku/resources/testresources"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"google.golang.org/protobuf/testing/protocmp"
)

func userLess(u1, u2 *pb.User) bool {
	return u1.Name < u2.Name
}

func seedAlice(t *testing.T) *Repository {
	return SeedRepository(t, []*pb.User{testresources.Alice}, []string{testresources.AlicePassword})
}

func seedAliceBob(t *testing.T) *Repository {
	return SeedRepository(
		t,
		[]*pb.User{testresources.Alice, testresources.Bob},
		[]string{testresources.AlicePassword, testresources.BobPassword},
	)
}

func seedAliceBobCarol(t *testing.T) *Repository {
	return SeedRepository(
		t,
		[]*pb.User{testresources.Alice, testresources.Bob, testresources.Carol},
		[]string{testresources.AlicePassword, testresources.BobPassword, testresources.CarolPassword},
	)
}

func TestNotFoundError_Error(t *testing.T) {
	for _, test := range []struct {
		name         string
		emailAddress string
		want         string
	}{
		{name: testresources.Alice.Name, want: fmt.Sprintf("user not found: %q", testresources.Alice.Name)},
		{emailAddress: testresources.Alice.EmailAddress, want: fmt.Sprintf("user email not found: %q", testresources.Alice.EmailAddress)},
	} {
		err := &NotFoundError{
			Name:         test.name,
			EmailAddress: test.emailAddress,
		}
		if got := err.Error(); got != test.want {
			t.Errorf("err.Error() = %q; want %q", got, test.want)
		}
	}
}

func TestNotFoundError_Is(t *testing.T) {
	for _, test := range []struct {
		err    *NotFoundError
		target error
		want   bool
	}{
		{
			err:    &NotFoundError{Name: testresources.Alice.Name},
			target: &NotFoundError{Name: testresources.Alice.Name},
			want:   true,
		},
		{
			err:    &NotFoundError{Name: testresources.Alice.Name},
			target: &NotFoundError{Name: testresources.Bob.Name},
			want:   false,
		},
		{
			err:    &NotFoundError{Name: testresources.Alice.Name},
			target: &NotFoundError{EmailAddress: testresources.Alice.EmailAddress},
			want:   false,
		},
		{
			err:    &NotFoundError{EmailAddress: testresources.Alice.EmailAddress},
			target: &NotFoundError{Name: testresources.Alice.Name},
			want:   false,
		},
		{
			err:    &NotFoundError{EmailAddress: testresources.Alice.EmailAddress},
			target: &NotFoundError{EmailAddress: testresources.Alice.EmailAddress},
			want:   true,
		},
		{
			err:    &NotFoundError{EmailAddress: testresources.Alice.EmailAddress},
			target: &NotFoundError{EmailAddress: testresources.Bob.EmailAddress},
			want:   false,
		},
		{
			err:    &NotFoundError{Name: testresources.Alice.Name},
			target: fmt.Errorf("user not found: %q", testresources.Alice.Name),
			want:   false,
		},
		{
			err:    &NotFoundError{EmailAddress: testresources.Alice.EmailAddress},
			target: fmt.Errorf("user email not found: %q", testresources.Alice.EmailAddress),
			want:   false,
		},
	} {
		if got := test.err.Is(test.target); got != test.want {
			t.Errorf("test.err.Is(%v) = %v; want %v", test.target, got, test.want)
		}
	}
}

func TestExistsError_Error(t *testing.T) {
	for _, test := range []struct {
		name         string
		emailAddress string
		want         string
	}{
		{name: testresources.Alice.Name, want: fmt.Sprintf("user exists: %q", testresources.Alice.Name)},
		{emailAddress: testresources.Alice.EmailAddress, want: fmt.Sprintf("user email exists: %q", testresources.Alice.EmailAddress)},
	} {
		err := &ExistsError{
			Name:         test.name,
			EmailAddress: test.emailAddress,
		}
		if got := err.Error(); got != test.want {
			t.Errorf("err.Error() = %q; want %q", got, test.want)
		}
	}
}

func TestExistsError_Is(t *testing.T) {
	for _, test := range []struct {
		err    *ExistsError
		target error
		want   bool
	}{
		{
			err:    &ExistsError{Name: testresources.Alice.Name},
			target: &ExistsError{Name: testresources.Alice.Name},
			want:   true,
		},
		{
			err:    &ExistsError{Name: testresources.Alice.Name},
			target: &ExistsError{Name: testresources.Bob.Name},
			want:   false,
		},
		{
			err:    &ExistsError{Name: testresources.Alice.Name},
			target: &ExistsError{EmailAddress: testresources.Alice.EmailAddress},
			want:   false,
		},
		{
			err:    &ExistsError{EmailAddress: testresources.Alice.EmailAddress},
			target: &ExistsError{Name: testresources.Alice.Name},
			want:   false,
		},
		{
			err:    &ExistsError{EmailAddress: testresources.Alice.EmailAddress},
			target: &ExistsError{EmailAddress: testresources.Alice.EmailAddress},
			want:   true,
		},
		{
			err:    &ExistsError{EmailAddress: testresources.Alice.EmailAddress},
			target: &ExistsError{EmailAddress: testresources.Bob.EmailAddress},
			want:   false,
		},
		{
			err:    &ExistsError{Name: testresources.Alice.Name},
			target: fmt.Errorf("user exists: %q", testresources.Alice.Name),
			want:   false,
		},
		{
			err:    &ExistsError{EmailAddress: testresources.Alice.EmailAddress},
			target: fmt.Errorf("user email exists: %q", testresources.Alice.EmailAddress),
			want:   false,
		},
	} {
		if got := test.err.Is(test.target); got != test.want {
			t.Errorf("test.err.Is(%v) = %v; want %v", test.target, got, test.want)
		}
	}
}

func TestWrongPasswordError_Error(t *testing.T) {
	err := &WrongPasswordError{Name: testresources.Alice.Name}
	if got, want := err.Error(), fmt.Sprintf("wrong password for user %q", testresources.Alice.Name); got != want {
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
			err:    &WrongPasswordError{Name: testresources.Alice.Name},
			target: &WrongPasswordError{Name: testresources.Alice.Name},
			want:   true,
		},
		{
			err:    &WrongPasswordError{Name: testresources.Alice.Name},
			target: &WrongPasswordError{Name: testresources.Bob.Name},
			want:   false,
		},
		{
			err:    &WrongPasswordError{Name: testresources.Alice.Name},
			target: fmt.Errorf("wrong password for user %q", testresources.Alice.Name),
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
			name:     testresources.Alice.Name,
			password: testresources.AlicePassword,
			want:     nil,
		},
		{
			desc:     "NotFound",
			name:     testresources.Bob.Name,
			password: testresources.BobPassword,
			want:     &NotFoundError{Name: testresources.Bob.Name},
		},
		{
			desc:     "WrongPassword",
			name:     testresources.Alice.Name,
			password: "wrong password",
			want:     &WrongPasswordError{Name: testresources.Alice.Name},
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
			name:     testresources.Alice.Name,
			wantUser: testresources.Alice,
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
			name:     testresources.Bob.Name,
			wantUser: nil,
			wantErr:  &NotFoundError{Name: testresources.Bob.Name},
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
			emailAddress: testresources.Alice.EmailAddress,
			wantUser:     testresources.Alice,
			wantErr:      nil,
		},
		{
			desc:         "EmptyEmailAddress",
			emailAddress: "",
			wantUser:     nil,
			wantErr:      &NotFoundError{EmailAddress: ""},
		},
		{
			desc:         "NotFound",
			emailAddress: testresources.Bob.EmailAddress,
			wantUser:     nil,
			wantErr:      &NotFoundError{EmailAddress: testresources.Bob.EmailAddress},
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
		testresources.Alice,
		testresources.Bob,
		testresources.Carol,
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
				testresources.Alice,
			},
		},
		{
			name:      "SeveralMatching",
			predicate: func(user *pb.User) bool { return user.DisplayName == "Alice" || user.DisplayName == "Bob" },
			want: []*pb.User{
				testresources.Alice,
				testresources.Bob,
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
			user:     testresources.Bob,
			password: testresources.BobPassword,
			want:     nil,
		},
		{
			name:     "EmptyPassword",
			user:     testresources.Bob,
			password: "",
			want:     ErrEmptyPassword,
		},
		{
			name:     "DuplicateEmail",
			user:     &pb.User{Name: testresources.Bob.Name, EmailAddress: testresources.Alice.EmailAddress, DisplayName: testresources.Bob.DisplayName},
			password: testresources.BobPassword,
			want:     &ExistsError{EmailAddress: testresources.Alice.EmailAddress},
		},
		{
			name:     "DuplicateName",
			user:     &pb.User{Name: testresources.Alice.Name, EmailAddress: testresources.Bob.EmailAddress, DisplayName: testresources.Bob.DisplayName},
			password: testresources.BobPassword,
			want:     &ExistsError{Name: testresources.Alice.Name},
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
		oldAlice := Clone(testresources.Alice)
		newAlice := Clone(oldAlice)
		newAlice.EmailAddress = "new-alice@example.com"
		newAlice.DisplayName = "New Alice"
		if err := r.UpdateUser(ctx, newAlice); err != nil {
			t.Errorf("r.UpdateUser(%v, %v) = %v; want nil", ctx, newAlice, err)
		}
		testCases := []testCase{
			{
				desc:     "ByName",
				name:     testresources.Alice.Name,
				email:    "",
				wantUser: newAlice,
				wantErr:  nil,
			},
			{
				desc:     "ByOldEmail",
				name:     "",
				email:    oldAlice.EmailAddress,
				wantUser: nil,
				wantErr:  &NotFoundError{EmailAddress: oldAlice.EmailAddress},
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
		oldAlice := Clone(testresources.Alice)
		newAlice := Clone(oldAlice)
		newAlice.EmailAddress = testresources.Bob.EmailAddress
		want := &ExistsError{EmailAddress: testresources.Bob.EmailAddress}
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
				name:     testresources.Bob.Name,
				email:    "",
				wantUser: testresources.Bob,
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
				wantUser: testresources.Bob,
				wantErr:  nil,
			},
			{
				desc:     "ByEmail_Bob",
				name:     "",
				email:    testresources.Bob.EmailAddress,
				wantUser: testresources.Bob,
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
				modify: func(alice *pb.User) { alice.Name = testresources.Bob.Name },
				want:   &NotFoundError{Name: testresources.Bob.Name},
			},
		} {
			t.Run(test.desc, func(t *testing.T) {
				updated := Clone(testresources.Alice)
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
		if err := r.DeleteUser(ctx, testresources.Alice.Name); err != nil {
			t.Fatalf("r.DeleteUser(%v, %q) = %v; want nil", ctx, testresources.Alice.Name, err)
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
				name:     testresources.Alice.Name,
				wantUser: nil,
				wantErr:  &NotFoundError{Name: testresources.Alice.Name},
			},
			{
				desc:     "LookupExisting",
				name:     testresources.Bob.Name,
				wantUser: testresources.Bob,
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
				name: testresources.Bob.Name,
				want: &NotFoundError{Name: testresources.Bob.Name},
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
