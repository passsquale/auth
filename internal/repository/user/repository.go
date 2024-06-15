package user

import (
	"context"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/passsquale/auth/internal/repository"
	"github.com/passsquale/auth/internal/repository/user/converter"
	"github.com/passsquale/auth/internal/repository/user/model"
	desc "github.com/passsquale/auth/pkg/user_v1"
)

const (
	usersTable = "users"

	idColumn        = "id"
	nameColumn      = "name"
	emailColumn     = "email"
	roleColumn      = "role"
	passwordColumn  = "password"
	createdAtColumn = "created_at"
	updatedAtColumn = "updated_at"
)

type repo struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) repository.UserRepository {
	return &repo{db: db}
}

func (r *repo) Create(ctx context.Context, cr *desc.CreateRequest) (int64, error) {
	builder := squirrel.Insert(usersTable).
		PlaceholderFormat(squirrel.Dollar).
		Columns(nameColumn, emailColumn, roleColumn, passwordColumn).
		Values(cr.Info.Name, cr.Info.Email, cr.Info.Role, cr.Password).
		Suffix(fmt.Sprintf("RETURNING %s", idColumn))
	query, args, err := builder.ToSql()
	if err != nil {
		return 0, err
	}
	var id int64
	err = r.db.QueryRow(ctx, query, args...).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *repo) Get(ctx context.Context, id int64) (*desc.User, error) {
	builder := squirrel.Select(idColumn, nameColumn, emailColumn, roleColumn,
		passwordColumn, createdAtColumn, updatedAtColumn).
		PlaceholderFormat(squirrel.Dollar).
		From(usersTable).
		Where(squirrel.Eq{idColumn: id})

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, err
	}
	var user model.User
	err = r.db.QueryRow(ctx, query, args...).Scan(&user.ID, &user.UserInfo.Name, &user.UserInfo.Email,
		&user.UserInfo.Role, &user.Password, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return converter.ToUserFromRepo(user), nil
}

func (r *repo) Update(ctx context.Context, wrap *desc.Updwrap) error {
	builder := squirrel.Update(usersTable).
		PlaceholderFormat(squirrel.Dollar).
		Set(nameColumn, wrap.Name).
		Set(emailColumn, wrap.Email).
		Set(roleColumn, wrap.Role).
		Where(squirrel.Eq{idColumn: wrap.Id})
	query, args, err := builder.ToSql()
	if err != nil {
		return err
	}
	_, err = r.db.Exec(ctx, query, args...)
	if err != nil {
		return err
	}
	return nil
}

func (r *repo) Delete(ctx context.Context, id int64) error {
	builder := squirrel.Delete(usersTable).
		PlaceholderFormat(squirrel.Dollar).
		Where(squirrel.Eq{idColumn: id})

	query, args, err := builder.ToSql()
	if err != nil {
		return err
	}
	_, err = r.db.Exec(ctx, query, args...)
	if err != nil {
		return err
	}
	return nil
}
