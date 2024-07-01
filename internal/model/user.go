package model

import (
	"database/sql"
	"github.com/dgrijalva/jwt-go"
	"time"
)

type UserRole int8

const (
	UNKNOWN UserRole = iota
	USER
	ADMIN
)

const (
	UserNameFieldCode = "username"
	IDFieldCode       = "id"
)

type User struct {
	ID        int64
	Info      UserInfo
	CreatedAt time.Time
	UpdatedAt sql.NullTime
	Password  string
}

type UserInfo struct {
	Username string
	Email    string
	Role     UserRole
}

type UserCreate struct {
	Info     UserInfo
	Password string
}

type UserUpdate struct {
	ID   int64
	Info UserInfo
}

type UserClaims struct {
	jwt.StandardClaims
	Username string
	Role     UserRole
}
