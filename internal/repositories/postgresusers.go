package repositories

import (
	"context"

	pb "github.com/Saser/strecku/api/v1"
	"github.com/jackc/pgx/v4/pgxpool"
)

type PostgresUsers struct {
	inmemory *InMemoryUsers

	pool *pgxpool.Pool
}

var _ Users = (*PostgresUsers)(nil)

func NewPostgresUsers(pool *pgxpool.Pool) *PostgresUsers {
	return &PostgresUsers{
		inmemory: NewInMemoryUsers(),

		pool: pool,
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
