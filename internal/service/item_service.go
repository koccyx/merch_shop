package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	"github.com/koccyx/avito_assignment/internal/lib/sl"
	"github.com/koccyx/avito_assignment/internal/storage/postgres"
)

type ItemServiceSt struct {
	log          *slog.Logger
	userRepo     UserRepository
	itemRepo     ItemRepository
	userItemRepo UserItemRepository
	db           *sql.DB
}

func (s *ItemServiceSt) PurchaseItem(ctx context.Context, userId, itemName string) error {
	const op = "service.item.PurchaseItem"

	log := s.log.With(
		slog.String("op", op),
		slog.String("userId", userId),
	)

	log.Info("purchasing item")

	prsdUsrId, err := uuid.Parse(userId)
	if err != nil {
		log.Error("error while parsing userId", sl.Err(err))
		return err
	}

	usr, err := s.userRepo.GetOne(ctx, prsdUsrId)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			log.Error("error while getting user", sl.Err(err))
			return fmt.Errorf("%s: %w", op, err)
		}

		log.Error("user not found")
		return fmt.Errorf("%s: %w", op, ErrNoEntry)
	}

	item, err := s.itemRepo.GetByName(ctx, itemName)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			log.Error("error while getting item", sl.Err(err))
			return fmt.Errorf("%s: %w", op, err)
		}

		log.Info("item not found")
		return fmt.Errorf("%s: %w", op, ErrNoEntry)
	}

	if usr.Balance < item.Price {
		return fmt.Errorf("%s: %w", op, ErrNotEnoughBalance)
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

	_, err = s.userRepo.PutCoins(ctx, tx, usr.Id, -item.Price)
	if err != nil {
		if errors.Is(err, postgres.ErrNotEnoughBalance) {
			return fmt.Errorf("%s: %w", op, ErrNotEnoughBalance)
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = s.userItemRepo.Create(ctx, tx, usr.Id, item.Id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("%s successfuly purchased item %s", usr.Username, item.Name)

	return nil
}

func NewItemService(tdRepo UserRepository, itemRepo ItemRepository, userItemRepo UserItemRepository, logger *slog.Logger, db *sql.DB) *ItemServiceSt {
	return &ItemServiceSt{
		userRepo:     tdRepo,
		itemRepo:     itemRepo,
		userItemRepo: userItemRepo,
		log:          logger,
		db:           db,
	}
}
