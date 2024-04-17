package repository

import (
	"context"

	"bank-api/internal/domain"
	"bank-api/internal/repository/queries"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type UserRepository interface {
	CreateUser(ctx context.Context, userInfo *domain.UserInfo) (*domain.User, error)
	GetUser(ctx context.Context, id int) (*domain.User, error)
	UserExistsByEmail(ctx context.Context, email string) (bool, error)
	UserExistsById(ctx context.Context, id int) (bool, error)
	GetUserIdByEmail(ctx context.Context, email string) (int, error)
	UpdateUser(ctx context.Context, id int, userInfo *domain.UserInfo) (*domain.User, error)
	DeleteUser(ctx context.Context, id int) error
}

type AccountRepository interface {
	UserExistsById(ctx context.Context, id int) (bool, error)
	GetCurrencyId(ctx context.Context, cur domain.Currency) (int, error)
	CurrencyExists(ctx context.Context, cur domain.Currency) (bool, error)
	AccountExists(ctx context.Context, id int) (bool, error)
	CreateAccount(ctx context.Context, userId int, cur domain.Currency) (*domain.Account, error)
	GetAccount(ctx context.Context, id int) (*domain.Account, error)
	UpdateAccount(ctx context.Context, id int, amount int) (*domain.Account, error)
	DeleteAccount(ctx context.Context, id int) error

	Transaction(ctx context.Context, accountId int, amount int, t domain.TransactionType) error
	Transfer(ctx context.Context, fromAccountId int, toAccountId int, amount int) error

	ListTransactions(ctx context.Context, accountId int) ([]*domain.Transaction, error)
}

type Repository interface {
	UserRepository
	AccountRepository
}
type repo struct {
	*queries.Queries
	pool   *pgxpool.Pool
	logger *zap.SugaredLogger
}

func New(pgxPool *pgxpool.Pool, logger *zap.SugaredLogger) Repository {
	return repo{
		Queries: queries.New(pgxPool),
		pool:    pgxPool,
		logger:  logger,
	}
}
