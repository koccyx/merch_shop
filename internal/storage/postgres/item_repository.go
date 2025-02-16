package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/koccyx/avito_assignment/internal/entities"
)

type ItemRepository struct {
	db *sql.DB
}

func (r *ItemRepository) GetOne(ctx context.Context, itemId uuid.UUID) (*entities.Item, error) {
	const op = "repo.postgres.item.GetOne"

	if itemId == uuid.Nil {
		return nil, fmt.Errorf("%s: missing id", op)
	}

	builder := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
	selectBuilder := builder.
		Select("id", "name", "price").
		From("items").
		Where(squirrel.Eq{"id": itemId})

	query, args, err := selectBuilder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var item entities.Item

	err = r.db.QueryRowContext(ctx, query, args...).Scan(&item.Id, &item.Name, &item.Price)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &item, nil
}

func (r *ItemRepository) GetAll(ctx context.Context, usrId uuid.UUID) ([]entities.Item, error) {
	const op = "repo.postgres.item.GetAll"

	if usrId == uuid.Nil {
		return nil, fmt.Errorf("%s: missing id", op)
	}

	builder := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
	selectBuilder := builder.
		Select("id", "name", "price").
		From("items")

	query, args, err := selectBuilder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	items := make([]entities.Item, 0)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	for rows.Next() {
		var item entities.Item
		err := rows.Scan(&item.Id, &item.Name)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		items = append(items, item)
	}
	return items, nil
}

func (r *ItemRepository) GetByName(ctx context.Context, name string) (*entities.Item, error) {
	const op = "repo.postgres.item.GetByName"

	builder := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
	selectBuilder := builder.
		Select("id", "name", "price").
		From("items").
		Where(squirrel.Eq{"name": name})

	query, args, err := selectBuilder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var item entities.Item

	err = r.db.QueryRowContext(ctx, query, args...).Scan(&item.Id, &item.Name, &item.Price)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &item, nil
}

func NewItemRepository(db *sql.DB) *ItemRepository {
	return &ItemRepository{
		db: db,
	}
}
