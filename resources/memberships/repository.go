package memberships

import (
	"context"
	"fmt"
	"testing"

	pb "github.com/Saser/strecku/api/v1"
	"github.com/Saser/strecku/resources/stores"
	"github.com/Saser/strecku/resources/users"
	"google.golang.org/protobuf/proto"
)

type MembershipNotFoundError struct {
	Name        string
	User, Store string
}

func (e *MembershipNotFoundError) Error() string {
	var query string
	switch {
	case e.Name != "":
		query = fmt.Sprintf("%q", e.Name)
	case e.User != "" && e.Store != "":
		query = fmt.Sprintf("between %q and %q", e.User, e.Store)
	}
	return fmt.Sprintf("membership not found: %s", query)
}

func (e *MembershipNotFoundError) Is(target error) bool {
	other, ok := target.(*MembershipNotFoundError)
	if !ok {
		return false
	}
	return e.Name == other.Name && e.User == other.User && e.Store == other.Store
}

type MembershipExistsError struct {
	Name        string
	User, Store string
}

func (e *MembershipExistsError) Error() string {
	var query string
	switch {
	case e.Name != "":
		query = fmt.Sprintf("%q", e.Name)
	case e.User != "" && e.Store != "":
		query = fmt.Sprintf("between %q and %q", e.User, e.Store)
	}
	return fmt.Sprintf("membership exists: %s", query)
}

func (e *MembershipExistsError) Is(target error) bool {
	other, ok := target.(*MembershipExistsError)
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
		return nil, &MembershipNotFoundError{Name: name}
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
		return nil, &MembershipNotFoundError{User: user, Store: store}
	}
	return r.LookupMembership(ctx, name)
}
