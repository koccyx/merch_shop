package models

// InfoResponse представляет ответ на запрос информации о монетах, инвентаре и истории транзакций.
type InfoResponse struct {
    Coins       int           `json:"coins"`
    Inventory   []InventoryItem `json:"inventory"`
    CoinHistory CoinHistory    `json:"coinHistory"`
}

// InventoryItem представляет элемент инвентаря.
type InventoryItem struct {
    Type     string `json:"type"`
    Quantity int    `json:"quantity"`
}

// CoinHistory представляет историю получения и отправки монет.
type CoinHistory struct {
    Received []CoinTransactionRecived `json:"received"`
    Sent     []CoinTransactionSent `json:"sent"`
}

// CoinTransaction представляет одну транзакцию монет.
type CoinTransactionRecived struct {
    FromUser string `json:"fromUser"`
    Amount   int    `json:"amount"`
}

type CoinTransactionSent struct {
    ToUser   string `json:"toUser"`
    Amount   int    `json:"amount"`
}

// ErrorResponse представляет ответ с ошибкой.
type ErrorResponse struct {
    Errors string `json:"errors"`
}

// AuthRequest представляет запрос на аутентификацию.
type AuthRequest struct {
    Username string `json:"username" validate:"required,min=4"`
    Password string `json:"password" validate:"required,min=6"`
}

// AuthResponse представляет ответ на успешную аутентификацию.
type AuthResponse struct {
    Token string `json:"token"`
}

// SendCoinRequest представляет запрос на отправку монет.
type SendCoinRequest struct {
    ToUser string `json:"toUser" validate:"required"`
    Amount int    `json:"amount" validate:"required"`
}