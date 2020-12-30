package repositories

import (
	"context"
	"database/sql"

	pb "github.com/Saser/strecku/api/v1"
)

type PostgresUsers struct {
	inmemory *InMemoryUsers

	db *sql.DB
}

var _ Users = (*PostgresUsers)(nil)

func NewPostgresUsers(db *sql.DB) *PostgresUsers {
	return &PostgresUsers{
		inmemory: NewInMemoryUsers(),

		db: db,
	}
}

func (u *PostgresUsers) Authenticate(ctx context.Context, name string, password string) error {
	return u.inmemory.Authenticate(ctx, name, password)
}

func (u *PostgresUsers) Lookup(ctx context.Context, name string) (*pb.User, error) {
	return u.inmemory.Lookup(ctx, name)
}

func (u *PostgresUsers) ResolveEmail(ctx context.Context, emailAddress string) (string, error) {
	return u.inmemory.ResolveEmail(ctx, emailAddress)
}

func (u *PostgresUsers) List(ctx context.Context) ([]*pb.User, error) {
	return u.inmemory.List(ctx)
}

func (u *PostgresUsers) Create(ctx context.Context, user *pb.User, password string) error {
	return u.inmemory.Create(ctx, user, password)
}

func (u *PostgresUsers) Update(ctx context.Context, user *pb.User) error {
	return u.inmemory.Update(ctx, user)
}

func (u *PostgresUsers) Delete(ctx context.Context, name string) error {
	return u.inmemory.Delete(ctx, name)
}
