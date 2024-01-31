package queries

import (
	"bank-api/internal/model"
	"context"
	"errors"
)

var (
	ErrNoSuchAccount = errors.New("no such account")
)

const currencyIdBySymbol = `
SELECT id FROM currency
WHERE symbol = $1
`

func (q *Queries) GetCurrencyId(ctx context.Context, cur model.Currency) (int, error) {
	var id int
	err := q.pool.QueryRow(ctx, currencyIdBySymbol, cur.Symbol).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

const currencySymbolById = `
SELECT symbol FROM currency
WHERE id = $1
`

func (q *Queries) GetCurrencySymbol(ctx context.Context, id int) (string, error) {
	var symbol string
	err := q.pool.QueryRow(ctx, currencySymbolById, id).Scan(&symbol)
	if err != nil {
		return "", err
	}
	return symbol, nil
}

const createAccount = `
INSERT INTO account (user_id, currency_id)
VALUES ($1, (SELECT id FROM currency WHERE symbol = $2))
RETURNING id, user_id, amount
`

func (q *Queries) CreateAccount(ctx context.Context, userId int, cur model.Currency) (*model.Account, error) {
	var account model.Account
	_ = q.pool.QueryRow(ctx, createAccount, userId, cur.Symbol).Scan(&account.Id, &account.UserId, &account.Amount)
	account.Cur = cur
	return &account, nil
}

const getAccount = `
SELECT account.id, account.user_id, currency.symbol, account.amount
FROM account
JOIN currency ON currency.id = account.currency_id
WHERE account.id = $1
`

func (q *Queries) GetAccount(ctx context.Context, accountId int) (*model.Account, error) {
	var account model.Account
	err := q.pool.QueryRow(ctx, getAccount, accountId).Scan(&account.Id, &account.UserId, &account.Cur.Symbol, &account.Amount)
	if err != nil {
		return nil, ErrNoSuchAccount
	}
	account.Cur.Id, _ = q.GetCurrencyId(ctx, account.Cur)
	return &account, nil
}

const updateAccount = `
UPDATE account
SET amount = $2
WHERE id = $1
RETURNING id, user_id, currency_id, amount
`

func (q *Queries) UpdateAccount(ctx context.Context, accountId int, amount int) (*model.Account, error) {
	var account model.Account
	err := q.pool.QueryRow(ctx, updateAccount, accountId, amount).Scan(&account.Id, &account.UserId, &account.Cur.Id, &account.Amount)
	if err != nil {
		return nil, ErrNoSuchAccount
	}
	account.Cur.Symbol, _ = q.GetCurrencySymbol(ctx, account.Cur.Id)
	return &account, nil
}

const deleteAccount = `
DELETE FROM account
WHERE id = $1
`

func (q *Queries) DeleteAccount(ctx context.Context, id int) error {
	if _, err := q.pool.Exec(ctx, deleteAccount, id); err != nil {
		return ErrNoSuchAccount
	}
	return nil
}
