package service_test

import (
	"context"
	"database/sql"
	"io"
	"log/slog"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/koccyx/avito_assignment/internal/entities"
	"github.com/koccyx/avito_assignment/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type RepositoryMock struct {
	mock.Mock
}

type MockItemRepository struct {
	mock.Mock
}

func (m *MockItemRepository) GetOne(ctx context.Context, itemId uuid.UUID) (*entities.Item, error) {
	args := m.Called(ctx, itemId)
	return args.Get(0).(*entities.Item), args.Error(1)
}

func (m *MockItemRepository) GetAll(ctx context.Context, usrId uuid.UUID) ([]entities.Item, error) {
	args := m.Called(ctx, usrId)
	return args.Get(0).([]entities.Item), args.Error(1)
}

func (m *MockItemRepository) GetByName(ctx context.Context, name string) (*entities.Item, error) {
	args := m.Called(ctx, name)
	return args.Get(0).(*entities.Item), args.Error(1)
}


type MockTransactionRepository struct {
	mock.Mock
}

func (m *MockTransactionRepository) Create(ctx context.Context, tx *sql.Tx, fromUserId uuid.UUID, toUserId uuid.UUID, amount int) (*uuid.UUID, error) {
	args := m.Called(ctx, tx, fromUserId, toUserId, amount)
	return args.Get(0).(*uuid.UUID), args.Error(1)
}

func (m *MockTransactionRepository) GetAllWithDirection(ctx context.Context, usrId uuid.UUID, direction entities.Direction) ([]entities.CoinTransactionInfo, error) {
	args := m.Called(ctx, usrId, direction)
	return args.Get(0).([]entities.CoinTransactionInfo), args.Error(1)
}

func (m *MockTransactionRepository) GetOne(ctx context.Context, transactionId uuid.UUID) (*entities.Transaction, error) {
	args := m.Called(ctx, transactionId)
	return args.Get(0).(*entities.Transaction), args.Error(1)
}


type MockUserItemsRepository struct {
	mock.Mock
}

func (m *MockUserItemsRepository) Create(ctx context.Context, tx *sql.Tx, userId uuid.UUID, itemId uuid.UUID) (*uuid.UUID, error) {
	args := m.Called(ctx, tx, userId, itemId)
	return args.Get(0).(*uuid.UUID), args.Error(1)
}

func (m *MockUserItemsRepository) GetAllInfoByUserId(ctx context.Context, usrId uuid.UUID) ([]entities.UserItem, error) {
	args := m.Called(ctx, usrId)
	return args.Get(0).([]entities.UserItem), args.Error(1)
}

func (m *MockUserItemsRepository) GetOne(ctx context.Context, userItemId uuid.UUID) (*entities.UserItem, error) {
	args := m.Called(ctx, userItemId)
	return args.Get(0).(*entities.UserItem), args.Error(1)
}

func TestPurchaseItem(t *testing.T) {
	db, sqlMock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error opening mock database: %s", err)
	}
	defer db.Close()

	sqlMock.ExpectBegin()
	sqlMock.ExpectCommit()

	mockUserRepo := new(MockUserRepository)
	mockItemRepo := new(MockItemRepository)
	mockUserItemRepo := new(MockUserItemsRepository)

	logger := slog.New(slog.NewJSONHandler(io.Discard, nil))
	service := service.NewItemService(mockUserRepo, mockItemRepo, mockUserItemRepo, logger, db)

	ctx := context.Background()
	userId := uuid.New()
	itemName := "t-shirt"
	testUser := &entities.User{
		Id:       userId,
		Username: "tst",
		Balance:  100,
	}
	testItem := &entities.Item{
		Id:    uuid.New(),
		Name:  itemName,
		Price: 80,
	}

	newUUID := uuid.New()

	mockUserRepo.On("GetOne", ctx, userId).Return(testUser, nil)
	mockItemRepo.On("GetByName", ctx, itemName).Return(testItem, nil)
	mockUserRepo.On("PutCoins", ctx, mock.Anything, userId, -testItem.Price).Return(1, nil)
	mockUserItemRepo.On("Create", ctx, mock.Anything, userId, testItem.Id).Return(&newUUID, nil)

	err = service.PurchaseItem(ctx, userId.String(), itemName)

	assert.NoError(t, err)

	mockUserRepo.AssertExpectations(t)
	mockItemRepo.AssertExpectations(t)
	mockUserItemRepo.AssertExpectations(t)
}