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
)

type TransactionService interface {
	ProcessTransaction(ctx context.Context, transaction *model.Transaction) error
	ProcessTransfer(ctx context.Context, transfer *model.Transfer) error
}

type transactionService struct {
	repo repository.AccountRepository
}

func NewTransactionService(repo repository.AccountRepository) TransactionService {
	return &transactionService{repo: repo}
}

func (s *transactionService) ProcessTransaction(ctx context.Context, transaction *model.Transaction) error {
	acc, err := s.repo.GetAccount(ctx, transaction.AccountId)
	if err != nil {
		return ErrInvalidAccount
	}

	if acc.UserId != transaction.UserId {
		return ErrInvalidAccount
	}

	if transaction.Amount <= 0 {
		return ErrInvalidAmount
	}

	if transaction.Type == model.Withdraw {
		if acc.Amount < transaction.Amount {
			return ErrNotEnoughMoney
		}

		if err = s.repo.Transaction(ctx, transaction.AccountId, acc.Amount-transaction.Amount); err != nil {
			return ErrInternal
		}
	} else if transaction.Type == model.Deposit {
		if err = s.repo.Transaction(ctx, transaction.AccountId, acc.Amount+transaction.Amount); err != nil {
			return ErrInternal
		}
	}

	return nil
}

func (s *transactionService) ProcessTransfer(ctx context.Context, transfer *model.Transfer) error {
	accFrom, err := s.repo.GetAccount(ctx, transfer.AccountId)
	if err != nil {
		return ErrInvalidAccount
	}
	if accFrom.UserId != transfer.UserId {
		return ErrInvalidAccount
	}

	accTo, err := s.repo.GetAccount(ctx, transfer.ToAccountId)
	if err != nil {
		return ErrInvalidAccount
	}

	if accFrom.Cur != accTo.Cur {
		return ErrInvalidAccount
	}

	if transfer.Amount <= 0 {
		return ErrInvalidAmount
	}

	if accFrom.Amount < transfer.Amount {
		return ErrNotEnoughMoney
	}

	if err = s.repo.Transfer(ctx, transfer.AccountId, transfer.ToAccountId, transfer.Amount); err != nil {
		return err
	}
	return nil
}
