package service

import (
	"context"
	"testing"
	"time"

	"bank-api/internal/domain"
	"bank-api/mocks"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestCreateUser(t *testing.T) {
	repoMock := mocks.NewMockUserRepository(gomock.NewController(t))

	userInfo := &domain.UserInfo{
		Name:  "Test User",
		Email: "test@example.com",
	}

	repoMock.EXPECT().UserExistsByEmail(gomock.Any(), userInfo.Email).Return(false, nil)

	repoMock.EXPECT().CreateUser(gomock.Any(), userInfo).Return(&domain.User{Id: 1, Name: "Test User", Email: "test@example.com", HashedPassword: "hash", CreatedAt: time.Now()}, nil)

	s := NewUserService(repoMock)

	user, err := s.CreateUser(context.Background(), userInfo)
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, 1, user.Id)
	assert.Equal(t, "Test User", user.Name)
	assert.Equal(t, "test@example.com", user.Email)
	assert.Equal(t, "hash", user.HashedPassword)
}

func TestCreateUser_UserAlreadyExists(t *testing.T) {
	mockRepo := mocks.NewMockUserRepository(gomock.NewController(t))

	userInfo := &domain.UserInfo{
		Name:  "Test User",
		Email: "test@example.com",
	}

	mockRepo.EXPECT().UserExistsByEmail(gomock.Any(), userInfo.Email).Return(true, nil)

	s := NewUserService(mockRepo)

	user, err := s.CreateUser(context.Background(), userInfo)
	assert.Nil(t, user)
	assert.ErrorIs(t, err, domain.ErrUserAlreadyExists)
}

func TestCreateUser_BadEmail(t *testing.T) {
	mockRepo := mocks.NewMockUserRepository(gomock.NewController(t))

	userInfo := &domain.UserInfo{
		Name:  "Test User",
		Email: "",
	}

	s := NewUserService(mockRepo)

	user, err := s.CreateUser(context.Background(), userInfo)
	assert.Nil(t, user)
	assert.Equal(t, domain.ErrInvalidEmail, err)
}

func TestUserService_GetUserById(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)

	user := &domain.User{
		Id:             1,
		Name:           "Test User",
		Email:          "test@example.com",
		HashedPassword: "hash",
		CreatedAt:      time.Now(),
	}

	mockRepo.EXPECT().UserExistsById(gomock.Any(), user.Id).Return(true, nil)
	mockRepo.EXPECT().GetUser(gomock.Any(), user.Id).Return(user, nil)

	s := NewUserService(mockRepo)

	user, err := s.GetUserById(context.Background(), user.Id)
	assert.NotNil(t, user)
	assert.Equal(t, 1, user.Id)
	assert.Equal(t, "Test User", user.Name)
	assert.Equal(t, "test@example.com", user.Email)

	mockRepo.EXPECT().UserExistsById(gomock.Any(), 2).Return(false, nil)

	user, err = s.GetUserById(context.Background(), 2)
	assert.ErrorIs(t, domain.ErrNoSuchUser, err)
	assert.Nil(t, user)
}

func TestUserService_GetUserByEmail(t *testing.T) {
	mockRepo := mocks.NewMockUserRepository(gomock.NewController(t))

	user := &domain.User{
		Id:             1,
		Name:           "Test User",
		Email:          "test@example.com",
		HashedPassword: "hash",
		CreatedAt:      time.Now(),
	}

	mockRepo.EXPECT().UserExistsByEmail(gomock.Any(), user.Email).Return(true, nil)
	mockRepo.EXPECT().GetUserIdByEmail(gomock.Any(), user.Email).Return(user.Id, nil)
	mockRepo.EXPECT().GetUser(gomock.Any(), user.Id).Return(user, nil)

	s := NewUserService(mockRepo)

	user, err := s.GetUserByEmail(context.Background(), user.Email)
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, 1, user.Id)
	assert.Equal(t, "Test User", user.Name)
	assert.Equal(t, "test@example.com", user.Email)

	mockRepo.EXPECT().UserExistsByEmail(gomock.Any(), user.Email).Return(false, nil)

	user, err = s.GetUserByEmail(context.Background(), user.Email)
	assert.Nil(t, user)
	assert.ErrorIs(t, domain.ErrNoSuchUser, err)
}

func TestUserService_UpdateUserInfo(t *testing.T) {
	mockRepo := mocks.NewMockUserRepository(gomock.NewController(t))

	user := &domain.User{
		Id:             1,
		Name:           "Kylie Jenner",
		Email:          "test@example.com",
		HashedPassword: "hash",
		CreatedAt:      time.Now(),
	}

	userInfo := &domain.UserInfo{
		Name:  "Kylie Chalamet",
		Email: "test@example.com",
	}

	mockRepo.EXPECT().UserExistsById(gomock.Any(), user.Id).Return(true, nil)
	mockRepo.EXPECT().UpdateUser(gomock.Any(), user.Id, userInfo).Return(user, nil)

	s := NewUserService(mockRepo)

	user, err := s.UpdateUserInfo(context.Background(), user.Id, userInfo)
	assert.NoError(t, err)
	assert.NotNil(t, user)
}

func TestUserService_DeleteUserById(t *testing.T) {
	mockRepo := mocks.NewMockUserRepository(gomock.NewController(t))

	user := &domain.User{
		Id:             1,
		Name:           "Test User",
		Email:          "test@example.com",
		HashedPassword: "hash",
		CreatedAt:      time.Now(),
	}

	mockRepo.EXPECT().UserExistsById(gomock.Any(), user.Id).Return(true, nil)
	mockRepo.EXPECT().DeleteUser(gomock.Any(), user.Id).Return(nil)

	s := NewUserService(mockRepo)

	err := s.DeleteUserById(context.Background(), user.Id)
	assert.NoError(t, err)

	mockRepo.EXPECT().UserExistsById(gomock.Any(), user.Id).Return(false, nil)

	err = s.DeleteUserById(context.Background(), user.Id)
	assert.ErrorIs(t, domain.ErrNoSuchUser, err)
}
