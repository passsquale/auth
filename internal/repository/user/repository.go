package user

import (
	"context"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/passsquale/auth/internal/client/db"
	"github.com/passsquale/auth/internal/model"
	"github.com/passsquale/auth/internal/repository"
	"github.com/passsquale/auth/internal/repository/user/converter"
	modelRepo "github.com/passsquale/auth/internal/repository/user/model"
	"time"
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
	db db.Client
}

func NewUserRepository(db db.Client) repository.UserRepository {
	return &repo{db: db}
}

func (r *repo) Create(ctx context.Context, userCreate *model.UserCreate) (int64, error) {
	builder := squirrel.Insert(usersTable).
		PlaceholderFormat(squirrel.Dollar).
		Columns(nameColumn, emailColumn, roleColumn, passwordColumn).
		Values(userCreate.Info.Username, userCreate.Info.Email, userCreate.Info.Role, userCreate.Password).
		Suffix(fmt.Sprintf("RETURNING %s", idColumn))
	query, args, err := builder.ToSql()
	if err != nil {
		return 0, err
	}
	q := db.Query{
		Name:     "user_repository.Create",
		QueryRaw: query,
	}
	var id int64
	err = r.db.DB().QueryRowContext(ctx, q, args...).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *repo) Get(ctx context.Context, id int64) (*model.User, error) {
	builder := squirrel.Select(idColumn, nameColumn, emailColumn, roleColumn,
		passwordColumn, createdAtColumn, updatedAtColumn).
		PlaceholderFormat(squirrel.Dollar).
		From(usersTable).
		Where(squirrel.Eq{idColumn: id}).
		Limit(1)

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, err
	}
	q := db.Query{
		Name:     "user_repository.Get",
		QueryRaw: query,
	}
	var user modelRepo.User
	err = r.db.DB().ScanOneContext(ctx, &user, q, args...)
	if err != nil {
		return nil, err
	}
	return converter.ToUserFromRepo(user), nil
}

func (r *repo) Update(ctx context.Context, userUpdate *model.UserUpdate) error {
	builder := squirrel.Update(usersTable).
		PlaceholderFormat(squirrel.Dollar).
		Set(nameColumn, userUpdate.Info.Username).
		Set(emailColumn, userUpdate.Info.Email).
		Set(roleColumn, userUpdate.Info.Role).
		Set(updatedAtColumn, time.Now()).
		Where(squirrel.Eq{idColumn: userUpdate.ID})
	query, args, err := builder.ToSql()
	if err != nil {
		return err
	}
	q := db.Query{
		Name:     "user_repository.Update",
		QueryRaw: query,
	}
	_, err = r.db.DB().ExecContext(ctx, q, args...)
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
	q := db.Query{
		Name:     "user_repository.Delete",
		QueryRaw: query,
	}
	_, err = r.db.DB().ExecContext(ctx, q, args...)
	if err != nil {
		return err
	}
	return nil
}
