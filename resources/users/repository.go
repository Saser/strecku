package users

import (
	"context"
	"errors"
	"fmt"
	"testing"

	pb "github.com/Saser/strecku/api/v1"
	"google.golang.org/protobuf/proto"
)

type Repository struct {
	users     map[string]*pb.User // name -> user
	passwords map[string]string   // name -> password (plaintext)
	names     map[string]string   // email address -> name
}

var ErrPasswordEmpty = errors.New("empty password")

type NotFoundError struct {
	Name         string
	EmailAddress string
}

func (e *NotFoundError) Error() string {
	var (
		msg   string
		query string
	)
	switch {
	case e.Name != "":
		msg = "user not found"
		query = e.Name
	case e.EmailAddress != "":
		msg = "user email not found"
		query = e.EmailAddress
	}
	return fmt.Sprintf("%s: %q", msg, query)
}

func (e *NotFoundError) Is(target error) bool {
	other, ok := target.(*NotFoundError)
	if !ok {
		return false
	}
	return e.Name == other.Name && e.EmailAddress == other.EmailAddress
}

type ExistsError struct {
	Name         string
	EmailAddress string
}

func (e *ExistsError) Error() string {
	var (
		msg   string
		query string
	)
	switch {
	case e.Name != "":
		msg = "user exists"
		query = e.Name
	case e.EmailAddress != "":
		msg = "user email exists"
		query = e.EmailAddress
	}
	return fmt.Sprintf("%s: %q", msg, query)
}

func (e *ExistsError) Is(target error) bool {
	other, ok := target.(*ExistsError)
	if !ok {
		return false
	}
	return e.Name == other.Name && e.EmailAddress == other.EmailAddress
}

type WrongPasswordError struct {
	Name string
}

func (e *WrongPasswordError) Error() string {
	return fmt.Sprintf("wrong password for user %q", e.Name)
}

func (e *WrongPasswordError) Is(target error) bool {
	other, ok := target.(*WrongPasswordError)
	if !ok {
		return false
	}
	return e.Name == other.Name
}

func NewRepository() *Repository {
	return newRepository(make(map[string]*pb.User), make(map[string]string), make(map[string]string))
}

func SeedRepository(t *testing.T, users []*pb.User, passwords []string) *Repository {
	t.Helper()
	userCount := len(users)
	if passwordCount := len(passwords); passwordCount != userCount {
		t.Fatalf("len(users), len(passwords) = %v, %v; want equal", userCount, passwordCount)
	}
	mUsers := make(map[string]*pb.User, userCount)
	mPasswords := make(map[string]string, userCount)
	mNames := make(map[string]string, userCount)
	for i, user := range users {
		if err := Validate(user); err != nil {
			t.Errorf("Validate(%v) = %v; want nil", user, err)
		}
		mUsers[user.Name] = user
		mPasswords[user.Name] = passwords[i]
		mNames[user.EmailAddress] = user.Name
	}
	if t.Failed() {
		t.FailNow()
	}
	return newRepository(mUsers, mPasswords, mNames)
}

func newRepository(users map[string]*pb.User, passwords map[string]string, names map[string]string) *Repository {
	return &Repository{
		users:     users,
		passwords: passwords,
		names:     names,
	}
}

func Clone(user *pb.User) *pb.User {
	return proto.Clone(user).(*pb.User)
}

func (r *Repository) Authenticate(_ context.Context, name string, password string) error {
	storedPassword, ok := r.passwords[name]
	if !ok {
		return &NotFoundError{Name: name}
	}
	if password != storedPassword {
		return &WrongPasswordError{Name: name}
	}
	return nil
}

func (r *Repository) LookupUser(_ context.Context, name string) (*pb.User, error) {
	if err := ValidateName(name); err != nil {
		return nil, err
	}
	user, ok := r.users[name]
	if !ok {
		return nil, &NotFoundError{Name: name}
	}
	return Clone(user), nil
}

func (r *Repository) LookupUserByEmail(ctx context.Context, emailAddress string) (*pb.User, error) {
	name, ok := r.names[emailAddress]
	if !ok {
		return nil, &NotFoundError{EmailAddress: emailAddress}
	}
	return r.LookupUser(ctx, name)
}

func (r *Repository) ListUsers(ctx context.Context) ([]*pb.User, error) {
	return r.FilterUsers(ctx, func(*pb.User) bool { return true })
}

func (r *Repository) FilterUsers(_ context.Context, predicate func(*pb.User) bool) ([]*pb.User, error) {
	var filtered []*pb.User
	for _, user := range r.users {
		if predicate(user) {
			filtered = append(filtered, Clone(user))
		}
	}
	return filtered, nil
}

func (r *Repository) CreateUser(_ context.Context, user *pb.User, password string) error {
	name := user.Name
	emailAddress := user.EmailAddress
	if _, exists := r.users[name]; exists {
		return &ExistsError{Name: user.Name}
	}
	if _, exists := r.names[emailAddress]; exists {
		return &ExistsError{EmailAddress: user.EmailAddress}
	}
	if password == "" {
		return ErrPasswordEmpty
	}
	r.users[name] = user
	r.passwords[name] = password
	r.names[emailAddress] = user.Name
	return nil
}

func (r *Repository) UpdateUser(_ context.Context, updated *pb.User) error {
	if err := Validate(updated); err != nil {
		return err
	}
	old, exists := r.users[updated.Name]
	if !exists {
		return &NotFoundError{Name: updated.Name}
	}
	if name, exists := r.names[updated.EmailAddress]; exists && name != updated.Name {
		return &ExistsError{EmailAddress: updated.EmailAddress}
	}
	delete(r.names, old.EmailAddress)
	r.names[updated.EmailAddress] = updated.Name
	r.users[updated.Name] = updated
	return nil
}

func (r *Repository) DeleteUser(_ context.Context, name string) error {
	if err := ValidateName(name); err != nil {
		return err
	}
	user, ok := r.users[name]
	if !ok {
		return &NotFoundError{Name: name}
	}
	delete(r.names, user.EmailAddress)
	delete(r.users, name)
	return nil
}
