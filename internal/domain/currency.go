package domain

import "errors"

var ErrNoSuchCurrency = errors.New("no such currency")

type Currency struct {
	Id     int
	Symbol string
}
