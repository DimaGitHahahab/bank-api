package bank

import (
	"bank-api/internal/model"
	"bank-api/mocks"
	"context"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestCreateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)

	userInfo := &model.UserInfo{
		Name:  "Test User",
		Email: "test@example.com",
	}

	mockRepo.EXPECT().GetUserIdByEmail(gomock.Any(), userInfo.Email).Return(0, ErrNoSuchUser)
	mockRepo.EXPECT().CreateUser(gomock.Any(), userInfo).Return(&model.User{Id: 1, Name: "Test User", Email: "test@example.com", HashedPassword: "hash", CreatedAt: time.Now()}, nil)

	s := NewUserService(mockRepo)

	user, err := s.CreateUser(context.Background(), userInfo)
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, 1, user.Id)
	assert.Equal(t, "Test User", user.Name)
	assert.Equal(t, "test@example.com", user.Email)
	assert.Equal(t, "hash", user.HashedPassword)
}

func TestCreateUser_UserAlreadyExists(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)

	userInfo := &model.UserInfo{
		Name:  "Test User",
		Email: "test@example.com",
	}

	mockRepo.EXPECT().GetUserIdByEmail(gomock.Any(), userInfo.Email).Return(0, ErrNoSuchUser)
	mockRepo.EXPECT().CreateUser(gomock.Any(), userInfo).Return(&model.User{Id: 1, Name: "Test User", Email: "test@example.com", HashedPassword: "hash", CreatedAt: time.Now()}, nil)
	mockRepo.EXPECT().GetUserIdByEmail(gomock.Any(), userInfo.Email).Return(1, nil)

	s := NewUserService(mockRepo)

	user, err := s.CreateUser(context.Background(), userInfo)
	assert.NoError(t, err)
	assert.NotNil(t, user)

	user, err = s.CreateUser(context.Background(), userInfo)
	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Equal(t, ErrUserAlreadyExists, err)
}

func TestCreateUser_BadEmail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)

	userInfo := &model.UserInfo{
		Name:  "Test User",
		Email: "",
	}

	s := NewUserService(mockRepo)

	user, err := s.CreateUser(context.Background(), userInfo)
	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Equal(t, ErrInvalidEmail, err)

	userInfo.Email = "test@"

	user, err = s.CreateUser(context.Background(), userInfo)
	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Equal(t, ErrInvalidEmail, err)
}

func TestUserService_GetUserById(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)

	user := &model.User{
		Id:             1,
		Name:           "Test User",
		Email:          "test@example.com",
		HashedPassword: "hash",
		CreatedAt:      time.Now(),
	}

	mockRepo.EXPECT().GetUser(gomock.Any(), 1).Return(user, nil)

	s := NewUserService(mockRepo)

	user, err := s.GetUserById(context.Background(), user.Id)
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, 1, user.Id)
	assert.Equal(t, "Test User", user.Name)
	assert.Equal(t, "test@example.com", user.Email)

	mockRepo.EXPECT().GetUser(gomock.Any(), 2).Return(nil, ErrNoSuchUser)

	user, err = s.GetUserById(context.Background(), 2)
	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Equal(t, ErrNoSuchUser, err)
}

func TestUserService_GetUserByEmail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)

	user := &model.User{
		Id:             1,
		Name:           "Test User",
		Email:          "test@example.com",
		HashedPassword: "hash",
		CreatedAt:      time.Now(),
	}

	mockRepo.EXPECT().GetUserIdByEmail(gomock.Any(), user.Email).Return(user.Id, nil)
	mockRepo.EXPECT().GetUser(gomock.Any(), user.Id).Return(user, nil)

	s := NewUserService(mockRepo)

	user, err := s.GetUserByEmail(context.Background(), user.Email)
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, 1, user.Id)
	assert.Equal(t, "Test User", user.Name)
	assert.Equal(t, "test@example.com", user.Email)

	mockRepo.EXPECT().GetUserIdByEmail(gomock.Any(), user.Email).Return(0, ErrNoSuchUser)

	user, err = s.GetUserByEmail(context.Background(), user.Email)
	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Equal(t, ErrNoSuchUser, err)
}

func TestUserService_UpdateUserInfo(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)

	user := &model.User{
		Id:             1,
		Name:           "Kylie Jenner",
		Email:          "test@example.com",
		HashedPassword: "hash",
		CreatedAt:      time.Now(),
	}

	userInfo := &model.UserInfo{
		Name:  "Kylie Chalamet",
		Email: "test@example.com",
	}

	mockRepo.EXPECT().GetUser(gomock.Any(), user.Id).Return(user, nil)
	mockRepo.EXPECT().UpdateUser(gomock.Any(), user.Id, userInfo).Return(user, nil)

	s := NewUserService(mockRepo)

	user, err := s.UpdateUserInfo(context.Background(), user.Id, userInfo)
	assert.NoError(t, err)
	assert.NotNil(t, user)

}

func TestUserService_DeleteUserById(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)

	user := &model.User{
		Id:             1,
		Name:           "Test User",
		Email:          "test@example.com",
		HashedPassword: "hash",
		CreatedAt:      time.Now(),
	}

	mockRepo.EXPECT().GetUser(gomock.Any(), user.Id).Return(user, nil)
	mockRepo.EXPECT().DeleteUser(gomock.Any(), user.Id).Return(nil)

	s := NewUserService(mockRepo)

	err := s.DeleteUserById(context.Background(), user.Id)
	assert.NoError(t, err)

	mockRepo.EXPECT().GetUser(gomock.Any(), user.Id).Return(nil, ErrNoSuchUser)

	err = s.DeleteUserById(context.Background(), user.Id)
	assert.Error(t, err)
	assert.Equal(t, ErrNoSuchUser, err)

}
