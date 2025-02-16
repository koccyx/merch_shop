package validators

import (
	"fmt"
	"testing"

	"github.com/koccyx/avito_assignment/internal/server/models"
	"github.com/stretchr/testify/assert"
)

func TestValidateUsername(t *testing.T) {
	tests := []struct {
		name     string
		username string
		expected error
	}{
		{
			name:     "valid username",
			username: "validuser",
			expected: nil,
		},
		{
			name:     "short username",
			username: "usr",
			expected: fmt.Errorf("username must be more then 4"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateUsername(tt.username)
			if tt.expected != nil {
				assert.EqualError(t, err, tt.expected.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidatePassword(t *testing.T) {
	tests := []struct {
		name     string
		password string
		expected error
	}{
		{
			name:     "valid password",
			password: "validpass",
			expected: nil,
		},
		{
			name:     "short password",
			password: "pass",
			expected: fmt.Errorf("password must be more then 5"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePassword(tt.password)
			if tt.expected != nil {
				assert.EqualError(t, err, tt.expected.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateCoinsAmount(t *testing.T) {
	tests := []struct {
		name     string
		amount   int
		expected error
	}{
		{
			name:     "valid amount",
			amount:   10,
			expected: nil,
		},
		{
			name:     "invalid amount",
			amount:   0,
			expected: fmt.Errorf("coins amount must be more then 0"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateCoinsAmount(tt.amount)
			if tt.expected != nil {
				assert.EqualError(t, err, tt.expected.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValdateCoinsTransactionRequest(t *testing.T) {
	tests := []struct {
		name     string
		request  models.SendCoinRequest
		expected error
	}{
		{
			name: "Valid transaction request",
			request: models.SendCoinRequest{
				ToUser:   "user2",
				Amount:   100,
			},
			expected: nil,
		},
		{
			name: "invalid transaction request (missing FromUser)",
			request: models.SendCoinRequest{
				Amount:   100,
			},
			expected: assert.AnError,
		},
		{
			name: "invalid transaction request (Amount < 1)",
			request: models.SendCoinRequest{
				ToUser:   "user2",
			},
			expected: assert.AnError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValdateCoinsTransactionRequest(tt.request)
			if tt.expected != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
