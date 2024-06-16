package model

import (
	"database/sql"
	"time"
)

type UserUpdate struct {
	ID       int64    `db:"id"`
	UserInfo UserInfo `db:""`
}

type UserCreate struct {
	UserInfo UserInfo `db:""`
	Password string   `db:"password"`
}

type User struct {
	ID        int64        `db:"id"`
	UserInfo  UserInfo     `db:""`
	Password  string       `db:"password"`
	CreatedAt time.Time    `db:"created_at"`
	UpdatedAt sql.NullTime `db:"updated_at"`
}

type UserInfo struct {
	Name  string `db:"username"`
	Email string `db:"email"`
	Role  int32  `db:"role"`
}
