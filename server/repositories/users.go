package repositories

import (
	"fmt"

	streckuv1 "github.com/Saser/strecku/saser/strecku/v1"
)

type Users struct {
	users map[string]*streckuv1.User // name -> user
}

type UserNotFoundError struct {
	Name string
}

func (e *UserNotFoundError) Error() string {
	return fmt.Sprintf("user not found: %q", e.Name)
}

func NewUsers() *Users {
	return newUsers(make(map[string]*streckuv1.User))
}

func newUsers(users map[string]*streckuv1.User) *Users {
	return &Users{users: users}
}

func (r *Users) LookupUser(name string) (*streckuv1.User, error) {
	user, ok := r.users[name]
	if !ok {
		return nil, &UserNotFoundError{Name: name}
	}
	return user, nil
}
