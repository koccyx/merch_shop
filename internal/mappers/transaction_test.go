package mappers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/koccyx/avito_assignment/internal/entities"
	"github.com/koccyx/avito_assignment/internal/server/models"
)

func TestMapTransactionEntityToCoinsHistory(t *testing.T) {
	tests := []struct {
		name           string
		received       []entities.CoinTransactionInfo
		sent           []entities.CoinTransactionInfo
		expectedOutput models.CoinHistory
	}{
		{
			name:     "empty transactions",
			received: []entities.CoinTransactionInfo{},
			sent:     []entities.CoinTransactionInfo{},
			expectedOutput: models.CoinHistory{
				Received: []models.CoinTransactionRecived{},
				Sent:     []models.CoinTransactionSent{},
			},
		},
		{
			name: "single received transaction",
			received: []entities.CoinTransactionInfo{
				{FromUser: "user1", Amount: 100},
			},
			sent: []entities.CoinTransactionInfo{},
			expectedOutput: models.CoinHistory{
				Received: []models.CoinTransactionRecived{
					{FromUser: "user1", Amount: 100},
				},
				Sent: []models.CoinTransactionSent{},
			},
		},
		{
			name: "single sent transaction",
			received: []entities.CoinTransactionInfo{},
			sent: []entities.CoinTransactionInfo{
				{ToUser: "user2", Amount: 50},
			},
			expectedOutput: models.CoinHistory{
				Received: []models.CoinTransactionRecived{},
				Sent: []models.CoinTransactionSent{
					{ToUser: "user2", Amount: 50},
				},
			},
		},
		{
			name: "multiple received and sent transactions",
			received: []entities.CoinTransactionInfo{
				{FromUser: "user1", Amount: 100},
				{FromUser: "user3", Amount: 200},
			},
			sent: []entities.CoinTransactionInfo{
				{ToUser: "user2", Amount: 50},
				{ToUser: "user4", Amount: 75},
			},
			expectedOutput: models.CoinHistory{
				Received: []models.CoinTransactionRecived{
					{FromUser: "user1", Amount: 100},
					{FromUser: "user3", Amount: 200},
				},
				Sent: []models.CoinTransactionSent{
					{ToUser: "user2", Amount: 50},
					{ToUser: "user4", Amount: 75},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := MapTransactionEntityToCoinsHistory(tt.received, tt.sent)
			assert.Equal(t, tt.expectedOutput, output)
		})
	}
}
