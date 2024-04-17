package queries

import (
	"context"
	"database/sql"
	"fmt"

	"bank-api/internal/domain"
)

const addTransactionEntry = `
INSERT INTO transaction (from_account_id, to_account_id, currency_id, amount)
VALUES ($1, $2, $3, $4)
`

const getCurrencyByAccountId = `
SELECT currency_id FROM account
WHERE id = $1
`

func (q *Queries) Transaction(ctx context.Context, accountId int, amount int, t domain.TransactionType) error {
	tx, err := q.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("error starting transaction: %w", err)
	}

	var accountBalance int
	if err = tx.QueryRow(ctx, getBalance, accountId).Scan(&accountBalance); err != nil {
		tx.Rollback(ctx)
		return fmt.Errorf("error getting account balance: %w", err)
	}
	if t == domain.Withdraw {
		_, err = tx.Exec(ctx, updateAccount, accountId, accountBalance-amount)
	} else {
		_, err = tx.Exec(ctx, updateAccount, accountId, accountBalance+amount)
	}
	if err != nil {
		tx.Rollback(ctx)
		return fmt.Errorf("error updating account balance: %w", err)
	}

	var currencyId int
	if err := tx.QueryRow(ctx, getCurrencyByAccountId, accountId).Scan(&currencyId); err != nil {
		tx.Rollback(ctx)
		return fmt.Errorf("error getting currency id: %w", err)
	}

	if t == domain.Withdraw {
		_, err = tx.Exec(ctx, addTransactionEntry, accountId, nil, currencyId, amount)
	} else if t == domain.Deposit {
		_, err = tx.Exec(ctx, addTransactionEntry, nil, accountId, currencyId, amount)
	}
	if err != nil {
		tx.Rollback(ctx)
		return fmt.Errorf("error adding transaction entry: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("error committing transaction: %w", err)
	}

	return nil
}

const getBalance = `
SELECT amount FROM account
WHERE id = $1 
`

func (q *Queries) Transfer(ctx context.Context, fromAccountId int, toAccountId int, amount int) error {
	tx, err := q.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("error starting transaction: %w", err)
	}

	var fromAccountBalance int
	if err = tx.QueryRow(ctx, getBalance, fromAccountId).Scan(&fromAccountBalance); err != nil {
		tx.Rollback(ctx)
		return fmt.Errorf("error getting account balance: %w", err)
	}

	fromAccountBalance -= amount

	if _, err = tx.Exec(ctx, updateAccount, fromAccountId, fromAccountBalance); err != nil {
		tx.Rollback(ctx)
		return fmt.Errorf("error updating account balance: %w", err)
	}

	var toAccountBalance int
	if err = tx.QueryRow(ctx, getBalance, toAccountId).Scan(&toAccountBalance); err != nil {
		tx.Rollback(ctx)
		return fmt.Errorf("error getting account balance: %w", err)
	}
	toAccountBalance += amount

	if _, err := tx.Exec(ctx, updateAccount, toAccountId, toAccountBalance); err != nil {
		tx.Rollback(ctx)
		return fmt.Errorf("error updating account balance: %w", err)
	}

	var currencyId int
	if err := tx.QueryRow(ctx, getCurrencyByAccountId, fromAccountId).Scan(&currencyId); err != nil {
		tx.Rollback(ctx)
		return fmt.Errorf("error getting currency id: %w", err)
	}
	_, err = tx.Exec(ctx, addTransactionEntry, fromAccountId, toAccountId, currencyId, amount)
	if err != nil {
		tx.Rollback(ctx)
		return fmt.Errorf("error adding transaction entry: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("error committing transaction: %w", err)
	}

	return nil
}

const listTransactions = `
SELECT from_account_id, to_account_id, currency_id, amount, created_at FROM transaction
WHERE from_account_id = $1 OR to_account_id = $1
`

func (q *Queries) ListTransactions(ctx context.Context, accountId int) ([]*domain.Transaction, error) {
	rows, err := q.pool.Query(ctx, listTransactions, accountId)
	if err != nil {
		return nil, fmt.Errorf("error getting transactions: %w", err)
	}
	defer rows.Close()

	var transactions []*domain.Transaction
	for rows.Next() {
		var transaction domain.Transaction
		var from, to sql.NullInt64
		if err := rows.Scan(&from, &to, &transaction.Cur.Id, &transaction.Amount, &transaction.Time); err != nil {
			return nil, fmt.Errorf("error getting transaction: %w", err)
		}
		if !from.Valid {
			transaction.Type = domain.Deposit
			transaction.ToAccountId = int(to.Int64)
		} else if !to.Valid {
			transaction.Type = domain.Withdraw
			transaction.FromAccountId = int(from.Int64)
		} else {
			transaction.Type = domain.Transfer
			transaction.FromAccountId = int(from.Int64)
			transaction.ToAccountId = int(to.Int64)
		}

		transaction.Cur.Symbol, err = q.GetCurrencySymbol(ctx, transaction.Cur.Id)
		if err != nil {
			return nil, fmt.Errorf("error getting currency symbol: %w", err)
		}

		transactions = append(transactions, &transaction)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error getting transactions: %w", err)
	}

	return transactions, nil
}
