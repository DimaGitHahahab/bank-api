package model

type Account struct {
	Id     int
	UserId int
	Cur    Currency
	Amount int
}
