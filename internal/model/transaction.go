package model

import "time"

type TransactionType int

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
