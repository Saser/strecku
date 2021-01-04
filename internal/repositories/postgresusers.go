package repositories

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Masterminds/squirrel"
	pb "github.com/Saser/strecku/api/v1"
	"github.com/Saser/strecku/internal/database"
	"github.com/Saser/strecku/resourcename"
	"github.com/Saser/strecku/resources/users"
	"github.com/google/uuid"
)

type PostgresUsers struct {
	db *sql.DB
	sb squirrel.StatementBuilderType
}

var _ Users = (*PostgresUsers)(nil)

func NewPostgresUsers(db *sql.DB) *PostgresUsers {
	return &PostgresUsers{
		db: db,
		sb: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}

func (u *PostgresUsers) Authenticate(ctx context.Context, name string, password string) error {
	return errors.New("method Authenticate not implemented")
}

func (u *PostgresUsers) Lookup(ctx context.Context, name string) (*pb.User, error) {
	id, err := users.ParseName(name)
	if err != nil {
		return nil, err
	}
	var user *pb.User
	if err := database.InTx(ctx, u.db, func(tx *sql.Tx) error {
		query := u.sb.
			Select("email_address", "display_name").
			From("users").
			Where(squirrel.Eq{
				"uuid": id,
			}).
			RunWith(tx)
		var (
			emailAddress string
			displayName  string
		)
		if err := query.ScanContext(ctx, &emailAddress, &displayName); err != nil {
			return err
		}
		user = &pb.User{
			Name:         name,
			EmailAddress: emailAddress,
			DisplayName:  displayName,
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return user, nil
}

func (u *PostgresUsers) ResolveEmail(ctx context.Context, emailAddress string) (string, error) {
	var name string
	if err := database.InTx(ctx, u.db, func(tx *sql.Tx) error {
		query := u.sb.
			Select("uuid").
			From("users").
			Where(squirrel.Eq{"email_address": emailAddress}).
			RunWith(tx)
		var id uuid.UUID
		if err := query.QueryRowContext(ctx).Scan(&id); err != nil {
			return err
		}
		var err error
		name, err = users.NameFormat.Format(resourcename.UUIDs{"user": id})
		return err
	}); err != nil {
		return "", err
	}
	return name, nil
}

func (u *PostgresUsers) List(ctx context.Context) ([]*pb.User, error) {
	var allUsers []*pb.User
	if err := database.InTx(ctx, u.db, func(tx *sql.Tx) error {
		query := u.sb.
			Select("uuid", "email_address", "display_name").
			From("users").
			RunWith(tx)
		rows, err := query.QueryContext(ctx)
		if err != nil {
			return err
		}
		for rows.Next() {
			var (
				id           uuid.UUID
				emailAddress string
				displayName  string
			)
			if err := rows.Scan(&id, &emailAddress, &displayName); err != nil {
				return err
			}
			name, _ := users.NameFormat.Format(resourcename.UUIDs{"user": id})
			allUsers = append(allUsers, &pb.User{
				Name:         name,
				EmailAddress: emailAddress,
				DisplayName:  displayName,
			})
		}
		return rows.Err()
	}); err != nil {
		return nil, err
	}
	return allUsers, nil
}

func (u *PostgresUsers) Create(ctx context.Context, user *pb.User, password string) error {
	id, err := users.ParseName(user.Name)
	if err != nil {
		return err
	}
	return database.InTx(ctx, u.db, func(tx *sql.Tx) error {
		query := u.sb.
			Insert("users").
			SetMap(map[string]interface{}{
				"uuid":          id,
				"deleted":       false,
				"email_address": user.EmailAddress,
				"display_name":  user.DisplayName,
			}).
			RunWith(tx)
		_, err := query.ExecContext(ctx)
		return err
	})
}

func (u *PostgresUsers) Update(ctx context.Context, user *pb.User) error {
	id, err := users.ParseName(user.Name)
	if err != nil {
		return err
	}
	return database.InTx(ctx, u.db, func(tx *sql.Tx) error {
		query := u.sb.
			Update("users").
			SetMap(map[string]interface{}{
				"email_address": user.EmailAddress,
				"display_name":  user.DisplayName,
			}).
			Where(squirrel.Eq{"uuid": id}).
			RunWith(tx)
		_, err := query.ExecContext(ctx)
		return err
	})
}

func (u *PostgresUsers) Delete(ctx context.Context, name string) error {
	id, err := users.ParseName(name)
	if err != nil {
		return err
	}
	if err := database.InTx(ctx, u.db, func(tx *sql.Tx) error {
		query := u.sb.
			Delete("users").
			Where(squirrel.Eq{"uuid": id}).
			RunWith(tx)
		_, err := query.ExecContext(ctx)
		return err
	}); err != nil {
		return err
	}
	return nil
}
