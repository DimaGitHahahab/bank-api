package model

type TransactionType int

const (
	Deposit TransactionType = iota
	Withdraw
)

type Transaction struct {
	UserId     int
	AccountId  int
	CurrencyId int
	Amount     int
	Type       TransactionType
}

type Transfer struct {
	UserId      int
	AccountId   int
	ToAccountId int
	CurrencyId  int
	Amount      int
}
