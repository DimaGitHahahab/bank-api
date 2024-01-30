package repository

import (
	"bank-api/internal/model"
	"bank-api/internal/repository/queries"
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type UserRepository interface {
	CreateUser(ctx context.Context, userInfo *model.UserInfo) (*model.User, error)
	GetUser(ctx context.Context, id int) (*model.User, error)
	GetUserIdByEmail(ctx context.Context, email string) (int, error)
	UpdateUser(ctx context.Context, id int, userInfo *model.UserInfo) (*model.User, error)
	DeleteUser(ctx context.Context, id int) error
}

type AccountRepository interface {
	GetCurrencyId(ctx context.Context, cur model.Currency) (int, error)
	CreateAccount(ctx context.Context, userId int, cur model.Currency) (*model.Account, error)
	GetAccount(ctx context.Context, id int) (*model.Account, error)
	UpdateAccount(ctx context.Context, id int, amount int) (*model.Account, error)
	DeleteAccount(ctx context.Context, id int) error

	Transaction(ctx context.Context, accountId int, amount int) error
	Transfer(ctx context.Context, fromAccountId int, toAccountId int, amount int) error
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
