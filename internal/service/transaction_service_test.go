package service

import (
	"bank-api/internal/domain"
	"bank-api/mocks"
	"context"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
	"time"
)

func TestProcessTransaction_Deposit(t *testing.T) {
	mockRepo := mocks.NewMockAccountRepository(gomock.NewController(t))

	mockRepo.EXPECT().AccountExists(gomock.Any(), 1).Return(true, nil)
	mockRepo.EXPECT().GetAccount(gomock.Any(), 1).Return(&domain.Account{Id: 1, UserId: 1, Cur: domain.Currency{
		Id:     1,
		Symbol: "RUB",
	}, Amount: 100}, nil)
	mockRepo.EXPECT().Transaction(gomock.Any(), 1, 100, domain.Deposit).Return(nil)

	s := NewTransactionService(mockRepo)

	transaction := &domain.Transaction{
		ToAccountId: 1,
		UserId:      1,
		Amount:      100,
		Type:        domain.Deposit,
	}

	err := s.ProcessTransaction(context.Background(), transaction)
	assert.NoError(t, err)
}

func TestProcessTransaction_Deposit_InvalidAmount(t *testing.T) {
	mockRepo := mocks.NewMockAccountRepository(gomock.NewController(t))

	s := NewTransactionService(mockRepo)

	transaction := &domain.Transaction{
		ToAccountId: 1,
		UserId:      1,
		Amount:      -100,
		Type:        domain.Deposit,
	}

	err := s.ProcessTransaction(context.Background(), transaction)
	assert.Error(t, err)

	transaction.Amount = 0
	err = s.ProcessTransaction(context.Background(), transaction)
	assert.Error(t, err)
}

func TestProcessTransaction_Deposit_InvalidAccount(t *testing.T) {
	mockRepo := mocks.NewMockAccountRepository(gomock.NewController(t))

	mockRepo.EXPECT().AccountExists(gomock.Any(), 1).Return(false, nil)

	s := NewTransactionService(mockRepo)

	transaction := &domain.Transaction{
		ToAccountId: 1,
		UserId:      1,
		Amount:      100,
		Type:        domain.Deposit,
	}

	err := s.ProcessTransaction(context.Background(), transaction)
	assert.ErrorIs(t, domain.ErrNoSuchAccount, err)
}

func TestProcessTransaction_Withdraw(t *testing.T) {
	mockRepo := mocks.NewMockAccountRepository(gomock.NewController(t))

	mockRepo.EXPECT().AccountExists(gomock.Any(), 1).Return(true, nil)
	mockRepo.EXPECT().GetAccount(gomock.Any(), 1).Return(&domain.Account{Id: 1, UserId: 1, Cur: domain.Currency{
		Id:     1,
		Symbol: "RUB",
	}, Amount: 200}, nil)
	mockRepo.EXPECT().Transaction(gomock.Any(), 1, 200, domain.Withdraw).Return(nil)

	s := NewTransactionService(mockRepo)

	transaction := &domain.Transaction{
		FromAccountId: 1,
		UserId:        1,
		Amount:        200,
		Type:          domain.Withdraw,
	}

	err := s.ProcessTransaction(context.Background(), transaction)
	assert.NoError(t, err)
}

func TestProcessTransaction_Withdraw_InvalidAmount(t *testing.T) {
	mockRepo := mocks.NewMockAccountRepository(gomock.NewController(t))

	s := NewTransactionService(mockRepo)

	transaction := &domain.Transaction{
		FromAccountId: 1,
		UserId:        1,
		Amount:        -100,
		Type:          domain.Withdraw,
	}

	err := s.ProcessTransaction(context.Background(), transaction)
	assert.ErrorIs(t, domain.ErrInvalidAmount, err)

	transaction.Amount = 0
	err = s.ProcessTransaction(context.Background(), transaction)
	assert.ErrorIs(t, domain.ErrInvalidAmount, err)
}

func TestProcessTransfer(t *testing.T) {
	mockRepo := mocks.NewMockAccountRepository(gomock.NewController(t))

	mockRepo.EXPECT().AccountExists(gomock.Any(), 1).Return(true, nil)
	mockRepo.EXPECT().GetAccount(gomock.Any(), 2).Return(&domain.Account{Id: 1, UserId: 1, Cur: domain.Currency{
		Id:     1,
		Symbol: "RUB",
	}, Amount: 100}, nil)
	mockRepo.EXPECT().AccountExists(gomock.Any(), 2).Return(true, nil)
	mockRepo.EXPECT().Transfer(gomock.Any(), 1, 2, 50).Return(nil)

	s := NewTransactionService(mockRepo)

	transaction := &domain.Transaction{
		FromAccountId: 1,
		ToAccountId:   2,
		UserId:        1,
		Amount:        50,
		Type:          domain.Transfer,
	}

	err := s.ProcessTransaction(context.Background(), transaction)
	assert.NoError(t, err)
}

func TestListTransactions(t *testing.T) {
	mockRepo := mocks.NewMockAccountRepository(gomock.NewController(t))

	mockRepo.EXPECT().UserExistsById(gomock.Any(), 1).Return(true, nil)
	trTimes := []time.Time{
		time.Now().Add(-time.Hour),
		time.Now(),
		time.Now().Add(time.Hour),
		time.Now().Add(2 * time.Hour),
	}
	trs := []*domain.Transaction{
		{
			FromAccountId: 1,
			ToAccountId:   2,
			UserId:        1,
			Amount:        50,
			Type:          domain.Transfer,
			Time:          trTimes[0],
		},
		{
			FromAccountId: 3,
			ToAccountId:   1,
			UserId:        2,
			Amount:        100,
			Type:          domain.Transfer,
			Time:          trTimes[1],
		},
		{
			FromAccountId: 1,
			ToAccountId:   0,
			UserId:        1,
			Amount:        100,
			Type:          domain.Withdraw,
			Time:          trTimes[2],
		},
		{
			FromAccountId: 0,
			ToAccountId:   1,
			UserId:        1,
			Amount:        100,
			Type:          domain.Deposit,
			Time:          trTimes[3],
		},
	}
	mockRepo.EXPECT().ListTransactions(gomock.Any(), 1).Return(trs, nil)

	s := NewTransactionService(mockRepo)

	transactions, err := s.ListTransactions(context.Background(), 1)
	assert.NoError(t, err)
	assert.NotNil(t, transactions)
	assert.Equal(t, 4, len(transactions))

	for i := 0; i < len(transactions); i++ {
		assertTransaction(t, trs[i], transactions[i])
	}
}

func assertTransaction(t *testing.T, expected *domain.Transaction, got *domain.Transaction) {
	assert.Equal(t, expected.FromAccountId, got.FromAccountId)
	assert.Equal(t, expected.ToAccountId, got.ToAccountId)
	assert.Equal(t, expected.UserId, got.UserId)
	assert.Equal(t, expected.Amount, got.Amount)
	assert.Equal(t, expected.Type, got.Type)
	assert.Equal(t, expected.Time, got.Time)
}
