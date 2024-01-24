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
	ID        int64
	Info      *UserInfo
	Password  string
	CreatedAt time.Time
	UpdatedAt sql.NullTime
}
type UserInfo struct {
	Name  string
	Email string
	Role  Role
}
