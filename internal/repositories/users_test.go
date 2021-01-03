package repositories

import (
	"context"
	"testing"

	pb "github.com/Saser/strecku/api/v1"
	"github.com/Saser/strecku/resourcename"
	"github.com/Saser/strecku/resources/testresources"
	"github.com/Saser/strecku/resources/users"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/stretchr/testify/suite"
	"google.golang.org/protobuf/testing/protocmp"
)

func userLess(u1, u2 *pb.User) bool {
	return u1.Name < u2.Name
}

type UsersTestSuite struct {
	suite.Suite
	r Users
}

func NewUsersTestSuite(r Users) *UsersTestSuite {
	return &UsersTestSuite{
		r: r,
	}
}

func (s *UsersTestSuite) seedUsers(ctx context.Context, t *testing.T, users []*pb.User, passwords []string) Users {
	t.Helper()
	SeedUsers(ctx, t, s.r, users, passwords)
	return s.r
}

func (s *UsersTestSuite) seedAlice(ctx context.Context, t *testing.T) Users {
	t.Helper()
	return s.seedUsers(
		ctx,
		t,
		[]*pb.User{
			testresources.Alice,
		},
		[]string{
			testresources.AlicePassword,
		},
	)
}

func (s *UsersTestSuite) seedAliceBob(ctx context.Context, t *testing.T) Users {
	t.Helper()
	return s.seedUsers(
		ctx,
		t,
		[]*pb.User{
			testresources.Alice,
			testresources.Bob,
		},
		[]string{
			testresources.AlicePassword,
			testresources.BobPassword,
		},
	)
}

func (s *UsersTestSuite) seedAliceBobCarol(ctx context.Context, t *testing.T) Users {
	t.Helper()
	return s.seedUsers(
		ctx,
		t,
		[]*pb.User{
			testresources.Alice,
			testresources.Bob,
			testresources.Carol,
		},
		[]string{
			testresources.AlicePassword,
			testresources.BobPassword,
			testresources.CarolPassword,
		},
	)
}

func (s *UsersTestSuite) TestAuthenticate() {
	t := s.T()
	ctx := context.Background()
	r := s.seedAlice(ctx, t)
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
			want:     &NotFound{Name: testresources.Bob.Name},
		},
		{
			desc:     "WrongPassword",
			name:     testresources.Alice.Name,
			password: "wrong password",
			want:     ErrUnauthenticated,
		},
	} {
		t.Run(test.desc, func(t *testing.T) {
			if got := r.Authenticate(ctx, test.name, test.password); !cmp.Equal(got, test.want, cmpopts.EquateErrors()) {
				t.Errorf("r.Authenticate(%v, %q, %q) = %v; want %v", ctx, test.name, test.password, got, test.want)
			}
		})
	}
}

func (s *UsersTestSuite) TestLookup() {
	t := s.T()
	ctx := context.Background()
	r := s.seedAlice(ctx, t)
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
			wantErr:  resourcename.ErrInvalidName,
		},
		{
			desc:     "NotFound",
			name:     testresources.Bob.Name,
			wantUser: nil,
			wantErr:  &NotFound{Name: testresources.Bob.Name},
		},
	} {
		t.Run(test.desc, func(t *testing.T) {
			user, err := r.Lookup(ctx, test.name)
			if diff := cmp.Diff(user, test.wantUser, protocmp.Transform()); diff != "" {
				t.Errorf("r.Lookup(ctx, %q) user != test.wantUser (-got +want)\n%s", test.name, diff)
			}
			if got, want := err, test.wantErr; !cmp.Equal(got, want, cmpopts.EquateErrors()) {
				t.Errorf("r.Lookup(ctx, %q) err = %v; want %v", test.name, got, want)
			}
		})
	}
}

