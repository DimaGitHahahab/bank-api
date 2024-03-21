package service

import (
	"bank-api/internal/domain"
	"bank-api/mocks"
	"context"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
)

func TestCreateAccount(t *testing.T) {
	mockRepo := mocks.NewMockAccountRepository(gomock.NewController(t))

	mockRepo.EXPECT().CurrencyExists(gomock.Any(), domain.Currency{Symbol: "RUB"}).Return(true, nil)
	mockRepo.EXPECT().UserExistsById(gomock.Any(), 1).Return(true, nil)
	mockRepo.EXPECT().GetCurrencyId(gomock.Any(), domain.Currency{Symbol: "RUB"}).Return(1, nil)
	mockRepo.EXPECT().CreateAccount(gomock.Any(), 1, domain.Currency{Id: 1, Symbol: "RUB"}).Return(&domain.Account{Id: 1, UserId: 1, Cur: domain.Currency{Id: 1, Symbol: "RUB"}, Amount: 0}, nil)

	s := NewAccountService(mockRepo)

	account, err := s.CreateAccount(context.Background(), 1, domain.Currency{Symbol: "RUB"})
	assert.NoError(t, err)
	assert.NotNil(t, account)
	assert.Equal(t, 1, account.Id)
	assert.Equal(t, 1, account.UserId)
	assert.Equal(t, domain.Currency{Id: 1, Symbol: "RUB"}, account.Cur)
	assert.Equal(t, 0, account.Amount)
}

func TestCreateAccount_NoSuchCurrency(t *testing.T) {

	mockRepo := mocks.NewMockAccountRepository(gomock.NewController(t))

	mockRepo.EXPECT().CurrencyExists(gomock.Any(), domain.Currency{Symbol: "currencyName"}).Return(false, nil)

	s := NewAccountService(mockRepo)

	account, err := s.CreateAccount(context.Background(), 1, domain.Currency{Symbol: "currencyName"})
	assert.Nil(t, account)
	assert.ErrorIs(t, domain.ErrNoSuchCurrency, err)
}

func TestGetAccount(t *testing.T) {
	mockRepo := mocks.NewMockAccountRepository(gomock.NewController(t))

	mockRepo.EXPECT().AccountExists(gomock.Any(), 1).Return(true, nil)
	mockRepo.EXPECT().UserExistsById(gomock.Any(), 1).Return(true, nil)
	mockRepo.EXPECT().GetAccount(gomock.Any(), 1).Return(&domain.Account{Id: 1, UserId: 1, Cur: domain.Currency{
		Id:     1,
		Symbol: "RUB",
	}, Amount: 0}, nil)

	s := NewAccountService(mockRepo)

	account, err := s.GetAccount(context.Background(), 1, 1)
	assert.NoError(t, err)
	assert.NotNil(t, account)
	assert.Equal(t, 1, account.Id)
	assert.Equal(t, 1, account.UserId)
	assert.Equal(t, domain.Currency{
		Id:     1,
		Symbol: "RUB",
	}, account.Cur)
	assert.Equal(t, 0, account.Amount)

}

func TestGetAccount_WrongUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockAccountRepository(ctrl)

	mockRepo.EXPECT().AccountExists(gomock.Any(), 1).Return(true, nil)
	mockRepo.EXPECT().UserExistsById(gomock.Any(), 1).Return(true, nil)
	mockRepo.EXPECT().GetAccount(gomock.Any(), 1).Return(&domain.Account{Id: 1, UserId: 2, Cur: domain.Currency{
		Id:     1,
		Symbol: "RUB",
	}, Amount: 0}, nil)

	s := NewAccountService(mockRepo)

	account, err := s.GetAccount(context.Background(), 1, 1)
	assert.Nil(t, account)
	assert.ErrorIs(t, domain.ErrInvalidAccount, err)
}

func TestUpdateAccount(t *testing.T) {
	mockRepo := mocks.NewMockAccountRepository(gomock.NewController(t))

	mockRepo.EXPECT().AccountExists(gomock.Any(), 1).Return(true, nil)
	mockRepo.EXPECT().UserExistsById(gomock.Any(), 1).Return(true, nil)
	mockRepo.EXPECT().GetAccount(gomock.Any(), 1).Return(&domain.Account{Id: 1, UserId: 1, Cur: domain.Currency{
		Id:     1,
		Symbol: "RUB",
	}, Amount: 0}, nil)

	mockRepo.EXPECT().UpdateAccount(gomock.Any(), 1, 100).Return(&domain.Account{Id: 1, UserId: 1, Cur: domain.Currency{
		Id:     1,
		Symbol: "RUB",
	}, Amount: 100}, nil)

	s := NewAccountService(mockRepo)

	account, err := s.UpdateAccount(context.Background(), 1, 1, 100)
	assert.NotNil(t, account)
	assert.NoError(t, err)
	assert.Equal(t, 1, account.Id)
	assert.Equal(t, 1, account.UserId)
	assert.Equal(t, domain.Currency{
		Id:     1,
		Symbol: "RUB",
	}, account.Cur)
	assert.Equal(t, 100, account.Amount)
}

func TestUpdateAccount_WrongUser(t *testing.T) {
	mockRepo := mocks.NewMockAccountRepository(gomock.NewController(t))

	mockRepo.EXPECT().AccountExists(gomock.Any(), 1).Return(true, nil)
	mockRepo.EXPECT().UserExistsById(gomock.Any(), 1).Return(true, nil)
	mockRepo.EXPECT().GetAccount(gomock.Any(), 1).Return(&domain.Account{Id: 1, UserId: 2, Cur: domain.Currency{
		Id:     1,
		Symbol: "RUB",
	}, Amount: 0}, nil)

	s := NewAccountService(mockRepo)

	account, err := s.UpdateAccount(context.Background(), 1, 1, 100)
	assert.Nil(t, account)
	assert.ErrorIs(t, domain.ErrInvalidAccount, err)
}

func TestDeleteAccount(t *testing.T) {
	mockRepo := mocks.NewMockAccountRepository(gomock.NewController(t))

	mockRepo.EXPECT().UserExistsById(gomock.Any(), 1).Return(true, nil)
	mockRepo.EXPECT().AccountExists(gomock.Any(), 1).Return(true, nil)
	mockRepo.EXPECT().GetAccount(gomock.Any(), 1).Return(&domain.Account{Id: 1, UserId: 1, Cur: domain.Currency{
		Id:     1,
		Symbol: "RUB",
	}, Amount: 0}, nil)
	mockRepo.EXPECT().DeleteAccount(gomock.Any(), 1).Return(nil)

	s := NewAccountService(mockRepo)

	err := s.DeleteAccount(context.Background(), 1, 1)
	assert.NoError(t, err)
}

func TestDeleteAccount_WrongUser(t *testing.T) {
	mockRepo := mocks.NewMockAccountRepository(gomock.NewController(t))

	mockRepo.EXPECT().UserExistsById(gomock.Any(), 1).Return(true, nil)
	mockRepo.EXPECT().AccountExists(gomock.Any(), 1).Return(true, nil)
	mockRepo.EXPECT().GetAccount(gomock.Any(), 1).Return(&domain.Account{Id: 1, UserId: 2, Cur: domain.Currency{
		Id:     1,
		Symbol: "RUB",
	}, Amount: 0}, nil)

	s := NewAccountService(mockRepo)

	err := s.DeleteAccount(context.Background(), 1, 1)
	assert.ErrorIs(t, domain.ErrInvalidAccount, err)
}
