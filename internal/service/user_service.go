package service

import (
	"context"
	"fmt"
	"net/mail"

	"bank-api/internal/domain"
	"bank-api/internal/repository"

	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	CreateUser(ctx context.Context, info *domain.UserInfo) (*domain.User, error)
	GetUserById(ctx context.Context, id int) (*domain.User, error)
	GetUserByEmail(ctx context.Context, email string) (*domain.User, error)
	UpdateUserInfo(ctx context.Context, id int, info *domain.UserInfo) (*domain.User, error)
	DeleteUserById(ctx context.Context, id int) error
	AuthenticateUser(ctx context.Context, u *domain.UserInfo) (*domain.User, error)
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
}

func (s *userService) CreateUser(ctx context.Context, new *domain.UserInfo) (*domain.User, error) {
	if _, err := mail.ParseAddress(new.Email); err != nil {
		return nil, domain.ErrInvalidEmail
	}

	ok, err := s.repo.UserExistsByEmail(ctx, new.Email)
	if err != nil {
		return nil, fmt.Errorf("error checking if user exists: %w", err)
	}
	if ok {
		return nil, domain.ErrUserAlreadyExists
	}

	err = hashPassword(new)
	if err != nil {
		return nil, err
	}

	user, err := s.repo.CreateUser(ctx, new)
	if err != nil {
		return nil, fmt.Errorf("error creating user: %w", err)
	}

	return user, nil
}

func hashPassword(u *domain.UserInfo) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hash)
	return nil
}

func (s *userService) GetUserById(ctx context.Context, id int) (*domain.User, error) {
	ok, err := s.repo.UserExistsById(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("can't check if user exists: %w", err)
	}

	if !ok {
		return nil, domain.ErrNoSuchUser
	}

	account, err := s.repo.GetUser(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("can't get user by id: %w", err)
	}

	return account, nil
}

func (s *userService) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	if _, err := mail.ParseAddress(email); err != nil {
		return nil, err
	}

	ok, err := s.repo.UserExistsByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("can't check by email if user exists: %w", err)
	}
	if !ok {
		return nil, domain.ErrNoSuchUser
	}

	id, err := s.repo.GetUserIdByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("can' get user id by email: %w", err)
	}

	user, err := s.repo.GetUser(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("can't get user: %w", err)
	}

	return user, nil
}

func (s *userService) UpdateUserInfo(ctx context.Context, id int, newInfo *domain.UserInfo) (*domain.User, error) {
	if !(newInfo.Name != "" || newInfo.Email != "") {
		return nil, domain.ErrEmptyUserInfo
	}

	ok, err := s.repo.UserExistsById(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("can't check by id if user exists: %w", err)
	}
	if !ok {
		return nil, domain.ErrNoSuchUser
	}

	var info domain.UserInfo
	if newInfo.Name != "" {
		info.Name = newInfo.Name
	}
	if newInfo.Email != "" {
		if _, err = mail.ParseAddress(newInfo.Email); err != nil {
			return nil, domain.ErrInvalidEmail
		}
		info.Email = newInfo.Email
	}

	user, err := s.repo.UpdateUser(ctx, id, &info)
	if err != nil {
		return nil, fmt.Errorf("can't update user: %w", err)
	}

	return user, nil
}

func (s *userService) DeleteUserById(ctx context.Context, id int) error {
	ok, err := s.repo.UserExistsById(ctx, id)
	if err != nil {
		return fmt.Errorf("can't check by id if user exists: %w", err)
	}
	if !ok {
		return domain.ErrNoSuchUser
	}

	err = s.repo.DeleteUser(ctx, id)
	if err != nil {
		return fmt.Errorf("can't delete user: %w", err)
	}

	return nil
}

func (s *userService) AuthenticateUser(ctx context.Context, login *domain.UserInfo) (*domain.User, error) {
	user, err := s.GetUserByEmail(ctx, login.Email)
	if err != nil {
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(login.Password)); err != nil {
		return nil, domain.ErrWrongPassword
	}

	return user, nil
}
