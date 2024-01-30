package bank

import (
	"bank-api/internal/model"
	"bank-api/internal/repository"
	"context"
	"errors"
	"net/mail"
)

var (
	ErrEmptyUserInfo     = errors.New("empty account info")
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrSameAccountInfo   = errors.New("same account info")
	ErrInvalidEmail      = errors.New("invalid email")
)

type UserService interface {
	CreateUser(ctx context.Context, info *model.UserInfo) (*model.User, error)
	GetUserById(ctx context.Context, id int) (*model.User, error)
	GetUserByEmail(ctx context.Context, email string) (*model.User, error)
	UpdateUserInfo(ctx context.Context, id int, info *model.UserInfo) (*model.User, error)
	DeleteUserById(ctx context.Context, id int) error
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
}
func (s *userService) CreateUser(ctx context.Context, new *model.UserInfo) (*model.User, error) {

	_, err := s.repo.GetUserIdByEmail(ctx, new.Email)
	if err == nil {
		return nil, ErrUserAlreadyExists
	}

	if _, err := mail.ParseAddress(new.Email); err != nil {
		return nil, ErrInvalidEmail
	}
	user, err := s.repo.CreateUser(ctx, new)
	if err != nil {

	}

	return user, nil
}

func (s *userService) GetUserById(ctx context.Context, id int) (*model.User, error) {

	account, err := s.repo.GetUser(ctx, id)
	if err != nil {
		return nil, err
	}

	return account, nil
}

func (s *userService) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	if _, err := mail.ParseAddress(email); err != nil {
		return nil, err
	}
	id, err := s.repo.GetUserIdByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	user, err := s.repo.GetUser(ctx, id)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userService) UpdateUserInfo(ctx context.Context, id int, newInfo *model.UserInfo) (*model.User, error) {
	if !(newInfo.Name != "" || newInfo.Email != "") {
		return nil, ErrEmptyUserInfo
	}

	oldInfo, err := s.repo.GetUser(ctx, id)
	if err != nil {
		return nil, err
	}

	if newInfo.Name == oldInfo.Name && newInfo.Email == oldInfo.Email {
		return nil, ErrSameAccountInfo
	}

	var info model.UserInfo
	if newInfo.Name != "" {
		info.Name = newInfo.Name
	}
	if newInfo.Email != "" {
		info.Email = newInfo.Email
	}

	user, err := s.repo.UpdateUser(ctx, id, &info)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userService) DeleteUserById(ctx context.Context, id int) error {
	_, err := s.repo.GetUser(ctx, id)
	if err != nil {
		return err
	}

	err = s.repo.DeleteUser(ctx, id)
	if err != nil {
		return err
	}

	return nil
}
