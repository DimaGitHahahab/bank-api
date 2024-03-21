package domain

import (
	"errors"
	"time"
)

type TransactionType int

var (
	ErrNotEnoughMoney = errors.New("not enough money")
	ErrInvalidAmount  = errors.New("invalid amount")
)

const (
	Deposit TransactionType = iota
	Withdraw

	Transfer
)

type Transaction struct {
	UserId        int
	FromAccountId int
	ToAccountId   int
	Cur           Currency
	Amount        int
	Type          TransactionType
	Time          time.Time
}
