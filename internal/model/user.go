package model

import (
	"database/sql"
	"time"
)

type Role int32

const (
	Role_UNKNOWN Role = 0
	Role_USER    Role = 1
	Role_ADMIN   Role = 2
)

type UserUpdate struct {
	ID       int64
	UserInfo UserInfo
}

type UserCreate struct {
	UserInfo UserInfo
	Password string
}

type User struct {
	ID        int64
	UserInfo  UserInfo
	Password  string
	CreatedAt time.Time
	UpdatedAt sql.NullTime
}

type UserInfo struct {
	Name  string
	Email string
	Role  Role
}
