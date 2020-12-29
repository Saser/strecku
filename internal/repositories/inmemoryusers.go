package repositories

import (
	"context"

	pb "github.com/Saser/strecku/api/v1"
	"github.com/Saser/strecku/resources/users"
)

type InMemoryUsers struct {
	users     map[string]*pb.User // name -> user
	passwords map[string]string   // name -> password (plaintext)
	names     map[string]string   // email address -> name
}

var _ Users = (*InMemoryUsers)(nil)

func NewInMemoryUsers() *InMemoryUsers {
	return &InMemoryUsers{
		users:     make(map[string]*pb.User),
		passwords: make(map[string]string),
		names:     make(map[string]string),
	}
}

func (u *InMemoryUsers) Authenticate(ctx context.Context, name string, password string) error {
	if err := users.ValidateName(name); err != nil {
		return err
	}
	stored, ok := u.passwords[name]
	if !ok {
		return &NotFound{Name: name}
	}
	if password != stored {
		return ErrUnauthenticated
	}
	return nil
}

func (u *InMemoryUsers) Lookup(ctx context.Context, name string) (*pb.User, error) {
	if err := users.ValidateName(name); err != nil {
		return nil, err
	}
	user, ok := u.users[name]
	if !ok {
		return nil, &NotFound{Name: name}
	}
	return users.Clone(user), nil
}

func (u *InMemoryUsers) ResolveEmail(ctx context.Context, emailAddress string) (string, error) {
	name, ok := u.names[emailAddress]
	if !ok {
		return "", &EmailAddressNotFound{EmailAddress: emailAddress}
	}
	return name, nil
}

func (u *InMemoryUsers) List(ctx context.Context) ([]*pb.User, error) {
	allUsers := make([]*pb.User, 0, len(u.users))
	for _, user := range u.users {
		allUsers = append(allUsers, users.Clone(user))
	}
	return allUsers, nil
}

func (u *InMemoryUsers) Create(ctx context.Context, user *pb.User, password string) error {
	if err := users.Validate(user); err != nil {
		return err
	}
	if err := users.ValidatePassword(password); err != nil {
		return err
	}
	if _, exists := u.users[user.Name]; exists {
		return &Exists{Name: user.Name}
	}
	if _, exists := u.names[user.EmailAddress]; exists {
		return &EmailAddressExists{EmailAddress: user.EmailAddress}
	}
	u.users[user.Name] = users.Clone(user)
	u.passwords[user.Name] = password
	u.names[user.EmailAddress] = user.Name
	return nil
}

func (u *InMemoryUsers) Update(ctx context.Context, user *pb.User) error {
	if err := users.Validate(user); err != nil {
		return err
	}
	old, exists := u.users[user.Name]
	if !exists {
		return &NotFound{Name: user.Name}
	}
	if name, exists := u.names[user.EmailAddress]; exists && name != user.Name {
		return &EmailAddressExists{EmailAddress: user.EmailAddress}
	}
	delete(u.names, old.EmailAddress)
	u.names[user.EmailAddress] = user.Name
	u.users[user.Name] = users.Clone(user)
	return nil
}

func (u *InMemoryUsers) Delete(ctx context.Context, name string) error {
	user, err := u.Lookup(ctx, name)
	if err != nil {
		return err
	}
	delete(u.names, user.EmailAddress)
	delete(u.users, user.Name)
	return nil
}
