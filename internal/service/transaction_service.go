package service

import (
	"bank-api/internal/domain"
	"bank-api/internal/repository"
	"context"
	"fmt"
)

type TransactionService interface {
	ProcessTransaction(ctx context.Context, transaction *domain.Transaction) error
	ListTransactions(ctx context.Context, accountId int) ([]*domain.Transaction, error)
}

type transactionService struct {
	repo repository.AccountRepository
}

func NewTransactionService(repo repository.AccountRepository) TransactionService {
	return &transactionService{repo: repo}
}

func (s *transactionService) ProcessTransaction(ctx context.Context, transaction *domain.Transaction) error {
	if transaction.Amount <= 0 {
		return domain.ErrInvalidAmount
	}

	switch transaction.Type {
	case domain.Deposit:
		return s.processDeposit(ctx, transaction)
	case domain.Withdraw:
		return s.processWithdraw(ctx, transaction)
	case domain.Transfer:
		return s.processTransfer(ctx, transaction)
	}

	return nil
}

func (s *transactionService) processDeposit(ctx context.Context, transaction *domain.Transaction) error {
	ok, err := s.repo.AccountExists(ctx, transaction.ToAccountId)
	if err != nil {
		return fmt.Errorf("can't check if such an account exists: %w", err)
	}
	if !ok {
		return domain.ErrNoSuchAccount
	}

	accTo, err := s.repo.GetAccount(ctx, transaction.ToAccountId)
	if err != nil {
		return fmt.Errorf("can't get account: %w", err)
	}
	if accTo.UserId != transaction.UserId {
		return domain.ErrInvalidAccount
	}

	if err := s.repo.Transaction(ctx, transaction.ToAccountId, transaction.Amount, transaction.Type); err != nil {
		return fmt.Errorf("can't perform transaction: %w", err)
	}

	return nil
}

func (s *transactionService) processWithdraw(ctx context.Context, transaction *domain.Transaction) error {
	ok, err := s.repo.AccountExists(ctx, transaction.FromAccountId)
	if err != nil {
		return fmt.Errorf("can't check if such an account exists: %w", err)
	}
	if !ok {
		return domain.ErrNoSuchAccount
	}

	accFrom, err := s.repo.GetAccount(ctx, transaction.FromAccountId)
	if err != nil {
		return fmt.Errorf("can't get account: %w", err)
	}
	if accFrom.UserId != transaction.UserId {
		return domain.ErrInvalidAccount
	}

	if accFrom.Amount < transaction.Amount {
		return domain.ErrNotEnoughMoney
	}

	if err := s.repo.Transaction(ctx, transaction.FromAccountId, transaction.Amount, transaction.Type); err != nil {
		return fmt.Errorf("can't process transaction: %w", err)
	}

	return nil
}

func (s *transactionService) processTransfer(ctx context.Context, transaction *domain.Transaction) error {
	ok, err := s.repo.AccountExists(ctx, transaction.FromAccountId)
	if err != nil {
		return fmt.Errorf("can't check if such an account exists: %w", err)
	}
	if !ok {
		return domain.ErrNoSuchAccount
	}

	accFrom, err := s.repo.GetAccount(ctx, transaction.ToAccountId)
	if err != nil {
		return fmt.Errorf("can't get account: %w", err)
	}
	if accFrom.UserId != transaction.UserId {
		return domain.ErrInvalidAccount
	}

	if accFrom.Amount < transaction.Amount {
		return domain.ErrNotEnoughMoney
	}

	ok, err = s.repo.AccountExists(ctx, transaction.ToAccountId)
	if err != nil {
		return fmt.Errorf("can't check if such an account exists: %w", err)
	}
	if !ok {
		return domain.ErrNoSuchAccount
	}

	if err := s.repo.Transfer(ctx, transaction.FromAccountId, transaction.ToAccountId, transaction.Amount); err != nil {
		return fmt.Errorf("can't process transaction: %w", err)
	}

	return nil
}

func (s *transactionService) ListTransactions(ctx context.Context, userId int) ([]*domain.Transaction, error) {
	ok, err := s.repo.UserExistsById(ctx, userId)
	if err != nil {
		return nil, fmt.Errorf("can't check if such a user exists: %w", err)
	}
	if !ok {
		return nil, domain.ErrNoSuchUser
	}

	trs, err := s.repo.ListTransactions(ctx, userId)
	if err != nil {
		return nil, fmt.Errorf("can't list transactions: %w", err)
	}

	return trs, nil
}
