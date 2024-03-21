package service

import (
	"bank-api/internal/domain"
	"bank-api/internal/repository"
	"context"
	"fmt"
)

type AccountService interface {
	CreateAccount(ctx context.Context, userId int, cur domain.Currency) (*domain.Account, error)
	GetAccount(ctx context.Context, userId int, accountId int) (*domain.Account, error)
	UpdateAccount(ctx context.Context, userId int, accountId int, amount int) (*domain.Account, error)
	DeleteAccount(ctx context.Context, userId int, accountId int) error
}

type accountService struct {
	repo repository.AccountRepository
}

func NewAccountService(repo repository.AccountRepository) AccountService {
	return &accountService{repo: repo}
}

func (s *accountService) CreateAccount(ctx context.Context, userId int, cur domain.Currency) (*domain.Account, error) {
	ok, err := s.repo.CurrencyExists(ctx, cur)
	if err != nil {
		return nil, fmt.Errorf("can't check if such currency exists: %w", err)
	}
	if !ok {
		return nil, domain.ErrNoSuchCurrency
	}

	ok, err = s.repo.UserExistsById(ctx, userId)
	if err != nil {
		return nil, fmt.Errorf("can't check if such a user exists: %w", err)
	}
	if !ok {
		return nil, domain.ErrNoSuchUser
	}

	id, err := s.repo.GetCurrencyId(ctx, cur)
	if err != nil {
		return nil, fmt.Errorf("can't get currency by id: %w", err)
	}

	cur.Id = id
	account, err := s.repo.CreateAccount(ctx, userId, cur)
	if err != nil {
		return nil, fmt.Errorf("can't create account: %w", err)
	}

	return account, nil
}

func (s *accountService) GetAccount(ctx context.Context, userId int, accountId int) (*domain.Account, error) {
	ok, err := s.repo.AccountExists(ctx, accountId)
	if err != nil {
		return nil, fmt.Errorf("can't check if account exists: %w", err)
	}
	if !ok {
		return nil, domain.ErrNoSuchAccount
	}

	ok, err = s.repo.UserExistsById(ctx, userId)
	if err != nil {
		return nil, fmt.Errorf("can't check if such a user exists: %w", err)
	}
	if !ok {
		return nil, domain.ErrNoSuchUser
	}

	account, err := s.repo.GetAccount(ctx, accountId)
	if err != nil {
		return nil, fmt.Errorf("can't get account: %w", err)
	}

	if account.UserId != userId {
		return nil, domain.ErrInvalidAccount
	}

	return account, nil
}

func (s *accountService) UpdateAccount(ctx context.Context, userId int, accountId int, amount int) (*domain.Account, error) {
	ok, err := s.repo.AccountExists(ctx, accountId)
	if err != nil {
		return nil, fmt.Errorf("can't check if account exists: %w", err)
	}
	if !ok {
		return nil, domain.ErrNoSuchAccount
	}

	ok, err = s.repo.UserExistsById(ctx, userId)
	if err != nil {
		return nil, fmt.Errorf("can't check if such a user exists: %w", err)
	}
	if !ok {
		return nil, domain.ErrNoSuchUser
	}

	account, err := s.repo.GetAccount(ctx, accountId)
	if err != nil {
		return nil, fmt.Errorf("can't get account by id: %w", err)
	}

	if account.UserId != userId {
		return nil, domain.ErrInvalidAccount
	}
	account, err = s.repo.UpdateAccount(ctx, accountId, amount)
	if err != nil {
		return nil, fmt.Errorf("can't update account: %w", err)
	}

	return account, nil
}

func (s *accountService) DeleteAccount(ctx context.Context, userId int, accountId int) error {
	ok, err := s.repo.UserExistsById(ctx, userId)
	if err != nil {
		return fmt.Errorf("can't check if such a user exists")
	}
	if !ok {
		return domain.ErrNoSuchUser
	}

	ok, err = s.repo.AccountExists(ctx, accountId)
	if err != nil {
		return fmt.Errorf("can't check if such an account exists: %w", err)
	}
	if !ok {
		return domain.ErrNoSuchAccount
	}

	account, err := s.repo.GetAccount(ctx, accountId)
	if err != nil {
		return fmt.Errorf("can't get account by id: %w", err)
	}

	if userId != account.UserId {
		return domain.ErrInvalidAccount
	}

	err = s.repo.DeleteAccount(ctx, accountId)
	if err != nil {
		return fmt.Errorf("can't delete account by id: %w", err)
	}

	return nil
}
