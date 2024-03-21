package domain

import (
	"errors"
	"time"
)

var (
	ErrEmptyUserInfo     = errors.New("empty user info")
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrInvalidEmail      = errors.New("invalid email")
	ErrNoSuchUser        = errors.New("no such user")
	ErrEmptyPassword     = errors.New("empty password")
)

type UserInfo struct {
	Name     string
	Email    string
	Password string
}

type User struct {
	Id             int
	Name           string
	Email          string
	HashedPassword string
	CreatedAt      time.Time
}
