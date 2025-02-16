package entities

import (
	"time"

	"github.com/google/uuid"
)

type Direction string

const (
	FromDirection Direction = "FROM"
	ToDirection   Direction = "TO"
)

type Transaction struct {
	Id         uuid.UUID `db:"id"`
	FromUserId uuid.UUID `db:"user_id"`
	ToUserId   uuid.UUID `db:"item_id"`
	Amount     int       `db:"amount"`
	Created_at time.Time `db:"created_at"`
}

type CoinTransactionInfo struct {
	FromUser string
	ToUser   string
	Amount   int
}
