package model

import (
	"database/sql"
	"time"
)

type Role int32

const (
	Role_USER  Role = 0
	Role_ADMIN Role = 1
)

type User struct {
	ID        int64        `db:"id"`
	Info      *Info        `db:""`
	Password  string       `db:"password"`
	CreatedAt time.Time    `db:"created_at"`
	UpdatedAt sql.NullTime `db:"updated_at"`
}
type Info struct {
	Name  string `db:"name"`
	Email string `db:"email"`
	Role  Role   `db:"role"`
}
