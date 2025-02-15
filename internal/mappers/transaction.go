package mappers

import (
	"github.com/koccyx/avito_assignment/internal/entities"
	"github.com/koccyx/avito_assignment/internal/http/models"
)

func MapTransactionEntityToCoinsHistory(recived, sent []entities.CoinTransactionInfo) models.CoinHistory {
	trRecived := make([]models.CoinTransactionRecived, 0, len(recived))
	trSent := make([]models.CoinTransactionSent, 0, len(sent))

	for _, tr := range recived {
		trRecived = append(trRecived, models.CoinTransactionRecived{
			FromUser: tr.FromUser,
			Amount: tr.Amount,
		})
	}

	for _, tr := range sent {
		trSent = append(trSent, models.CoinTransactionSent{
			ToUser: tr.ToUser,
			Amount: tr.Amount,
		})
	}

	return models.CoinHistory{
		Received: trRecived,
		Sent: trSent,
	}
}