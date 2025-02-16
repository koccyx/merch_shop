package postgres

import (
	"context"
	"database/sql"
	"fmt"
	
	"github.com/koccyx/avito_assignment/internal/entities"
	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
)

type TransactionRepository struct {
	db *sql.DB
}

func (r *TransactionRepository) Create(ctx context.Context, tx *sql.Tx, fromUserId uuid.UUID, toUserId uuid.UUID, amount int) (*uuid.UUID, error) {	
	const op = "repo.postgres.transaction.Create" 

	if fromUserId == uuid.Nil {
		return nil, fmt.Errorf("%s: missing from user id", op)
	}
	if toUserId == uuid.Nil {
		return nil, fmt.Errorf("%s: missing to user id", op)
	}

	id := uuid.New()

	builder := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
    createBuilder := builder.
	Insert("transactions").
	Columns("id","from_user_id", "to_user_id", "amount").
	Values(id, fromUserId, toUserId, amount)

    sql, args, err := createBuilder.ToSql()
    if err != nil {
        return nil, fmt.Errorf("%s: %w", op, err)
    }
	
	_, err = tx.ExecContext(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	
    return &id, nil
}


func (r *TransactionRepository) GetAllWithDirection(ctx context.Context, usrId uuid.UUID, direction entities.Direction) ([]entities.CoinTransactionInfo, error) {
	const op = "repo.postgres.transaction.GetAll" 

	if usrId == uuid.Nil {
		return nil, fmt.Errorf("%s: missing id", op)
	}

	var eq string
	
	if direction == entities.FromDirection {
		eq = "t.from_user_id"
	} else if direction == entities.ToDirection {
		eq = "t.to_user_id"
	}
	
    builder := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
    selectBuilder := builder.
	Select(
		"u1.username AS user_from_name",
		"u2.username AS user_to_name",
		"t.amount",
	).
	From("transactions t").
	Join("users u1 ON t.from_user_id = u1.id").
	Join("users u2 ON t.to_user_id = u2.id").
	Where(squirrel.Eq{eq: usrId}).
	OrderBy("t.created_at DESC")
	
    query, args, err := selectBuilder.ToSql()
    if err != nil {
        return nil, fmt.Errorf("%s: %w", op, err)
    }

	transactions := make([]entities.CoinTransactionInfo, 0)
	
    rows, err := r.db.QueryContext(ctx, query, args...)
    if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
    }
	
	for rows.Next() {
		var transaction entities.CoinTransactionInfo

		err := rows.Scan(&transaction.FromUser, &transaction.ToUser,  &transaction.Amount)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		transactions = append(transactions, transaction)
	}
    return transactions, nil
}

func (r *TransactionRepository) GetOne(ctx context.Context, transactionId uuid.UUID) (*entities.Transaction, error) {
	const op = "repo.postgres.transaction.GetOne" 
	
    builder := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
    selectBuilder := builder.
	Select("id", "from_user_id", "to_user_id", "amount", "created_at").
	From("transactions").
	Where(squirrel.Eq{"id": transactionId})

    query, args, err := selectBuilder.ToSql()
    if err != nil {
        return nil, fmt.Errorf("%s: %w", op, err)
    }

	var transaction entities.Transaction

    err = r.db.QueryRowContext(ctx, query, args...).Scan(&transaction.Id, &transaction.FromUserId,  &transaction.ToUserId, &transaction.Amount, &transaction.Created_at)
    if err != nil {
        return nil, fmt.Errorf("%s: %w", op, err)
    }

    return &transaction, nil
}

func NewTransactionRepository(db *sql.DB) *TransactionRepository {
	return &TransactionRepository{
		db: db,
	}
}
