package storage

import (
	"context"
	"database/sql"

	"github.com/google/uuid"

	"github.com/koccyx/avito_assignment/internal/entities"
	"github.com/koccyx/avito_assignment/internal/storage/postgres"
)

type UserRepository interface {
	GetOne(ctx context.Context, usrId uuid.UUID) (*entities.User, error)
	Create(ctx context.Context, tx *sql.Tx, username, password string) (*uuid.UUID, error)
	GetByName(ctx context.Context, username string) (*entities.User, error)
	PutCoins(ctx context.Context, tx *sql.Tx, userId uuid.UUID, amount int) (int, error)
	GetUserItemsInfo(ctx context.Context, userId uuid.UUID) ([]entities.InventoryItem, error)
}

type UserItemRepository interface {
	Create(ctx context.Context, tx *sql.Tx, userId uuid.UUID, itemId uuid.UUID) (*uuid.UUID, error)
	GetAllInfoByUserId(ctx context.Context, usrId uuid.UUID) ([]entities.UserItem, error)
	GetOne(ctx context.Context, userItemId uuid.UUID) (*entities.UserItem, error)
}

type TransactionRepository interface {
	Create(ctx context.Context, tx *sql.Tx, fromUserId uuid.UUID, toUserId uuid.UUID, amount int) (*uuid.UUID, error)
	GetAllWithDirection(ctx context.Context, usrId uuid.UUID, direction entities.Direction) ([]entities.CoinTransactionInfo, error)
	GetOne(ctx context.Context, transactionId uuid.UUID) (*entities.Transaction, error)
}

type ItemRepository interface {
	GetOne(ctx context.Context, itemId uuid.UUID) (*entities.Item, error)
	GetAll(ctx context.Context, usrId uuid.UUID) ([]entities.Item, error)
	GetByName(ctx context.Context, name string) (*entities.Item, error)
}

type Repository struct {
	User        UserRepository
	Item        ItemRepository
	UserItem    UserItemRepository
	Transaction TransactionRepository
}

func NewRepository(db *sql.DB) *Repository {
	user := postgres.NewUserRepository(db)
	item := postgres.NewItemRepository(db)
	userItem := postgres.NewUserItemRepository(db)
	transaction := postgres.NewTransactionRepository(db)

	return &Repository{
		User:        user,
		Item:        item,
		UserItem:    userItem,
		Transaction: transaction,
	}
}
