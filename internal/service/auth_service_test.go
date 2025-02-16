package service_test

import (
	"context"
	"database/sql"
	"io"
	"log/slog"
	"testing"

	"github.com/google/uuid"
	"github.com/koccyx/avito_assignment/internal/entities"
	"github.com/koccyx/avito_assignment/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)


type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) GetByName(ctx context.Context, username string) (*entities.User, error) {
	args := m.Called(ctx, username)
	return args.Get(0).(*entities.User), args.Error(1)
}

func (m *MockUserRepository) Create(ctx context.Context, tx *sql.Tx, username, password string) (*uuid.UUID, error) {
	args := m.Called(ctx, tx, username, password)
	return args.Get(0).(*uuid.UUID), args.Error(1)
}

func (m *MockUserRepository) GetOne(ctx context.Context, userId uuid.UUID) (*entities.User, error) {
	args := m.Called(ctx, userId)
	return args.Get(0).(*entities.User), args.Error(1)
}

func (m *MockUserRepository) GetUserItemsInfo(ctx context.Context, userId uuid.UUID) ([]entities.InventoryItem, error) {
	args := m.Called(ctx, userId)
	return args.Get(0).([]entities.InventoryItem), args.Error(1)
}

func (m *MockUserRepository) PutCoins(ctx context.Context, tx *sql.Tx, userId uuid.UUID, amount int) (int, error) {
	args := m.Called(ctx, tx, userId, amount)
	return args.Get(0).(int), args.Error(1)
}

func (m *MockUserRepository) GetAllInfoByUserId(ctx context.Context, userId uuid.UUID) ([]entities.UserItem, error) {
	args := m.Called(ctx, userId)
	return args.Get(0).([]entities.UserItem), args.Error(1)
}

type MockJWTService struct {
	mock.Mock
}

func (m *MockJWTService) NewToken(userId string, secret string) (string, error) {
	args := m.Called(userId, secret)
	return args.String(0), args.Error(1)
}

func (m *MockJWTService) ParseToken(token string, secret string) (string, error) {
	args := m.Called(token, secret)
	return args.String(0), args.Error(1)
}

func TestAuthServiceAuth(t *testing.T) {
	userRepo := new(MockUserRepository)
	logger := slog.New(slog.NewJSONHandler(io.Discard, nil))
	secret := "testsecret"

	username := "testuser"
	password := "password123"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	user := &entities.User{
		Id:       uuid.New(),
		Username: username,
		Password: string(hashedPassword),
	}

	userRepo.On("GetByName", mock.Anything, username).Return(user, nil)

	authService := service.NewAuthService(userRepo, logger, secret, nil)

	token, err := authService.Auth(context.Background(), username, password)

	assert.NoError(t, err)
	assert.NotEmpty(t, token)
	userRepo.AssertExpectations(t)
}

func TestAuthServiceAuthInvalidCredentials(t *testing.T) {
	userRepo := new(MockUserRepository)
	logger := slog.New(slog.NewJSONHandler(io.Discard, nil))
	secret := "testsecret"

	username := "testuser"
	password := "password123"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	user := &entities.User{
		Id:       uuid.New(),
		Username: username,
		Password: string(hashedPassword),
	}

	userRepo.On("GetByName", mock.Anything, username).Return(user, nil)

	authService := service.NewAuthService(userRepo, logger, secret, nil)

	token, err := authService.Auth(context.Background(), username, "wrongpassword")

	assert.Error(t, err)
	assert.Empty(t, token)
	userRepo.AssertExpectations(t)
}

