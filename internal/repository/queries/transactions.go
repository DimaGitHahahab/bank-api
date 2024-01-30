package queries

import (
	"context"
	"errors"
)

var (
	ErrTransactionFailed = errors.New("transaction failed")
)

func (q *Queries) Transaction(ctx context.Context, accountId int, amount int) error {
	tx, err := q.pool.Begin(ctx)
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		}
	}()
	if err != nil {
		return ErrTransactionFailed
	}

	if _, err := tx.Exec(ctx, updateAccount, accountId, amount); err != nil {
		return ErrTransactionFailed
	}

	if err := tx.Commit(ctx); err != nil {
		return ErrTransactionFailed
	}

	return nil
}

const getBalance = `
SELECT amount FROM account
WHERE id = $1 
`

func (q *Queries) Transfer(ctx context.Context, fromAccountId int, toAccountId int, amount int) error {
	tx, err := q.pool.Begin(ctx)
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		}
	}()
	if err != nil {
		return ErrTransactionFailed
	}

	var fromAccountBalance int
	if err = tx.QueryRow(ctx, getBalance, fromAccountId).Scan(&fromAccountBalance); err != nil {
		return ErrTransactionFailed
	}

	fromAccountBalance -= amount

	if _, err = tx.Exec(ctx, updateAccount, fromAccountId, fromAccountBalance); err != nil {
		return ErrTransactionFailed
	}

	var toAccountBalance int
	if err = tx.QueryRow(ctx, getBalance, toAccountId).Scan(&toAccountBalance); err != nil {
		return ErrTransactionFailed
	}
	toAccountBalance += amount

	if _, err := tx.Exec(ctx, updateAccount, toAccountId, toAccountBalance); err != nil {
		return ErrTransactionFailed
	}

	if err := tx.Commit(ctx); err != nil {
		return ErrTransactionFailed
	}

	return nil
}
