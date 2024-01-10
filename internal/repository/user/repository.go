package user

import (
	"context"
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/passsquale/auth/internal/repository"
	"github.com/passsquale/auth/internal/repository/user/converter"
	"github.com/passsquale/auth/internal/repository/user/model"
	desc "github.com/passsquale/auth/pkg/user_v1"
)

const (
	tableName = "users"

	idColumn        = "id"
	nameColumn      = "name"
	emailColumn     = "email"
	passwordColumn  = "password"
	roleColumn      = "role"
	createdAtColumn = "created_at"
	updatedAtColumn = "updated_at"
)

type repo struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) repository.UserRepository {
	return &repo{db: db}
}

func (r *repo) Create(ctx context.Context, info *desc.CreateRequest) (int64, error) {
	builder := sq.Insert(tableName).
		Columns(nameColumn, emailColumn, roleColumn, passwordColumn).
		Values(info.Info.Name, info.Info.Email, info.Password, info.Info.Role).
		Suffix("RETURNING id")
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
func (r *repo) Get(ctx context.Context, id int64) (*desc.GetResponse, error) {
	builder := sq.Select(idColumn, nameColumn, emailColumn,
		roleColumn, passwordColumn, createdAtColumn, updatedAtColumn).
		From(tableName).
		Where(sq.Eq{idColumn: id}).
		Limit(1)
	query, args, err := builder.ToSql()
	if err != nil {
		return nil, err
	}
	var user model.User
	err = r.db.QueryRow(ctx, query, args...).Scan(&user.ID, &user.Info.Name, &user.Info.Email,
		&user.Info.Role, &user.Password, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return converter.ToUserFromRepo(&user), nil
}
