package bank

import (
	"bank-api/internal/model"
	"bank-api/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestProcessTransaction_Deposit(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockAccountRepository(ctrl)

	mockRepo.EXPECT().GetAccount(gomock.Any(), 1).Return(&model.Account{Id: 1, UserId: 1, Cur: model.Currency{
		Id:     1,
		Symbol: "RUB",
	}, Amount: 0}, nil)
	mockRepo.EXPECT().Transaction(gomock.Any(), 1, 100, model.Deposit).Return(nil)

	s := NewTransactionService(mockRepo)

	transaction := &model.Transaction{
		ToAccountId: 1,
		UserId:      1,
		Amount:      100,
		Type:        model.Deposit,
	}

	err := s.ProcessTransaction(nil, transaction)
	assert.NoError(t, err)
}

func TestProcessTransaction_Deposit_InvalidAmount(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockAccountRepository(ctrl)

	s := NewTransactionService(mockRepo)

	transaction := &model.Transaction{
		ToAccountId: 1,
		UserId:      1,
		Amount:      -100,
		Type:        model.Deposit,
	}

	err := s.ProcessTransaction(nil, transaction)
	assert.Error(t, err)

	transaction.Amount = 0
	err = s.ProcessTransaction(nil, transaction)
	assert.Error(t, err)
}

func TestProcessTransaction_Deposit_InvalidAccount(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockAccountRepository(ctrl)

	mockRepo.EXPECT().GetAccount(gomock.Any(), 1).Return(nil, ErrInvalidAccount)

	s := NewTransactionService(mockRepo)

	transaction := &model.Transaction{
		ToAccountId: 1,
		UserId:      1,
		Amount:      100,
		Type:        model.Deposit,
	}

	err := s.ProcessTransaction(nil, transaction)
	assert.Error(t, err)
}

func TestProcessTransaction_Withdraw(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockAccountRepository(ctrl)

	mockRepo.EXPECT().GetAccount(gomock.Any(), 1).Return(&model.Account{Id: 1, UserId: 1, Cur: model.Currency{
		Id:     1,
		Symbol: "RUB",
	}, Amount: 400}, nil)
	mockRepo.EXPECT().Transaction(gomock.Any(), 1, 200, model.Withdraw).Return(nil)

	s := NewTransactionService(mockRepo)

	transaction := &model.Transaction{
		FromAccountId: 1,
		UserId:        1,
		Amount:        200,
		Type:          model.Withdraw,
	}

	err := s.ProcessTransaction(nil, transaction)
	assert.NoError(t, err)
}

func TestProcessTransaction_Withdraw_InvalidAmount(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockAccountRepository(ctrl)

	s := NewTransactionService(mockRepo)

	transaction := &model.Transaction{
		FromAccountId: 1,
		UserId:        1,
		Amount:        -100,
		Type:          model.Withdraw,
	}

	err := s.ProcessTransaction(nil, transaction)
	assert.Error(t, err)

	transaction.Amount = 0
	err = s.ProcessTransaction(nil, transaction)
	assert.Error(t, err)
}

func TestProcessTransfer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockAccountRepository(ctrl)

	mockRepo.EXPECT().GetAccount(gomock.Any(), 1).Return(&model.Account{Id: 1, UserId: 1, Cur: model.Currency{
		Id:     1,
		Symbol: "RUB",
	}, Amount: 100}, nil)
	mockRepo.EXPECT().GetAccount(gomock.Any(), 2).Return(&model.Account{Id: 2, UserId: 2, Cur: model.Currency{
		Id:     1,
		Symbol: "RUB",
	}, Amount: 200}, nil)
	mockRepo.EXPECT().Transfer(gomock.Any(), 1, 2, 50).Return(nil)

	s := NewTransactionService(mockRepo)

	transaction := &model.Transaction{
		FromAccountId: 1,
		ToAccountId:   2,
		UserId:        1,
		Amount:        50,
		Type:          model.Transfer,
	}

	err := s.ProcessTransaction(nil, transaction)
	assert.NoError(t, err)
}

func TestListTransactions(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockAccountRepository(ctrl)

	mockRepo.EXPECT().GetAccount(gomock.Any(), 1).Return(&model.Account{Id: 1, UserId: 1, Cur: model.Currency{
		Id:     1,
		Symbol: "RUB",
	}, Amount: 100}, nil)

	trTimes := []time.Time{
		time.Now().Add(-time.Hour),
		time.Now(),
		time.Now().Add(time.Hour),
		time.Now().Add(2 * time.Hour),
	}
	trs := []*model.Transaction{
		{
			FromAccountId: 1,
			ToAccountId:   2,
			UserId:        1,
			Amount:        50,
			Type:          model.Transfer,
			Time:          trTimes[0],
		},
		{
			FromAccountId: 3,
			ToAccountId:   1,
			UserId:        2,
			Amount:        100,
			Type:          model.Transfer,
			Time:          trTimes[1],
		},
		{
			FromAccountId: 1,
			ToAccountId:   0,
			UserId:        1,
			Amount:        100,
			Type:          model.Withdraw,
			Time:          trTimes[2],
		},
		{
			FromAccountId: 0,
			ToAccountId:   1,
			UserId:        1,
			Amount:        100,
			Type:          model.Deposit,
			Time:          trTimes[3],
		},
	}
	mockRepo.EXPECT().ListTransactions(gomock.Any(), 1).Return(trs, nil)

	s := NewTransactionService(mockRepo)

	transactions, err := s.ListTransactions(nil, 1)
	assert.NoError(t, err)
	assert.NotNil(t, transactions)
	assert.Equal(t, 4, len(transactions))

	for i := 0; i < len(transactions); i++ {
		assertTransaction(t, trs[i], transactions[i])
	}
}

func assertTransaction(t *testing.T, expected *model.Transaction, got *model.Transaction) {
	assert.Equal(t, expected.FromAccountId, got.FromAccountId)
	assert.Equal(t, expected.ToAccountId, got.ToAccountId)
	assert.Equal(t, expected.UserId, got.UserId)
	assert.Equal(t, expected.Amount, got.Amount)
	assert.Equal(t, expected.Type, got.Type)
	assert.Equal(t, expected.Time, got.Time)
}
