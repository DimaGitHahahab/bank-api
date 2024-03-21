package queries

import (
	"bank-api/internal/domain"
	"context"
	"fmt"
)

const createUser = `
INSERT INTO "user" (name, email, password)
VALUES ($1, $2, $3) RETURNING id, name, email, password, created_at
`

func (q *Queries) CreateUser(ctx context.Context, newUserInfo *domain.UserInfo) (*domain.User, error) {

	var user domain.User
	err := q.pool.QueryRow(ctx, createUser, newUserInfo.Name, newUserInfo.Email, newUserInfo.Password).Scan(&user.Id, &user.Name, &user.Email, &user.HashedPassword, &user.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("error creating user: %w", err)
	}
	return &user, nil
}

const getUser = `
SELECT id, name, email, password,created_at
FROM "user"
WHERE id = $1
`

func (q *Queries) GetUser(ctx context.Context, id int) (*domain.User, error) {
	var user domain.User
	if err := q.pool.QueryRow(ctx, getUser, id).Scan(&user.Id, &user.Name, &user.Email, &user.HashedPassword, &user.CreatedAt); err != nil {
		return nil, fmt.Errorf("error getting user: %w", err)
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
		return 0, fmt.Errorf("error getting user id by his email: %w", err)
	}
	return userId, nil
}

const ExistsByEmail = `
SELECT EXISTS (
	SELECT 1
	FROM "user"
	WHERE email = $1
)
`

func (q *Queries) UserExistsByEmail(ctx context.Context, email string) (bool, error) {
	var exists bool
	if err := q.pool.QueryRow(ctx, ExistsByEmail, email).Scan(&exists); err != nil {
		return false, fmt.Errorf("error while checking if user exists by email: %w", err)
	}
	return exists, nil
}

const ExistsById = `
SELECT EXISTS (
	SELECT 1
	FROM "user"
	WHERE id = $1
)
`

func (q *Queries) UserExistsById(ctx context.Context, id int) (bool, error) {
	var exists bool
	if err := q.pool.QueryRow(ctx, ExistsById, id).Scan(&exists); err != nil {
		return false, fmt.Errorf("error while checking if user exists by id: %w", err)
	}
	return exists, nil
}

const UpdateUser = `
UPDATE "user"
SET name = $2, email = $3
WHERE id = $1
RETURNING id, name, email, created_at
`

func (q *Queries) UpdateUser(ctx context.Context, id int, userInfo *domain.UserInfo) (*domain.User, error) {
	var user domain.User
	if err := q.pool.QueryRow(ctx, UpdateUser, id, userInfo.Name, userInfo.Email).Scan(&user.Id, &user.Name, &user.Email, &user.CreatedAt); err != nil {
		return nil, fmt.Errorf("error while updating user: %w", err)
	}
	return &user, nil
}

const DeleteUser = `
DELETE FROM "user"
WHERE id = $1
`

func (q *Queries) DeleteUser(ctx context.Context, id int) error {
	if _, err := q.pool.Exec(ctx, DeleteUser, id); err != nil {
		return fmt.Errorf("error while deleting user: %w", err)
	}
	return nil
}
