package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	"github.com/koccyx/avito_assignment/internal/entities"
	"github.com/koccyx/avito_assignment/internal/lib/sl"
	"github.com/koccyx/avito_assignment/internal/mappers"
	"github.com/koccyx/avito_assignment/internal/server/models"
	"github.com/koccyx/avito_assignment/internal/storage/postgres"
)

type UserServiceSt struct {
	log             *slog.Logger
	db              *sql.DB
	userRepo        UserRepository
	itemRepo        ItemRepository
	transactionRepo TransactionRepository
	userItemRepo    UserItemRepository
}

func (s *UserServiceSt) TransferCoins(ctx context.Context, userFromId, usernameTo string, amount int) error {
	const op = "service.user.TransferCoins"

	log := s.log.With(
		slog.String("op", op),
		slog.String("From", userFromId),
	)

	log.Info("transfering coins")

	prsdUsrFromId, err := uuid.Parse(userFromId)
	if err != nil {
		log.Error("error while parsing userId", sl.Err(err))
		return err
	}

	fromUsr, err := s.userRepo.GetOne(ctx, prsdUsrFromId)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			log.Error("error while getting from user", sl.Err(err))
			return fmt.Errorf("%s: %w", op, err)
		}

		log.Info("from user not found")
		return fmt.Errorf("%s: %w", op, ErrNoEntry)
	}

	if fromUsr.Balance < amount {
		log.Error("not enough balance", sl.Err(ErrNotEnoughBalance))
		return fmt.Errorf("%s: %w", op, ErrNotEnoughBalance)
	}

	toUsr, err := s.userRepo.GetByName(ctx, usernameTo)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			log.Error("error while getting to user", sl.Err(err))
			return fmt.Errorf("%s: %w", op, err)
		}

		log.Error("to user not found")
		return fmt.Errorf("%s: %w", op, ErrNoEntry)
	}

	if fromUsr.Id == toUsr.Id {
		log.Error("cant transfer coins to the same user", sl.Err(ErrSameUserTransfer))
		return fmt.Errorf("%s: %w", op, ErrSameUserTransfer)
	}

	tx, err := s.db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
	})
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
		tx.Commit()
	}()

	_, err = s.userRepo.PutCoins(ctx, tx, fromUsr.Id, -amount)
	if err != nil {
		if errors.Is(err, postgres.ErrNotEnoughBalance) {
			return fmt.Errorf("%s: %w", op, ErrNotEnoughBalance)
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = s.userRepo.PutCoins(ctx, tx, toUsr.Id, amount)
	if err != nil {
		if errors.Is(err, postgres.ErrNotEnoughBalance) {
			return fmt.Errorf("%s: %w", op, ErrNotEnoughBalance)
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = s.transactionRepo.Create(ctx, tx, fromUsr.Id, toUsr.Id, amount)
	if err != nil {
		log.Error("error during transaction", sl.Err(err))
		return fmt.Errorf("%s: %w", op, fmt.Errorf("error during transaction creation"))
	}

	log.Info("%s successfuly transfered coins from %s", fromUsr.Username, toUsr.Username)

	return nil
}

func (s *UserServiceSt) Info(ctx context.Context, userId string) (*models.InfoResponse, error) {
	const op = "service.user.Info"

	log := s.log.With(
		slog.String("op", op),
		slog.String("From", userId),
	)

	log.Info("transfering coins")

	prsdUsrId, err := uuid.Parse(userId)
	if err != nil {
		log.Error("error while parsing userId", sl.Err(err))
		return nil, err
	}

	usr, err := s.userRepo.GetOne(ctx, prsdUsrId)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			log.Error("error while getting from user", sl.Err(err))
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		log.Info("from user not found")
		return nil, fmt.Errorf("%s: %w", op, ErrNoEntry)
	}

	usrItems, err := s.userRepo.GetUserItemsInfo(ctx, usr.Id)
	if err != nil {
		log.Error("error while getting userItems Info")
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	inventory := mappers.MapInventoryItemsEntityToModel(usrItems)

	sentTr, err := s.transactionRepo.GetAllWithDirection(ctx, usr.Id, entities.FromDirection)
	if err != nil {
		log.Error("error while getting userItems")
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	recivedTr, err := s.transactionRepo.GetAllWithDirection(ctx, usr.Id, entities.ToDirection)
	if err != nil {
		log.Error("error while getting userItems")
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	coinsHistory := mappers.MapTransactionEntityToCoinsHistory(recivedTr, sentTr)

	log.Info(fmt.Sprintf("successfuly got info about %s", usr.Username))

	return &models.InfoResponse{
		CoinHistory: coinsHistory,
		Coins:       usr.Balance,
		Inventory:   inventory,
	}, nil
}

func NewUserService(tdRepo UserRepository, itemRepo ItemRepository, userItemRepo UserItemRepository, transactionRepo TransactionRepository, logger *slog.Logger, db *sql.DB) *UserServiceSt {
	return &UserServiceSt{
		userRepo:        tdRepo,
		itemRepo:        itemRepo,
		userItemRepo:    userItemRepo,
		transactionRepo: transactionRepo,
		log:             logger,
		db:              db,
	}
}
