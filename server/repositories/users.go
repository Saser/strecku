package repositories

import (
	"fmt"

	streckuv1 "github.com/Saser/strecku/saser/strecku/v1"
)

type Users struct {
	users map[string]*streckuv1.User // name -> user
	names map[string]string          // email address -> name
}

type UserNotFoundError struct {
	Name         string
	EmailAddress string
}

func (e *UserNotFoundError) Error() string {
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

type UserExistsError struct {
	EmailAddress string
}

func (e *UserExistsError) Error() string {
	return fmt.Sprintf("duplicate user email address: %q", e.EmailAddress)
}

func NewUsers() *Users {
	return newUsers(make(map[string]*streckuv1.User), make(map[string]string))
}

func newUsers(users map[string]*streckuv1.User, names map[string]string) *Users {
	return &Users{
		users: users,
		names: names,
	}
}

func (r *Users) LookupUser(name string) (*streckuv1.User, error) {
	user, ok := r.users[name]
	if !ok {
		return nil, &UserNotFoundError{Name: name}
	}
	return user, nil
}

func (r *Users) LookupUserByEmail(emailAddress string) (*streckuv1.User, error) {
	name, ok := r.names[emailAddress]
	if !ok {
		return nil, &UserNotFoundError{EmailAddress: emailAddress}
	}
	return r.LookupUser(name)
}

func (r *Users) ListUsers() ([]*streckuv1.User, error) {
	return r.FilterUsers(func(*streckuv1.User) bool { return true })
}

func (r *Users) FilterUsers(predicate func(*streckuv1.User) bool) ([]*streckuv1.User, error) {
	var filtered []*streckuv1.User
	for _, user := range r.users {
		if predicate(user) {
			filtered = append(filtered, user)
		}
	}
	return filtered, nil
}

func (r *Users) CreateUser(user *streckuv1.User) error {
	if _, exists := r.names[user.EmailAddress]; exists {
		return &UserExistsError{EmailAddress: user.EmailAddress}
	}
	r.users[user.Name] = user
	r.names[user.EmailAddress] = user.Name
	return nil
}
