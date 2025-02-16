package service_test

import (
	"context"
	"io"
	"log/slog"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/koccyx/avito_assignment/internal/entities"
	"github.com/koccyx/avito_assignment/internal/server/models"
	"github.com/koccyx/avito_assignment/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)
func TestTransferCoins(t *testing.T) {
	db, sqlMock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error opening mock database: %s", err)
	}
	defer db.Close()

	sqlMock.ExpectBegin()
	sqlMock.ExpectCommit()

	mockUserRepo := new(MockUserRepository)
	mockTransactionRepo := new(MockTransactionRepository)

	logger := slog.New(slog.NewJSONHandler(io.Discard, nil))
	service := service.NewUserService(mockUserRepo, nil, nil, mockTransactionRepo, logger, db)

	ctx := context.Background()
	userFromId := uuid.New()
	usernameTo := "receiver"
	amount := 50

	userFrom := &entities.User{
		Id:       userFromId,
		Username: "sender",
		Balance:  100,
	}
	userTo := &entities.User{
		Id:       uuid.New(),
		Username: usernameTo,
		Balance:  50,
	}

	mockUserRepo.On("GetOne", ctx, userFromId).Return(userFrom, nil)
	mockUserRepo.On("GetByName", ctx, usernameTo).Return(userTo, nil)
	mockUserRepo.On("PutCoins", ctx, mock.Anything, userFromId, -amount).Return(1, nil)
	mockUserRepo.On("PutCoins", ctx, mock.Anything, userTo.Id, amount).Return(1, nil)
	mockTransactionRepo.On("Create", ctx, mock.Anything, userFromId, userTo.Id, amount).Return(&uuid.UUID{}, nil)

	err = service.TransferCoins(ctx, userFromId.String(), usernameTo, amount)

	assert.NoError(t, err)
	mockUserRepo.AssertExpectations(t)
	mockTransactionRepo.AssertExpectations(t)
}

func TestInfo(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error opening mock database: %s", err)
	}
	defer db.Close()

	mockUserRepo := new(MockUserRepository)
	mockTransactionRepo := new(MockTransactionRepository)

	logger := slog.New(slog.NewJSONHandler(io.Discard, nil))
	service := service.NewUserService(mockUserRepo, nil, nil, mockTransactionRepo, logger, db)

	ctx := context.Background()
	userId := uuid.New()

	user := &entities.User{
		Id:       userId,
		Username: "test_user",
		Balance:  200,
	}

	transactions := []entities.CoinTransactionInfo{}
	items := []entities.InventoryItem{}

	mockUserRepo.On("GetOne", ctx, userId).Return(user, nil)
	mockUserRepo.On("GetUserItemsInfo", ctx, userId).Return(items, nil)
	mockTransactionRepo.On("GetAllWithDirection", ctx, userId, entities.FromDirection).Return(transactions, nil)
	mockTransactionRepo.On("GetAllWithDirection", ctx, userId, entities.ToDirection).Return(transactions, nil)

	info, err := service.Info(ctx, userId.String())

	assert.NoError(t, err)
	assert.Equal(t, 200, info.Coins)
	assert.Empty(t, info.Inventory)
	assert.Equal(t, models.CoinHistory{
		Received: []models.CoinTransactionRecived{},
		Sent: []models.CoinTransactionSent{},
	},info.CoinHistory)

	mockUserRepo.AssertExpectations(t)
	mockTransactionRepo.AssertExpectations(t)
}