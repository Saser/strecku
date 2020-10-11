package memberships

import (
	"context"
	"errors"
	"fmt"
	"testing"

	pb "github.com/Saser/strecku/api/v1"
	"github.com/Saser/strecku/resources/stores"
	"github.com/Saser/strecku/resources/users"
	"google.golang.org/protobuf/proto"
)

var (
	ErrUpdateUser  = errors.New("user cannot be updated")
	ErrUpdateStore = errors.New("store cannot be updated")
)

type NotFoundError struct {
	Name        string
	User, Store string
}

func (e *NotFoundError) Error() string {
	var query string
	switch {
	case e.Name != "":
		query = fmt.Sprintf("%q", e.Name)
	case e.User != "" && e.Store != "":
		query = fmt.Sprintf("between %q and %q", e.User, e.Store)
	}
	return fmt.Sprintf("membership not found: %s", query)
}

func (e *NotFoundError) Is(target error) bool {
	other, ok := target.(*NotFoundError)
	if !ok {
		return false
	}
	return e.Name == other.Name && e.User == other.User && e.Store == other.Store
}

type ExistsError struct {
	Name        string
	User, Store string
}

func (e *ExistsError) Error() string {
	var query string
	switch {
	case e.Name != "":
		query = fmt.Sprintf("%q", e.Name)
	case e.User != "" && e.Store != "":
		query = fmt.Sprintf("between %q and %q", e.User, e.Store)
	}
	return fmt.Sprintf("membership exists: %s", query)
}

func (e *ExistsError) Is(target error) bool {
	other, ok := target.(*ExistsError)
	if !ok {
		return false
	}
	return e.Name == other.Name && e.User == other.User && e.Store == other.Store
}

func Clone(membership *pb.Membership) *pb.Membership {
	return proto.Clone(membership).(*pb.Membership)
}

type composite struct {
	user  string
	store string
}

type Repository struct {
	memberships map[string]*pb.Membership // name -> membership
	names       map[composite]string      // (user name, store name) -> membership name
}

func NewRepository() *Repository {
	return newRepository(make(map[string]*pb.Membership), make(map[composite]string))
}

func SeedRepository(t *testing.T, memberships []*pb.Membership) *Repository {
	t.Helper()
	mMemberships := make(map[string]*pb.Membership, len(memberships))
	names := make(map[composite]string)
	for _, membership := range memberships {
		if err := Validate(membership); err != nil {
			t.Errorf("Validate(%v) = %v; want nil", membership, err)
		}
		if err := users.ValidateName(membership.User); err != nil {
			t.Errorf("users.ValidateName(%q) err = %v; want nil", membership.User, err)
		}
		if err := stores.ValidateName(membership.Store); err != nil {
			t.Errorf("stores.ValidateName(%q) err = %v; want nil", membership.Store, err)
		}
		mMemberships[membership.Name] = membership
		names[composite{user: membership.User, store: membership.Store}] = membership.Name
	}
	if t.Failed() {
		t.FailNow()
	}
	return newRepository(mMemberships, names)
}

func newRepository(memberships map[string]*pb.Membership, names map[composite]string) *Repository {
	return &Repository{
		memberships: memberships,
		names:       names,
	}
}

func (r *Repository) LookupMembership(_ context.Context, name string) (*pb.Membership, error) {
	if err := ValidateName(name); err != nil {
		return nil, err
	}
	membership, ok := r.memberships[name]
	if !ok {
		return nil, &NotFoundError{Name: name}
	}
	return membership, nil
}

func (r *Repository) LookupMembershipBetween(ctx context.Context, user string, store string) (*pb.Membership, error) {
	if err := users.ValidateName(user); err != nil {
		return nil, err
	}
	if err := stores.ValidateName(store); err != nil {
		return nil, err
	}
	name, ok := r.names[composite{user: user, store: store}]
	if !ok {
		return nil, &NotFoundError{User: user, Store: store}
	}
	return r.LookupMembership(ctx, name)
}

func (r *Repository) ListMemberships(ctx context.Context) ([]*pb.Membership, error) {
	return r.FilterMemberships(ctx, func(*pb.Membership) bool { return true })
}

func (r *Repository) FilterMemberships(_ context.Context, predicate func(*pb.Membership) bool) ([]*pb.Membership, error) {
	var filtered []*pb.Membership
	for _, membership := range r.memberships {
		if predicate(membership) {
			filtered = append(filtered, membership)
		}
	}
	return filtered, nil
}

func (r *Repository) CreateMembership(_ context.Context, membership *pb.Membership) error {
	name := membership.Name
	if _, exists := r.memberships[name]; exists {
		return &ExistsError{Name: name}
	}
	key := composite{user: membership.User, store: membership.Store}
	if _, exists := r.names[key]; exists {
		return &ExistsError{
			User:  membership.User,
			Store: membership.Store,
		}
	}
	r.memberships[name] = membership
	r.names[key] = name
	return nil
}

func (r *Repository) UpdateMembership(_ context.Context, updated *pb.Membership) error {
	membership, ok := r.memberships[updated.Name]
	if !ok {
		return &NotFoundError{Name: updated.Name}
	}
	if updated.User != membership.User {
		return ErrUpdateUser
	}
	if updated.Store != membership.Store {
		return ErrUpdateStore
	}
	r.memberships[updated.Name] = updated
	return nil
}

func (r *Repository) DeleteMembership(_ context.Context, name string) error {
	if err := ValidateName(name); err != nil {
		return err
	}
	membership, ok := r.memberships[name]
	if !ok {
		return &NotFoundError{Name: name}
	}
	delete(r.memberships, name)
	delete(r.names, composite{user: membership.User, store: membership.Store})
	return nil
}
