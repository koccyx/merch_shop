package jwt

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type TokenTest struct {
	name   string
	userId string
}

func TestCreateAndParseToken(t *testing.T) {
	secretKey := "avito_assignment"

	tests := []TokenTest{
		{
			name:   "success test",
			userId: "12312assdas",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := NewToken(tt.userId, secretKey)
			require.NoError(t, err)

			prsdToken, err := ParseToken(token, secretKey)
			require.NoError(t, err)

			require.Equal(t, tt.userId, prsdToken)
		})
	}
}

func TestNewToken(t *testing.T) {
	secret := "my_secret"
	userId := "12345"

	token, err := NewToken(userId, secret)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	parsedUserId, err := ParseToken(token, secret)
	require.NoError(t, err)
	assert.Equal(t, userId, parsedUserId)
}

func TestParseToken(t *testing.T) {
	secret := "my_secret"
	userId := "12345"

	token, err := NewToken(userId, secret)
	require.NoError(t, err)

	parsedUserId, err := ParseToken(token, secret)
	require.NoError(t, err)
	assert.Equal(t, userId, parsedUserId)

	_, err = ParseToken(token, "wrong_secret")
	require.Error(t, err)

	_, err = ParseToken("invalid_token", secret)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "token contains an invalid number of segments")
}
