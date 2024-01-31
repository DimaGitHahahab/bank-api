package bank

import (
	"bank-api/internal/model"
	"bank-api/mocks"
	"context"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreateAccount(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockAccountRepository(ctrl)

	mockRepo.EXPECT().GetCurrencyId(gomock.Any(), model.Currency{Symbol: "RUB"}).Return(1, nil)
	mockRepo.EXPECT().CreateAccount(gomock.Any(), 1, model.Currency{
		Id:     1,
		Symbol: "RUB",
	}).Return(&model.Account{Id: 1, UserId: 1, Cur: model.Currency{
		Id:     1,
		Symbol: "RUB",
	}, Amount: 0}, nil)

	s := NewAccountService(mockRepo)

	account, err := s.CreateAccount(context.Background(), 1, model.Currency{Symbol: "RUB"})
	assert.NoError(t, err)
	assert.NotNil(t, account)
	assert.Equal(t, 1, account.Id)
	assert.Equal(t, 1, account.UserId)
	assert.Equal(t, model.Currency{Id: 1, Symbol: "RUB"}, account.Cur)
	assert.Equal(t, 0, account.Amount)
}

func TestCreateAccount_NoSuchCurrency(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockAccountRepository(ctrl)

	mockRepo.EXPECT().GetCurrencyId(gomock.Any(), model.Currency{Symbol: "currencyName"}).Return(0, ErrNoSuchCurrency)

	s := NewAccountService(mockRepo)

	account, err := s.CreateAccount(context.Background(), 1, model.Currency{Symbol: "currencyName"})
	assert.Error(t, err)
	assert.Nil(t, account)
}

func TestGetAccount(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockAccountRepository(ctrl)

	mockRepo.EXPECT().GetAccount(gomock.Any(), 1).Return(&model.Account{Id: 1, UserId: 1, Cur: model.Currency{
		Id:     1,
		Symbol: "RUB",
	}, Amount: 0}, nil)

	s := NewAccountService(mockRepo)

	account, err := s.GetAccount(context.Background(), 1, 1)
	assert.NoError(t, err)
	assert.NotNil(t, account)
	assert.Equal(t, 1, account.Id)
	assert.Equal(t, 1, account.UserId)
	assert.Equal(t, model.Currency{
		Id:     1,
		Symbol: "RUB",
	}, account.Cur)
	assert.Equal(t, 0, account.Amount)
}

func TestGetAccount_WrongUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockAccountRepository(ctrl)

	mockRepo.EXPECT().GetAccount(gomock.Any(), 1).Return(&model.Account{Id: 1, UserId: 2, Cur: model.Currency{
		Id:     1,
		Symbol: "RUB",
	}, Amount: 0}, nil)

	s := NewAccountService(mockRepo)

	account, err := s.GetAccount(context.Background(), 1, 1)
	assert.Error(t, err)
	assert.Nil(t, account)
}

func TestUpdateAccount(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockAccountRepository(ctrl)

	mockRepo.EXPECT().GetAccount(gomock.Any(), 1).Return(&model.Account{Id: 1, UserId: 1, Cur: model.Currency{
		Id:     1,
		Symbol: "RUB",
	}, Amount: 0}, nil)
	mockRepo.EXPECT().UpdateAccount(gomock.Any(), 1, 100).Return(&model.Account{Id: 1, UserId: 1, Cur: model.Currency{
		Id:     1,
		Symbol: "RUB",
	}, Amount: 100}, nil)

	s := NewAccountService(mockRepo)

	account, err := s.UpdateAccount(context.Background(), 1, 1, 100)
	assert.NoError(t, err)
	assert.NotNil(t, account)
	assert.Equal(t, 1, account.Id)
	assert.Equal(t, 1, account.UserId)
	assert.Equal(t, model.Currency{
		Id:     1,
		Symbol: "RUB",
	}, account.Cur)
	assert.Equal(t, 100, account.Amount)
}

func TestUpdateAccount_WrongUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockAccountRepository(ctrl)
	mockRepo.EXPECT().GetAccount(gomock.Any(), 1).Return(&model.Account{Id: 1, UserId: 2, Cur: model.Currency{
		Id:     1,
		Symbol: "RUB",
	}, Amount: 0}, nil)

	s := NewAccountService(mockRepo)

	account, err := s.UpdateAccount(context.Background(), 1, 1, 100)
	assert.Error(t, err)
	assert.Nil(t, account)
}

func TestDeleteAccount(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockAccountRepository(ctrl)

	//mockRepo.EXPECT().GetAccount(gomock.Any(), 1).Return(&model.Account{Id: 1, UserId: 1, Cur: "RUB", Amount: 0}, nil)
	mockRepo.EXPECT().GetAccount(gomock.Any(), 1).Return(&model.Account{Id: 1, UserId: 1, Cur: model.Currency{
		Id:     1,
		Symbol: "RUB",
	}, Amount: 0}, nil)

	mockRepo.EXPECT().DeleteAccount(gomock.Any(), 1).Return(nil)

	s := NewAccountService(mockRepo)

	err := s.DeleteAccount(context.Background(), 1, 1)
	assert.NoError(t, err)
}

func TestDeleteAccount_WrongUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockAccountRepository(ctrl)

	mockRepo.EXPECT().GetAccount(gomock.Any(), 1).Return(&model.Account{Id: 1, UserId: 2, Cur: model.Currency{
		Id:     1,
		Symbol: "RUB",
	}, Amount: 0}, nil)

	s := NewAccountService(mockRepo)

	err := s.DeleteAccount(context.Background(), 1, 1)
	assert.Error(t, err)
}
