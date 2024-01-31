package bank

import (
	"bank-api/internal/model"
	"bank-api/internal/repository"
	"context"
	"errors"
)

var (
	ErrNotEnoughMoney = errors.New("not enough money")
	ErrInvalidAmount  = errors.New("invalid amount")

	ErrNoTransactions = errors.New("no transactions")
)

type TransactionService interface {
	ProcessTransaction(ctx context.Context, transaction *model.Transaction) error
	ListTransactions(ctx context.Context, accountId int) ([]*model.Transaction, error)
}

type transactionService struct {
	repo repository.AccountRepository
}

func NewTransactionService(repo repository.AccountRepository) TransactionService {
	return &transactionService{repo: repo}
}

func (s *transactionService) ProcessTransaction(ctx context.Context, transaction *model.Transaction) error {
	if transaction.Amount <= 0 {
		return ErrInvalidAmount
	}

	switch transaction.Type {
	case model.Deposit:
		return s.processDeposit(ctx, transaction)
	case model.Withdraw:
		return s.processWithdraw(ctx, transaction)
	case model.Transfer:
		return s.processTransfer(ctx, transaction)
	}

	return nil
}

func (s *transactionService) processDeposit(ctx context.Context, transaction *model.Transaction) error {
	accTo, err := s.repo.GetAccount(ctx, transaction.ToAccountId)
	if err != nil {
		return ErrInvalidAccount
	}
	if accTo.UserId != transaction.UserId {
		return ErrInvalidAccount
	}

	if err := s.repo.Transaction(ctx, transaction.ToAccountId, transaction.Amount, transaction.Type); err != nil {
		return ErrInternal
	}

	return nil
}

func (s *transactionService) processWithdraw(ctx context.Context, transaction *model.Transaction) error {
	accFrom, err := s.repo.GetAccount(ctx, transaction.FromAccountId)
	if err != nil {
		return ErrInvalidAccount
	}
	if accFrom.UserId != transaction.UserId || accFrom.Amount < transaction.Amount {
		return ErrInvalidAccount
	}

	if err := s.repo.Transaction(ctx, transaction.FromAccountId, transaction.Amount, transaction.Type); err != nil {
		return ErrInternal
	}

	return nil
}

func (s *transactionService) processTransfer(ctx context.Context, transaction *model.Transaction) error {
	accFrom, err := s.repo.GetAccount(ctx, transaction.FromAccountId)
	if err != nil || accFrom.UserId != transaction.UserId {
		return ErrInvalidAccount
	}

	accTo, err := s.repo.GetAccount(ctx, transaction.ToAccountId)
	if err != nil || accFrom.Cur != accTo.Cur || accFrom.Amount < transaction.Amount {
		return ErrInvalidAccount
	}

	if err := s.repo.Transfer(ctx, transaction.FromAccountId, transaction.ToAccountId, transaction.Amount); err != nil {
		return ErrInternal
	}

	return nil
}

func (s *transactionService) ListTransactions(ctx context.Context, accountId int) ([]*model.Transaction, error) {
	acc, err := s.repo.GetAccount(ctx, accountId)
	if err != nil {
		return nil, ErrInvalidAccount
	}
	if acc.UserId != accountId {
		return nil, ErrInvalidAccount
	}
	trs, err := s.repo.ListTransactions(ctx, accountId)
	if err != nil {
		return nil, ErrInternal
	}
	if len(trs) == 0 {
		return nil, ErrNoTransactions
	}
	return trs, nil
}
