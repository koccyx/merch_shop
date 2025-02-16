package mappers

import (
	"github.com/koccyx/avito_assignment/internal/entities"
	"github.com/koccyx/avito_assignment/internal/server/models"
)

func MapInventoryItemsEntityToModel(items []entities.InventoryItem) []models.InventoryItem {
	invenoryModels := make([]models.InventoryItem, 0, len(items))

	for _, item := range items {
		invenoryModels = append(invenoryModels, models.InventoryItem{
			Type:     item.Name,
			Quantity: item.Amount,
		})
	}

	return invenoryModels
}
