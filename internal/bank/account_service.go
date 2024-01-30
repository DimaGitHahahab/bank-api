package bank

import (
	"bank-api/internal/model"
	"bank-api/internal/repository"
	"context"
	"errors"
)

var (
	ErrInvalidAccount = errors.New("invalid account")
	ErrNoSuchCurrency = errors.New("no such currency")
)

type AccountService interface {
	CreateAccount(ctx context.Context, userId int, cur model.Currency) (*model.Account, error)
	GetAccount(ctx context.Context, userId int, accountId int) (*model.Account, error)
	UpdateAccount(ctx context.Context, userId int, accountId int, amount int) (*model.Account, error)
	DeleteAccount(ctx context.Context, userId int, accountId int) error
}

type accountService struct {
	repo repository.AccountRepository
}

func NewAccountService(repo repository.AccountRepository) AccountService {
	return &accountService{repo: repo}
}

func (s *accountService) CreateAccount(ctx context.Context, userId int, cur model.Currency) (*model.Account, error) {
	_, err := s.repo.GetCurrencyId(ctx, cur)
	if err != nil {
		return nil, ErrNoSuchCurrency
	}
	account, err := s.repo.CreateAccount(ctx, userId, cur)
	if err != nil {
		return nil, err
	}

	return account, nil
}

func (s *accountService) GetAccount(ctx context.Context, userId int, accountId int) (*model.Account, error) {
	account, err := s.repo.GetAccount(ctx, accountId)
	if err != nil {
		return nil, err
	}
	if account.UserId != userId {
		return nil, ErrInvalidAccount
	}
	return account, nil
}

func (s *accountService) UpdateAccount(ctx context.Context, userId int, accountId int, amount int) (*model.Account, error) {
	account, err := s.repo.GetAccount(ctx, accountId)
	if err != nil {
		return nil, err
	}
	if account.UserId != userId {
		return nil, ErrInvalidAccount
	}
	account, err = s.repo.UpdateAccount(ctx, accountId, amount)
	if err != nil {
		return nil, err
	}

	return account, nil
}

func (s *accountService) DeleteAccount(ctx context.Context, userId int, accountId int) error {
	acc, err := s.repo.GetAccount(ctx, accountId)
	if err != nil {
		return err
	}
	if acc.UserId != userId {
		return ErrInvalidAccount
	}

	err = s.repo.DeleteAccount(ctx, accountId)
	if err != nil {
		return err
	}

	return nil
}