func (s *UsersTestSuite) TestResolveEmail() {
	t := s.T()
	ctx := context.Background()
	r := s.seedAlice(ctx, t)
	for _, test := range []struct {
		desc         string
		emailAddress string
		wantName     string
		wantUser     *pb.User
		wantErr      error
	}{
		{
			desc:         "OK",
			emailAddress: testresources.Alice.EmailAddress,
			wantName:     testresources.Alice.Name,
			wantUser:     testresources.Alice,
			wantErr:      nil,
		},
		{
			desc:         "EmptyEmailAddress",
			emailAddress: "",
			wantName:     "",
			wantUser:     nil,
			wantErr:      &EmailAddressNotFound{EmailAddress: ""},
		},
		{
			desc:         "NotFound",
			emailAddress: testresources.Bob.EmailAddress,
			wantName:     "",
			wantUser:     nil,
			wantErr:      &EmailAddressNotFound{EmailAddress: testresources.Bob.EmailAddress},
		},
	} {
		t.Run(test.desc, func(t *testing.T) {
			name, err := r.ResolveEmail(ctx, test.emailAddress)
			if name != test.wantName {
				t.Errorf("r.ResolveEmail(ctx, %q) name = %q; want %q", test.emailAddress, name, test.wantName)
			}
			if got, want := err, test.wantErr; !cmp.Equal(got, want, cmpopts.EquateErrors()) {
				t.Errorf("r.ResolveEmail(ctx, %q) err = %v; want %v", test.emailAddress, got, want)
			}
		})
	}
}

func (s *UsersTestSuite) TestList() {
	t := s.T()
	ctx := context.Background()
	allUsers := []*pb.User{
		testresources.Alice,
		testresources.Bob,
		testresources.Carol,
	}
	r := s.seedAliceBobCarol(ctx, t)
	users, err := r.List(ctx)
	if diff := cmp.Diff(users, allUsers, protocmp.Transform(), cmpopts.SortSlices(userLess)); diff != "" {
		t.Errorf("r.List(ctx) users != allUsers (-got +want)\n%s", diff)
	}
	if err != nil {
		t.Errorf("r.List(ctx) err = %v; want nil", err)
	}
}

func (s *UsersTestSuite) TestCreate() {
	t := s.T()
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
			want:     users.ErrPasswordEmpty,
		},
		{
			name:     "DuplicateEmail",
			user:     &pb.User{Name: testresources.Bob.Name, EmailAddress: testresources.Alice.EmailAddress, DisplayName: testresources.Bob.DisplayName},
			password: testresources.BobPassword,
			want:     &EmailAddressExists{EmailAddress: testresources.Alice.EmailAddress},
		},
		{
			name:     "DuplicateName",
			user:     &pb.User{Name: testresources.Alice.Name, EmailAddress: testresources.Bob.EmailAddress, DisplayName: testresources.Bob.DisplayName},
			password: testresources.BobPassword,
			want:     &Exists{Name: testresources.Alice.Name},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			r := s.seedAlice(ctx, t)
			if got := r.Create(ctx, test.user, test.password); !cmp.Equal(got, test.want, cmpopts.EquateErrors()) {
				t.Errorf("r.Create(ctx, %v, %q) = %v; want %v", test.user, test.password, got, test.want)
			}
		})
	}
}

