package entities

import (
	"github.com/google/uuid"
)

type Item struct {
	Id    uuid.UUID `db:"id"`
	Name  string    `db:"name"`
	Price int       `db:"price"`
}

type InventoryItem struct {
	Name   string
	Amount int
}
