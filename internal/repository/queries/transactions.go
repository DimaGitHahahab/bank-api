package queries

import (
	"bank-api/internal/model"
	"context"
	"database/sql"
	"errors"
)

var (
	ErrTransactionFailed = errors.New("transaction failed")
	ErrNoTransactions    = errors.New("no transactions")
)

const addTransactionEntry = `
INSERT INTO transaction (from_account_id, to_account_id, currency_id, amount)
VALUES ($1, $2, $3, $4)
`

const getCurrencyByAccountId = `
SELECT currency_id FROM account
WHERE id = $1
`

func (q *Queries) Transaction(ctx context.Context, accountId int, amount int, t model.TransactionType) error {
	tx, err := q.pool.Begin(ctx)
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		}
	}()
	if err != nil {
		return ErrTransactionFailed
	}

	var accountBalance int
	if err = tx.QueryRow(ctx, getBalance, accountId).Scan(&accountBalance); err != nil {
		return ErrTransactionFailed
	}
	if t == model.Withdraw {
		_, err = tx.Exec(ctx, updateAccount, accountId, accountBalance-amount)
	} else {
		_, err = tx.Exec(ctx, updateAccount, accountId, accountBalance+amount)
	}

	var currencyId int
	if err := tx.QueryRow(ctx, getCurrencyByAccountId, accountId).Scan(&currencyId); err != nil {
		return ErrTransactionFailed
	}

	if t == model.Withdraw {
		_, err = tx.Exec(ctx, addTransactionEntry, accountId, nil, currencyId, amount)
		if err != nil {
			return ErrTransactionFailed
		}
	} else if t == model.Deposit {
		_, err = tx.Exec(ctx, addTransactionEntry, nil, accountId, currencyId, amount)
		if err != nil {
			return ErrTransactionFailed
		}
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

	var currencyId int
	if err := tx.QueryRow(ctx, getCurrencyByAccountId, fromAccountId).Scan(&currencyId); err != nil {
		return ErrTransactionFailed
	}
	_, err = tx.Exec(ctx, addTransactionEntry, fromAccountId, toAccountId, currencyId, amount)
	if err != nil {
		return ErrTransactionFailed
	}

	if err := tx.Commit(ctx); err != nil {
		return ErrTransactionFailed
	}

	return nil
}

const listTransactions = `
SELECT from_account_id, to_account_id, currency_id, amount, created_at FROM transaction
WHERE from_account_id = $1 OR to_account_id = $1
`

func (q *Queries) ListTransactions(ctx context.Context, accountId int) ([]*model.Transaction, error) {
	rows, err := q.pool.Query(ctx, listTransactions, accountId)
	if err != nil {
		return nil, ErrInternal
	}
	defer rows.Close()

	var transactions []*model.Transaction
	for rows.Next() {
		var transaction model.Transaction
		var from, to sql.NullInt64
		if err := rows.Scan(&from, &to, &transaction.Cur.Id, &transaction.Amount, &transaction.Time); err != nil {
			return nil, ErrInternal
		}
		if !from.Valid {
			transaction.Type = model.Deposit
			transaction.ToAccountId = int(to.Int64)
		} else if !to.Valid {
			transaction.Type = model.Withdraw
			transaction.FromAccountId = int(from.Int64)
		} else {
			transaction.Type = model.Transfer
			transaction.FromAccountId = int(from.Int64)
			transaction.ToAccountId = int(to.Int64)
		}

		transaction.Cur.Symbol, _ = q.GetCurrencySymbol(ctx, transaction.Cur.Id)

		transactions = append(transactions, &transaction)
	}
	return transactions, nil
}
