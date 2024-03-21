package domain

import "errors"

var (
	ErrInvalidAccount = errors.New("invalid account")
	ErrNoSuchAccount  = errors.New("no such account")
)

type Account struct {
	Id     int
	UserId int
	Cur    Currency
	Amount int
}
