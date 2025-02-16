package mappers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/koccyx/avito_assignment/internal/entities"
	"github.com/koccyx/avito_assignment/internal/server/models"
)

func TestMapInventoryItemsEntityToModel(t *testing.T) {
	tests := []struct {
		name           string
		input          []entities.InventoryItem
		expectedOutput []models.InventoryItem
	}{
		{
			name: "empty input",
			input: []entities.InventoryItem{},
			expectedOutput: []models.InventoryItem{},
		},
		{
			name: "single item",
			input: []entities.InventoryItem{
				{Name: "item1", Amount: 10},
			},
			expectedOutput: []models.InventoryItem{
				{Type: "item1", Quantity: 10},
			},
		},
		{
			name: "multiple items",
			input: []entities.InventoryItem{
				{Name: "item1", Amount: 10},
				{Name: "item2", Amount: 5},
			},
			expectedOutput: []models.InventoryItem{
				{Type: "item1", Quantity: 10},
				{Type: "item2", Quantity: 5},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := MapInventoryItemsEntityToModel(tt.input)
			assert.Equal(t, tt.expectedOutput, output)
		})
	}
}