func (s *UsersTestSuite) TestUpdate() {
	t := s.T()
	ctx := context.Background()
	t.Run("OK", func(t *testing.T) {
		r := s.seedAlice(ctx, t)
		oldAlice := users.Clone(testresources.Alice)
		newAlice := users.Clone(oldAlice)
		newAlice.EmailAddress = "new-alice@example.com"
		newAlice.DisplayName = "New Alice"
		if err := r.Update(ctx, newAlice); err != nil {
			t.Fatalf("r.Update(ctx, %v) = %v; want nil", newAlice, err)
		}

		// Verify that the new version of Alice can be looked up.
		user, err := r.Lookup(ctx, newAlice.Name)
		if err != nil {
			t.Errorf("r.Lookup(ctx, %q) err = %v; want nil", newAlice.Name, err)
		}
		if diff := cmp.Diff(user, newAlice, protocmp.Transform()); diff != "" {
			t.Errorf("r.Lookup(ctx, %q) user != newAlice (-got +want)\n%s", newAlice.Name, diff)
		}

		// Verify that the new email address resolves to Alice's name.
		name, err := r.ResolveEmail(ctx, newAlice.EmailAddress)
		if err != nil {
			t.Errorf("r.ResolveEmail(ctx, %q) err = %v; want nil", newAlice.EmailAddress, err)
		}
		if name != newAlice.Name {
			t.Errorf("r.ResolveEmail(ctx, %q) name = %q; want %q", newAlice.EmailAddress, name, newAlice.Name)
		}

		// Verify that resolving the old email address fails.
		_, err = r.ResolveEmail(ctx, oldAlice.EmailAddress)
		if want := (&EmailAddressNotFound{EmailAddress: oldAlice.EmailAddress}); !cmp.Equal(err, want, cmpopts.EquateErrors()) {
			t.Errorf("r.ResolveEmail(ctx, %q) err = %v; want %v", oldAlice.EmailAddress, err, want)
		}
	})

	t.Run("Errors", func(t *testing.T) {
		r := s.seedAliceBob(ctx, t)
		for _, test := range []struct {
			desc string
			user *pb.User
			want error
		}{
			{
				desc: "ExistingEmailAddress",
				user: func() *pb.User {
					newAlice := users.Clone(testresources.Alice)
					newAlice.EmailAddress = testresources.Bob.EmailAddress
					return newAlice
				}(),
				want: &EmailAddressExists{EmailAddress: testresources.Bob.EmailAddress},
			},
		} {
			t.Run(test.desc, func(t *testing.T) {
				if err := r.Update(ctx, test.user); !cmp.Equal(err, test.want, cmpopts.EquateErrors()) {
					t.Errorf("r.Update(ctx, %v) = %v; want %v", test.user, err, test.want)
				}
			})
		}
	})
}

func (s *UsersTestSuite) TestDelete() {
	t := s.T()
	ctx := context.Background()

	t.Run("OK", func(t *testing.T) {
		r := s.seedAlice(ctx, t)
		if err := r.Delete(ctx, testresources.Alice.Name); err != nil {
			t.Fatalf("r.Delete(ctx, %q) = %v; want nil", testresources.Alice.Name, err)
		}

		// Verify that looking up the deleted user fails.
		_, err := r.Lookup(ctx, testresources.Alice.Name)
		if want := (&NotFound{Name: testresources.Alice.Name}); !cmp.Equal(err, want, cmpopts.EquateErrors()) {
			t.Errorf("r.Lookup(ctx, %q) err = %v; want %v", testresources.Alice.Name, err, want)
		}

		// Verify that resolving the email address of the deleted user
		// fails.
		_, err = r.ResolveEmail(ctx, testresources.Alice.EmailAddress)
		if want := (&EmailAddressNotFound{EmailAddress: testresources.Alice.EmailAddress}); !cmp.Equal(err, want, cmpopts.EquateErrors()) {
			t.Errorf("r.ResolveEmail(ctx, %q) err = %v; want %v", testresources.Alice.EmailAddress, err, want)
		}
	})

	t.Run("Errors", func(t *testing.T) {
		for _, test := range []struct {
			desc string
			name string
			want error
		}{
			{
				desc: "EmptyName",
				name: "",
				want: resourcename.ErrInvalidName,
			},
			{
				desc: "InvalidName",
				name: testresources.Beer.Name, // name of a product
				want: resourcename.ErrInvalidName,
			},
			{
				desc: "NotFound",
				name: testresources.Bob.Name,
				want: &NotFound{Name: testresources.Bob.Name},
			},
		} {
			t.Run(test.desc, func(t *testing.T) {
				r := s.seedAlice(ctx, t)
				if got := r.Delete(ctx, test.name); !cmp.Equal(got, test.want, cmpopts.EquateErrors()) {
					t.Errorf("r.Delete(ctx, %q) = %v; want %v", test.name, got, test.want)
				}
			})
		}
	})
}
