package queries

import (
	"bank-api/internal/model"
	"context"
	"errors"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
)

var (
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrNoSuchUser        = errors.New("no such user")
)

const createUser = `
INSERT INTO "user" (name, email, password)
VALUES ($1, $2, $3) RETURNING id, name, email, password, created_at
`

func (q *Queries) CreateUser(ctx context.Context, newUserInfo *model.UserInfo) (*model.User, error) {

	var user model.User
	err := q.pool.QueryRow(ctx, createUser, newUserInfo.Name, newUserInfo.Email, newUserInfo.Password).Scan(&user.Id, &user.Name, &user.Email, &user.HashedPassword, &user.CreatedAt)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			return nil, ErrUserAlreadyExists
		}
	}
	return &user, nil
}

const getUser = `
SELECT id, name, email, password,created_at
FROM "user"
WHERE id = $1
`

func (q *Queries) GetUser(ctx context.Context, id int) (*model.User, error) {
	var user model.User
	if err := q.pool.QueryRow(ctx, getUser, id).Scan(&user.Id, &user.Name, &user.Email, &user.HashedPassword, &user.CreatedAt); err != nil {
		if err != nil {
			return nil, ErrNoSuchUser
		}
	}
	return &user, nil
}

const getUserIdByEmail = `
SELECT id
FROM "user"
WHERE email = $1
`

func (q *Queries) GetUserIdByEmail(ctx context.Context, email string) (int, error) {
	var userId int
	if err := q.pool.QueryRow(ctx, getUserIdByEmail, email).Scan(&userId); err != nil {
		return 0, ErrNoSuchUser
	}
	return userId, nil
}

const UpdateUser = `
UPDATE "user"
SET name = $2, email = $3
WHERE id = $1
RETURNING id, name, email, created_at
`

func (q *Queries) UpdateUser(ctx context.Context, id int, userInfo *model.UserInfo) (*model.User, error) {
	var user model.User
	if err := q.pool.QueryRow(ctx, UpdateUser, id, userInfo.Name, userInfo.Email).Scan(&user.Id, &user.Name, &user.Email, &user.CreatedAt); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			return nil, ErrUserAlreadyExists
		}
		return nil, ErrNoSuchUser
	}
	return &user, nil
}

const DeleteUser = `
DELETE FROM "user"
WHERE id = $1
`

func (q *Queries) DeleteUser(ctx context.Context, id int) error {
	if _, err := q.pool.Exec(ctx, DeleteUser, id); err != nil {
		return ErrNoSuchUser
	}
	return nil
}
