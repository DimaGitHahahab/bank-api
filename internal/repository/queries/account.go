package queries

import (
	"bank-api/internal/domain"
	"context"
	"fmt"
)

const currencyIdBySymbol = `
SELECT id FROM currency
WHERE symbol = $1
`

func (q *Queries) GetCurrencyId(ctx context.Context, cur domain.Currency) (int, error) {
	var id int
	err := q.pool.QueryRow(ctx, currencyIdBySymbol, cur.Symbol).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("error getting currency id: %w", err)
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
		return "", fmt.Errorf("error getting currency symbol: %w", err)
	}
	return symbol, nil
}

const currencyExists = `
SELECT EXISTS (
	SELECT 1
	FROM currency
	WHERE symbol = $1
)
`

func (q *Queries) CurrencyExists(ctx context.Context, cur domain.Currency) (bool, error) {
	var exists bool
	if err := q.pool.QueryRow(ctx, currencyExists, cur.Symbol).Scan(&exists); err != nil {
		return false, fmt.Errorf("error checking if currency exists: %w", err)
	}
	return exists, nil
}

const createAccount = `
INSERT INTO account (user_id, currency_id)
VALUES ($1, (SELECT id FROM currency WHERE symbol = $2))
RETURNING id, user_id, amount
`

func (q *Queries) CreateAccount(ctx context.Context, userId int, cur domain.Currency) (*domain.Account, error) {
	var account domain.Account
	if err := q.pool.QueryRow(ctx, createAccount, userId, cur.Symbol).Scan(&account.Id, &account.UserId, &account.Amount); err != nil {
		return nil, fmt.Errorf("error creating account: %w", err)
	}
	account.Cur = cur
	return &account, nil
}

const getAccount = `
SELECT account.id, account.user_id, currency.symbol, account.amount
FROM account
JOIN currency ON currency.id = account.currency_id
WHERE account.id = $1
`

func (q *Queries) GetAccount(ctx context.Context, accountId int) (*domain.Account, error) {
	var account domain.Account
	err := q.pool.QueryRow(ctx, getAccount, accountId).Scan(&account.Id, &account.UserId, &account.Cur.Symbol, &account.Amount)
	if err != nil {
		return nil, fmt.Errorf("error getting account: %w", err)
	}
	account.Cur.Id, err = q.GetCurrencyId(ctx, account.Cur)
	if err != nil {
		return nil, fmt.Errorf("error getting currency id: %w", err)
	}
	return &account, nil
}

const accountExists = `
SELECT EXISTS (
	SELECT 1
	FROM account
	WHERE id = $1
)
`

func (q *Queries) AccountExists(ctx context.Context, id int) (bool, error) {
	var exists bool
	if err := q.pool.QueryRow(ctx, accountExists, id).Scan(&exists); err != nil {
		return false, fmt.Errorf("error checking if account exists: %w", err)
	}
	return exists, nil
}

const updateAccount = `
UPDATE account
SET amount = $2
WHERE id = $1
RETURNING id, user_id, currency_id, amount
`

func (q *Queries) UpdateAccount(ctx context.Context, accountId int, amount int) (*domain.Account, error) {
	var account domain.Account
	err := q.pool.QueryRow(ctx, updateAccount, accountId, amount).Scan(&account.Id, &account.UserId, &account.Cur.Id, &account.Amount)
	if err != nil {
		return nil, fmt.Errorf("error updating account: %w", err)
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
		return fmt.Errorf("error deleting account: %w", err)
	}
	return nil
}
